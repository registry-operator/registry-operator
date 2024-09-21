package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/blang/semver/v4"
)

const (
	defaultVersion = "main"
	latestVersion  = "latest"
)

// parseVersion parses the given version string and returns a short version and a boolean indicating
// if it's a prerelease.
func parseVersion(version string) (string, bool, error) {
	if version == defaultVersion {
		return version, false, nil
	}

	sv, err := semver.ParseTolerant(strings.TrimPrefix(version, "v"))
	if err != nil {
		return "", false, fmt.Errorf("failed to parse version: %w", err)
	}

	shortVersion := fmt.Sprintf("v%d.%d", sv.Major, sv.Minor)
	isPrerelease := len(sv.Pre) > 0

	log.Printf("Parsed version: %s (is_prerelease: %v)", shortVersion, isPrerelease)
	return shortVersion, isPrerelease, nil
}

// runCommand runs a system command and pipes its output to stdout and stderr.
func runCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run command %s: %w", name, err)
	}

	return nil
}

// mike runs the "mike" command with the provided arguments.
func mike(args ...string) error {
	return runCommand("mike", args...)
}

// isInitial checks if the repository is in an initial state by verifying the gh-pages branch.
func isInitial() (bool, error) {
	log.Println("Fetching gh-pages from origin")
	if err := runCommand("git", "remote", "update"); err != nil {
		return false, fmt.Errorf("failed to update remote: %w", err)
	}

	log.Println("Checking gh-pages reference")
	cmd := exec.Command("git", "show-ref", "refs/remotes/origin/gh-pages")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()

	log.Printf("Show-ref exit code: %v", cmd.ProcessState.ExitCode())
	return err != nil, nil
}

// deployLatest deploys the latest version and sets it as the default.
func deployLatest(version string) error {
	log.Println("Deploying latest version")
	if err := mike("deploy", "--push", "--update-aliases", version, latestVersion); err != nil {
		return fmt.Errorf("failed to deploy latest: %w", err)
	}

	log.Println("Setting default version to latest")
	if err := mike("set-default", "--push", latestVersion); err != nil {
		return fmt.Errorf("failed to set default: %w", err)
	}

	return nil
}

// deployVersion deploys a specific version.
func deployVersion(version string, prerelease bool) error {
	if prerelease {
		log.Printf("Deploying prerelease version: %s", version)
		return mike("deploy", "--push", version)
	}

	log.Printf("Deploying release version: %s", version)
	if version == defaultVersion {
		return mike("deploy", "--push", "--update-aliases", version)
	}

	return mike("deploy", "--push", "--update-aliases", version, latestVersion)
}

func main() {
	versionFlag := flag.String("version", defaultVersion, "Tagged version to be built")
	flag.Parse()

	version, prerelease, err := parseVersion(*versionFlag)
	if err != nil {
		log.Fatalf("Error parsing version: %v", err)
	}

	initial, err := isInitial()
	if err != nil {
		log.Fatalf("Error checking repository state: %v", err)
	}

	if initial {
		log.Println("Initial release detected")
		if err := deployLatest(version); err != nil {
			log.Fatalf("Error deploying initial release: %v", err)
		}
	} else {
		log.Printf("Not an initial release (prerelease=%v, version=%s)", prerelease, version)
		if err := deployVersion(version, prerelease); err != nil {
			log.Fatalf("Error deploying version: %v", err)
		}
	}
}
