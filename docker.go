// Copyright 2019 the gamma-cli authors
// This file is part of the gamma-cli library in the Orbs project.
//
// This source code is licensed under the MIT license found in the LICENSE file in the root directory of this source tree.
// The above notice should be included in all copies or substantial portions of the software.

package main

import (
	"bufio"
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

func commandStartLocal(requiredOptions []string) {
	commandStartLocalContainer(gammaHandlerOptions(), requiredOptions)

	if prismEnabled() {
		commandStartLocalContainer(prismHandlerOptions(), requiredOptions)
	}
}

func commandStopLocal(requiredOptions []string) {
	commandStopLocalContainer(gammaHandlerOptions(), requiredOptions)
	commandStopLocalContainer(prismHandlerOptions(), requiredOptions)
}

func commandUpgrade(requiredOptions []string) {
	commandUpgradeImage(gammaHandlerOptions(), requiredOptions)
	commandUpgradeImage(prismHandlerOptions(), requiredOptions)
}

func commandStartLocalContainer(dockerOptions handlerOptions, requiredOptions []string) {
	version := verifyDockerInstalled(dockerOptions, dockerOptions.dockerRegistryTagsUrl)

	if !doesFileExist(*flagKeyFile) {
		commandGenerateTestKeys(nil)
	}

	if isDockerContainerRunning(dockerOptions.containerName) {
		log(`
*********************************************************************************
              %s is already running!

  Run 'gamma-cli help' in terminal to learn how to interact with this instance.
              
**********************************************************************************
`, dockerOptions.name)
		exit()
	}

	if err := createDockerNetwork(); err != nil {
		die("could not create docker network gamma: %s", err)
	}

	p := fmt.Sprintf("%d:%d", dockerOptions.port, dockerOptions.containerPort)
	run := fmt.Sprintf("%s:%s", dockerOptions.dockerRepo, version)
	args := []string {
		"run", "-d",
		"--name", dockerOptions.containerName,
		"-p", p,
		"--network", "gamma",
	}
	for _, value := range dockerOptions.env {
		args = append(args, "-e", value)
	}
	args = append(args, run)
	args = append(args, dockerOptions.dockerCmd...)

	out, err := exec.Command("docker",  args...).CombinedOutput()
	if err != nil {
		die("Could not exec 'docker run' command.\n\n%s", out)
	}

	if !isDockerContainerRunning(dockerOptions.containerName) {
		die("Could not run docker image.")
	}

	if *flagWait {
		waitUntilDockerIsReadyAndListening(IS_READY_TOTAL_WAIT_TIMEOUT)
	}

	log(`
*********************************************************************************
                 %s %s is running!

  Local blockchain instance started and listening on port %d.
  Run 'gamma-cli help' in terminal to learn how to interact with this instance.
              
**********************************************************************************
`, dockerOptions.name, version, dockerOptions.port)
}

func commandStopLocalContainer(dockerOptions handlerOptions, requiredOptions []string) {
	verifyDockerInstalled(dockerOptions, dockerOptions.dockerRegistryTagsUrl)

	out, err := exec.Command("docker", "stop", dockerOptions.containerName).CombinedOutput()
	if err != nil {
		log("%s server is already stopped.\n", dockerOptions.name)
	}

	out, err = exec.Command("docker", "rm", "-f", dockerOptions.containerName).CombinedOutput()
	if err != nil {
		log("Could not remove docker container.\n\n%s", out)
	}

	if isDockerContainerRunning(dockerOptions.containerName) {
		die("Could not stop docker container.")
	}

	log(`
*********************************************************************************
                    %s stopped.

  A local blockchain instance is running in-memory.
  The next time you start the instance, all contracts and state will disappear. 
              
**********************************************************************************
`, dockerOptions.name)
}

func commandUpgradeImage(dockerOptions handlerOptions, requiredOptions []string) {
	currentTag := verifyDockerInstalled(dockerOptions, dockerOptions.dockerRegistryTagsUrl)
	latestTag := getLatestDockerTag(dockerOptions.dockerRegistryTagsUrl)

	if !isExperimental() && cmpTags(latestTag, currentTag) <= 0 {
		log("Current %s stable version %s does not require upgrade.", dockerOptions.name, currentTag)
	} else {
		log("Downloading latest version %s:\n", latestTag)
		cmd := exec.Command("docker", "pull", fmt.Sprintf("%s:%s", dockerOptions.dockerRepo, latestTag))
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
		log("")
	}
}

func showLogs(requiredOptions []string) {
	dockerOptions := gammaHandlerOptions()
	verifyDockerInstalled(dockerOptions, dockerOptions.dockerRegistryTagsUrl)

	cmd := exec.Command("docker", "logs", "-f", "--tail=20", dockerOptions.containerName)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		die("could not read gamma server docker logs: %s", err)
	}

	if err := cmd.Start(); err != nil {
		die("could not start docker logs command: %s", err)
	}

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := scanner.Text()

		if testGammaLogLine(line)  {
			fmt.Println(line)
		}
	}
	if err := scanner.Err(); err != nil {
		die("error reading gamma server docker logs: %s", err)
	}
}

// TODO remove dockerRegistryUrl as separate parameter
func verifyDockerInstalled(dockerOptions handlerOptions, dockerRegistryTagUrl string) string {
	out, err := exec.Command("docker", "images", dockerOptions.dockerRepo).CombinedOutput()
	if err != nil {
		if runtime.GOOS == "darwin" {
			die("Docker is required but not running. Is it installed on your machine?\n\nInstall from:  https://docs.docker.com/docker-for-mac/install/")
		} else {
			die("Docker is required but not running. Is it installed on your machine?\n\nInstall from:  https://docs.docker.com/install/")
		}
	}

	existingTag := extractTagFromDockerImagesOutput(dockerOptions.dockerRepo, string(out))
	if existingTag != DOCKER_TAG_NOT_FOUND {
		return existingTag
	}

	latestTag := getLatestDockerTag(dockerRegistryTagUrl)

	log("%s image is not installed, downloading version %s:\n", dockerOptions.name, latestTag)
	cmd := exec.Command("docker", "pull", fmt.Sprintf("%s:%s", dockerOptions.dockerRepo, latestTag))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
	log("")

	out, err = exec.Command("docker", "images", dockerOptions.dockerRepo).CombinedOutput()
	if err != nil || strings.Count(string(out), "\n") == 1 {
		die("Could not download docker image.")
	}
	return extractTagFromDockerImagesOutput(dockerOptions.dockerRepo, string(out))
}

func isDockerContainerRunning(containerName string) bool {
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

func createDockerNetwork() error {
	out, err := exec.Command("docker", "network", "ls", "--filter", "name=gamma", "-q").CombinedOutput()
	if err != nil {
		return err
	}

	if len(out) == 0 {
		_, err := exec.Command("docker", "network", "create", "gamma").CombinedOutput()
		if err != nil {
			return err
		}
	}

	return nil
}

func prismEnabled() bool {
	return !*flagNoUi
}

func testGammaLogLine(line string) bool {
	return !strings.Contains(line, "info ") && !strings.HasPrefix(line, "error ") &&
		!strings.HasPrefix(line, "debug ") && !strings.HasPrefix(line, "\t") &&
		!strings.HasPrefix(line, "] ")
}