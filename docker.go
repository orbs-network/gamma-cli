package main

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

const DOCKER_REPO = "orbsnetwork/gamma"
const DOCKER_RUN = "orbsnetwork/gamma:%s"
const CONTAINER_NAME = "orbs-gamma-server"
const DOCKER_REGISTRY_TAGS_URL = "https://registry.hub.docker.com/v2/repositories/orbsnetwork/gamma/tags/"
const DOCKER_TAG_NOT_FOUND = "not found"
const DOCKER_TAG_EXPERIMENTAL = "experimental"

func commandStartLocal(requiredOptions []string) {
	gammaVersion := verifyDockerInstalled()

	if !doesFileExist(*flagKeyFile) {
		commandGenerateTestKeys(nil)
	}

	if isDockerGammaRunning() {
		log(`
*********************************************************************************
              Orbs Gamma personal blockchain is already running!

  Run 'gamma-cli help' in terminal to learn how to interact with this instance.
              
**********************************************************************************
`)
		exit()
	}

	p := fmt.Sprintf("%d:8080", *flagPort)
	run := fmt.Sprintf(DOCKER_RUN, gammaVersion)
	out, err := exec.Command("docker", "run", "-d", "--name", CONTAINER_NAME, "-p", p, run, "./gamma-server", "-override-config", *flagOverrideConfig).CombinedOutput()
	if err != nil {
		die("Could not exec 'docker run' command.\n\n%s", out)
	}

	if !isDockerGammaRunning() {
		die("Could not run docker image.")
	}

	if *flagWait {
		waitUntilDockerIsReadyAndListening(IS_READY_TOTAL_WAIT_TIMEOUT)
	}

	log(`
*********************************************************************************
                 Orbs Gamma %s personal blockchain is running!

  Local blockchain instance started and listening on port %d.
  Run 'gamma-cli help' in terminal to learn how to interact with this instance.
              
**********************************************************************************
`, gammaVersion, *flagPort)
}

func commandStopLocal(requiredOptions []string) {
	verifyDockerInstalled()

	out, err := exec.Command("docker", "stop", CONTAINER_NAME).CombinedOutput()
	if err != nil {
		log("Gamma server is already stopped.\n")
		exit()
	}

	out, err = exec.Command("docker", "rm", "-f", CONTAINER_NAME).CombinedOutput()
	if err != nil {
		die("Could not remove docker container.\n\n%s", out)
	}

	if isDockerGammaRunning() {
		die("Could not stop docker container.")
	}

	log(`
*********************************************************************************
                    Orbs Gamma personal blockchain stopped.

  A local blockchain instance is running in-memory.
  The next time you start the instance, all contracts and state will disappear. 
              
**********************************************************************************
`)
}

func commandUpgradeServer(requiredOptions []string) {
	currentTag := verifyDockerInstalled()
	latestTag := getLatestDockerTag()

	if !isExperimental() && cmpTags(latestTag, currentTag) <= 0 {
		log("Current Gamma server stable version %s does not require upgrade.\n", currentTag)
		exit()
	}

	log("Downloading latest version %s:\n", latestTag)
	cmd := exec.Command("docker", "pull", fmt.Sprintf("%s:%s", DOCKER_REPO, latestTag))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
	log("")
}

func verifyDockerInstalled() string {
	out, err := exec.Command("docker", "images", DOCKER_REPO).CombinedOutput()
	if err != nil {
		if runtime.GOOS == "darwin" {
			die("Docker is required but not running. Is it installed on your machine?\n\nInstall from:  https://docs.docker.com/docker-for-mac/install/")
		} else {
			die("Docker is required but not running. Is it installed on your machine?\n\nInstall from:  https://docs.docker.com/install/")
		}
	}

	existingTag := extractTagFromDockerImagesOutput(string(out))
	if existingTag != DOCKER_TAG_NOT_FOUND {
		return existingTag
	}

	latestTag := getLatestDockerTag()

	log("Orbs personal blockchain docker image is not installed, downloading version %s:\n", latestTag)
	cmd := exec.Command("docker", "pull", fmt.Sprintf("%s:%s", DOCKER_REPO, latestTag))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
	log("")

	out, err = exec.Command("docker", "images", DOCKER_REPO).CombinedOutput()
	if err != nil || strings.Count(string(out), "\n") == 1 {
		die("Could not download docker image.")
	}
	return extractTagFromDockerImagesOutput(string(out))
}

func isDockerGammaRunning() bool {
	out, err := exec.Command("docker", "ps", "-f", fmt.Sprintf("name=%s", CONTAINER_NAME)).CombinedOutput()
	if err != nil {
		return false
	}
	return strings.Count(string(out), "\n") > 1
}

func extractTagFromDockerImagesOutput(out string) string {
	pattern := fmt.Sprintf(`%s\s+(v\S+)`, regexp.QuoteMeta(DOCKER_REPO))
	if isExperimental() {
		pattern = fmt.Sprintf(`%s\s+(%s)`, regexp.QuoteMeta(DOCKER_REPO), regexp.QuoteMeta(DOCKER_TAG_EXPERIMENTAL))
	}
	re := regexp.MustCompile(pattern)
	res := re.FindStringSubmatch(out)
	if len(res) < 2 {
		return DOCKER_TAG_NOT_FOUND
	}
	return res[1]
}

func getLatestDockerTag() string {
	if isExperimental() {
		return DOCKER_TAG_EXPERIMENTAL
	}
	resp, err := http.Get(DOCKER_REGISTRY_TAGS_URL)
	if err != nil {
		die("Cannot connect to docker registry to get image list.")
	}
	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil || len(bytes) == 0 {
		die("Bad image list response from docker registry.")
	}
	tag, err := extractLatestTagFromDockerHubResponse(bytes)
	if err != nil {
		die("Cannot parse image list response from docker registry.")
	}
	return tag
}

type dockerHubTagsJson struct {
	Results []*struct {
		Name string
	}
}

func extractLatestTagFromDockerHubResponse(responseBytes []byte) (string, error) {
	var response *dockerHubTagsJson
	err := json.Unmarshal(responseBytes, &response)
	if err != nil {
		return "", err
	}
	maxTag := ""
	for _, result := range response.Results {
		if cmpTags(result.Name, maxTag) > 0 {
			maxTag = result.Name
		}
	}
	if maxTag == "" {
		return "", errors.New("no valid tags found")
	}
	return maxTag, nil
}

func isExperimental() bool {
	if *flagEnv == LOCAL_ENV_ID {
		return false
	}
	if *flagEnv == EXPERIMENTAL_ENV_ID {
		return true
	}
	// allow custom experimental configs (like for docker ci)
	return getEnvironmentFromConfigFile(*flagEnv).Experimental
}

func cmpTags(t1, t2 string) int {
	re := regexp.MustCompile(`v(\d+)\.(\d+)\.(\d+)`)
	m1 := re.FindStringSubmatch(t1)
	if len(m1) < 4 {
		return -1
	}
	m2 := re.FindStringSubmatch(t2)
	if len(m2) < 4 {
		return 1
	}
	diff := atoi(m1[1]) - atoi(m2[1])
	if diff != 0 {
		return diff
	}
	diff = atoi(m1[2]) - atoi(m2[2])
	if diff != 0 {
		return diff
	}
	return atoi(m1[3]) - atoi(m2[3])
}

func atoi(num string) int {
	res, err := strconv.Atoi(num)
	if err != nil {
		return 0
	}
	return res
}
