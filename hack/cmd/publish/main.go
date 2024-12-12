package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"sort"
	"strings"

	semver "github.com/Masterminds/semver/v3"
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

	sv, err := semver.NewVersion(strings.TrimPrefix(version, "v"))
	if err != nil {
		return "", false, fmt.Errorf("failed to parse version: %w", err)
	}

	shortVersion := fmt.Sprintf("v%d.%d", sv.Major(), sv.Minor())
	isPrerelease := len(sv.Prerelease()) > 0

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

// runCommandOutput runs a system command and returns its output as a string.
func runCommandOutput(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to run command %s: %w", name, err)
	}

	return string(output), nil
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

// cleanupOldVersions deletes the oldest versions if there are more than 3 versions.
func cleanupOldVersions() error {
	log.Println("Fetching deployed versions")
	output, err := runCommandOutput("mike", "list", "--json")
	if err != nil {
		return fmt.Errorf("failed to fetch deployed versions: %w", err)
	}

	var versions []struct {
		Version string   `json:"version"`
		Title   string   `json:"title"`
		Aliases []string `json:"aliases"`
	}
	if err := json.Unmarshal([]byte(output), &versions); err != nil {
		return fmt.Errorf("failed to parse JSON output: %w", err)
	}

	if len(versions) <= 4 {
		log.Println("No cleanup needed. Versions count:", len(versions))
		return nil
	}

	log.Println("Sorting versions")
	sortedVersions := make([]*semver.Version, 0, len(versions))
	for _, v := range versions {
		parsed, err := semver.NewVersion(strings.TrimPrefix(v.Version, "v"))
		if err != nil {
			log.Printf("Skipping invalid version: %s", v.Version)
			continue
		}
		sortedVersions = append(sortedVersions, parsed)
	}

	sort.Sort(semver.Collection(sortedVersions))

	log.Printf("Deleting oldest versions. Total versions: %d", len(sortedVersions))
	for i := 0; i < len(sortedVersions)-3; i++ {
		oldVersion := sortedVersions[i].Original()
		log.Printf("Deleting version: %s", oldVersion)
		if err := mike("delete", "--push", "v"+oldVersion); err != nil {
			return fmt.Errorf("failed to delete version %s: %w", oldVersion, err)
		}
	}

	return nil
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

	if err := cleanupOldVersions(); err != nil {
		log.Fatalf("Error during cleanup: %v", err)
	}
}
