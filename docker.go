// Copyright 2019 the gamma-cli authors
// This file is part of the gamma-cli library in the Orbs project.
//
// This source code is licensed under the MIT license found in the LICENSE file in the root directory of this source tree.
// The above notice should be included in all copies or substantial portions of the software.

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

const DOCKER_TAG_NOT_FOUND = "not found"
const DOCKER_TAG_EXPERIMENTAL = "experimental"

func commandStartLocal(dockerOptions handlerOptions, requiredOptions []string) {
	gammaVersion := verifyDockerInstalled(dockerOptions.dockerRepo, dockerOptions.dockerRegistryTagsUrl)

	if !doesFileExist(*flagKeyFile) {
		commandGenerateTestKeys(nil)
	}

	if isDockerGammaRunning(dockerOptions.containerName) {
		log(`
*********************************************************************************
              Orbs Gamma personal blockchain is already running!

  Run 'gamma-cli help' in terminal to learn how to interact with this instance.
              
**********************************************************************************
`)
		exit()
	}

	p := fmt.Sprintf("%d:%d", dockerOptions.port, dockerOptions.containerPort)
	run := fmt.Sprintf(dockerOptions.dockerRun, gammaVersion)
	out, err := exec.Command("docker", "run", "-d", "--name", dockerOptions.containerName, "-p", p, run).CombinedOutput()
	if err != nil {
		die("Could not exec 'docker run' command.\n\n%s", out)
	}

	if !isDockerGammaRunning(dockerOptions.containerName) {
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
`, gammaVersion, dockerOptions.port)
}

func commandStopLocal(dockerOptions handlerOptions, requiredOptions []string) {
	verifyDockerInstalled(dockerOptions.dockerRepo, dockerOptions.dockerRegistryTagsUrl)

	out, err := exec.Command("docker", "stop", dockerOptions.containerName).CombinedOutput()
	if err != nil {
		log("Gamma server is already stopped.\n")
		exit()
	}

	out, err = exec.Command("docker", "rm", "-f", dockerOptions.containerName).CombinedOutput()
	if err != nil {
		die("Could not remove docker container.\n\n%s", out)
	}

	if isDockerGammaRunning(dockerOptions.containerName) {
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

func commandUpgradeServer(dockerOptions handlerOptions, requiredOptions []string) {
	currentTag := verifyDockerInstalled(dockerOptions.dockerRepo, dockerOptions.dockerRegistryTagsUrl)
	latestTag := getLatestDockerTag(dockerOptions.dockerRegistryTagsUrl)

	if !isExperimental() && cmpTags(latestTag, currentTag) <= 0 {
		log("Current Gamma server stable version %s does not require upgrade.\n", currentTag)
		exit()
	}

	log("Downloading latest version %s:\n", latestTag)
	cmd := exec.Command("docker", "pull", fmt.Sprintf("%s:%s", dockerOptions.dockerRepo, latestTag))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
	log("")
}

func verifyDockerInstalled(dockerRepo string, dockerRegistryTagUrl string) string {
	out, err := exec.Command("docker", "images", dockerRepo).CombinedOutput()
	if err != nil {
		if runtime.GOOS == "darwin" {
			die("Docker is required but not running. Is it installed on your machine?\n\nInstall from:  https://docs.docker.com/docker-for-mac/install/")
		} else {
			die("Docker is required but not running. Is it installed on your machine?\n\nInstall from:  https://docs.docker.com/install/")
		}
	}

	existingTag := extractTagFromDockerImagesOutput(dockerRepo, string(out))
	if existingTag != DOCKER_TAG_NOT_FOUND {
		return existingTag
	}

	latestTag := getLatestDockerTag(dockerRegistryTagUrl)

	log("Orbs personal blockchain docker image is not installed, downloading version %s:\n", latestTag)
	cmd := exec.Command("docker", "pull", fmt.Sprintf("%s:%s", dockerRepo, latestTag))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
	log("")

	out, err = exec.Command("docker", "images", dockerRepo).CombinedOutput()
	if err != nil || strings.Count(string(out), "\n") == 1 {
		die("Could not download docker image.")
	}
	return extractTagFromDockerImagesOutput(dockerRepo, string(out))
}

func isDockerGammaRunning(containerName string) bool {
	out, err := exec.Command("docker", "ps", "-f", fmt.Sprintf("name=%s", containerName)).CombinedOutput()
	if err != nil {
		return false
	}
	return strings.Count(string(out), "\n") > 1
}

func extractTagFromDockerImagesOutput(dockerRepo string, out string) string {
	pattern := fmt.Sprintf(`%s\s+(v\S+)`, regexp.QuoteMeta(dockerRepo))
	if isExperimental() {
		pattern = fmt.Sprintf(`%s\s+(%s)`, regexp.QuoteMeta(dockerRepo), regexp.QuoteMeta(DOCKER_TAG_EXPERIMENTAL))
	}
	re := regexp.MustCompile(pattern)
	res := re.FindStringSubmatch(out)
	if len(res) < 2 {
		return DOCKER_TAG_NOT_FOUND
	}
	return res[1]
}

func getLatestDockerTag(dockerRegistryTagsUrl string) string {
	if isExperimental() {
		return DOCKER_TAG_EXPERIMENTAL
	}
	resp, err := http.Get(dockerRegistryTagsUrl)
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
