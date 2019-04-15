// Copyright 2019 the gamma-cli authors
// This file is part of the gamma-cli library in the Orbs project.
//
// This source code is licensed under the MIT license found in the LICENSE file in the root directory of this source tree.
// The above notice should be included in all copies or substantial portions of the software.

package test

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"time"
)

const CONFIG_FILENAME = "orbs-gamma-config.json"

var cachedGammaCliBinaryPath string
var downloadedLatestGammaServer bool

type gammaCli struct {
	port            string
	experimental    bool
	configOverrides string
}

func compileGammaCli() string {
	if gammaCliPathFromEnv := os.Getenv("GAMMA_CLI_PATH"); gammaCliPathFromEnv != "" {
		return gammaCliPathFromEnv
	}

	if cachedGammaCliBinaryPath != "" {
		return cachedGammaCliBinaryPath // cache compilation once per process
	}

	tempDir, err := ioutil.TempDir("", "gamma")
	if err != nil {
		panic(err)
	}

	binaryOutputPath := tempDir + "/gamma-cli"
	goCmd := path.Join(runtime.GOROOT(), "bin", "go")
	cmd := exec.Command(goCmd, "build", "-o", binaryOutputPath, ".")
	cmd.Dir = path.Join(getCurrentSourceFileDirPath(), "..")
	out, err := cmd.CombinedOutput()
	if err != nil {
		panic(fmt.Sprintf("compilation failed: %s\noutput:\n%s\n", err.Error(), out))
	} else {
		fmt.Printf("compiled gamma-cli successfully:\n %s\n", binaryOutputPath)
	}

	cachedGammaCliBinaryPath = binaryOutputPath
	return cachedGammaCliBinaryPath
}

func (g *gammaCli) Run(args ...string) (string, error) {
	if len(args) > 0 {
		// streamlined supprt for a different port
		args = append(args, "-port", g.port)

		// needed for dockerized tests in ci
		args = append(args, "-env", g.getGammaEnvironment())
		pathToConfig := path.Join(os.Getenv("HOME"), ".orbs", CONFIG_FILENAME)
		args = append(args, "-config", pathToConfig)

		if g.configOverrides != "" {
			args = append(args, "-override-config", g.configOverrides)
		}
	}
	out, err := exec.Command(compileGammaCli(), args...).CombinedOutput()
	return string(out), err
}

func GammaCli() *gammaCli {
	return GammaCliWithPort(8080)
}

func GammaCliWithPort(port int) *gammaCli {
	return &gammaCli{
		port: fmt.Sprintf("%d", port),
	}
}

func (g *gammaCli) WithConfigOverrides(override string) *gammaCli {
	g.configOverrides = override
	return g
}

func (g *gammaCli) WithStableServer() *gammaCli {
	g.experimental = false
	return g
}

func (g *gammaCli) WithExperimentalServer() *gammaCli {
	g.experimental = true
	return g
}

func (g *gammaCli) StartGammaServer() *gammaCli {
	out, err := g.Run("start-local", "-wait")
	if err != nil {
		panic(fmt.Sprintf("start Gamma server failed: %s\noutput:\n%s\n", err.Error(), out))
	}
	fmt.Println(out)
	return g
}

func (g *gammaCli) StopGammaServer() {
	out, err := g.Run("stop-local")
	if err != nil {
		panic(fmt.Sprintf("stop Gamma server failed: %s\noutput:\n%s\n", err.Error(), out))
	}
}

func (g *gammaCli) DownloadLatestGammaServer() *gammaCli {
	if downloadedLatestGammaServer {
		return g
	}
	downloadedLatestGammaServer = true

	start := time.Now()
	out, err := g.Run("upgrade-server")
	if err != nil {
		panic(fmt.Sprintf("download latest Gamma server failed: %s\noutput:\n%s\n", err.Error(), out))
	}
	delta := time.Now().Sub(start)
	fmt.Printf("upgraded gamma-server to latest version (this took %.3fs)\n", delta.Seconds())
	return g
}

func (g *gammaCli) getGammaEnvironment() string {
	if env := os.Getenv("GAMMA_ENVIRONMENT"); env != "" {
		if g.experimental {
			return env + "-experimental"
		}
		return env
	}

	if g.experimental {
		return "experimental"
	}
	return "local"
}

func extractTxIdFromSendTxOutput(out string) string {
	re := regexp.MustCompile(`\"TxId\":\s+\"(\w+)\"`)
	res := re.FindStringSubmatch(out)
	return res[1]
}

func getCurrentSourceFileDirPath() string {
	_, filename, _, _ := runtime.Caller(1)
	return filepath.Dir(filename)
}
