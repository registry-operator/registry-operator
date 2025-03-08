// Copyright 2025 The Registry Operator Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	semver "github.com/Masterminds/semver/v3"

	"sigs.k8s.io/yaml"
)

const (
	defaultConfigPath = "./docs/.crd-ref-docs.yaml"
	defaultImageName  = "controller"
	defaultImageRef   = "ghcr.io/registry-operator/registry-operator"
)

func runCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func makeCmd(args ...string) error {
	return runCommand("make", args...)
}

func gitCmd(args ...string) error {
	return runCommand("git", args...)
}

func branchPrep(version, fullVersion string) error {
	branch := fmt.Sprintf("release-%s", version)

	if branchExists(branch) {
		if err := switchToBranch(branch); err != nil {
			return err
		}
	} else {
		if err := createBranch(branch); err != nil {
			return err
		}
	}

	return gitCmd(
		"merge",
		"origin/main",
		"-m", fmt.Sprintf("chore(%s): merge changes for %s", version, fullVersion),
		"--signoff",
	)
}

func branchExists(branchName string) bool {
	cmd := exec.Command("git", "branch", "--list", branchName)
	output, _ := cmd.Output()
	return strings.Contains(string(output), branchName)
}

func switchToBranch(branchName string) error {
	return gitCmd("checkout", branchName)
}

func createBranch(branchName string) error {
	return gitCmd("checkout", "-b", branchName)
}

func parseVersion(version string) (string, error) {
	parsed, err := semver.NewVersion(strings.TrimPrefix(version, "v"))
	if err != nil {
		return "", fmt.Errorf("failed to parse version: %w", err)
	}

	return fmt.Sprintf("%d.%d", parsed.Major(), parsed.Minor()), nil
}

func getLatestKubernetesRelease() (string, error) {
	url := "https://api.github.com/repos/kubernetes/kubernetes/releases/latest"
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to fetch latest Kubernetes release: %w", err)
	}
	defer resp.Body.Close() //nolint:errcheck // best effort call

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("failed to fetch latest Kubernetes release, status code: %d", resp.StatusCode)
	}

	var releaseData struct {
		TagName string `json:"tag_name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&releaseData); err != nil {
		return "", fmt.Errorf("failed to decode Kubernetes release data: %w", err)
	}

	v, err := semver.NewVersion(releaseData.TagName)
	if err != nil {
		return "", fmt.Errorf("failed to parse Kubernetes release semver: %w", err)
	}

	return fmt.Sprintf("%d.%d", v.Major(), v.Minor()), nil
}

func replaceKubernetesVersion(filePath, newVersion string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	lines := strings.Split(string(content), "\n")
	for i, line := range lines {
		if strings.Contains(line, "kubernetesVersion") {
			parts := strings.Split(line, ":")
			lines[i] = fmt.Sprintf("%s: '%s'", parts[0], newVersion)
		}
	}

	tempFilePath := filePath + ".tmp"
	if err := os.WriteFile(tempFilePath, []byte(strings.Join(lines, "\n")), 0644); err != nil {
		return fmt.Errorf("failed to write temp file: %w", err)
	}

	return os.Rename(tempFilePath, filePath)
}

func createKustomization(resources []string, imageName, newImageName, newTag string) map[string]interface{} {
	return map[string]interface{}{
		"namespace":  "registry-operator-system",
		"namePrefix": "registry-operator-",
		"resources":  resources,
		"images": []map[string]string{
			{
				"name":    imageName,
				"newName": newImageName,
				"newTag":  newTag,
			},
		},
	}
}

func writeKustomization(kustomization map[string]interface{}, filepath string) error {
	content, err := yaml.Marshal(kustomization)
	if err != nil {
		return fmt.Errorf("failed to marshal kustomization: %w", err)
	}

	return os.WriteFile(filepath, content, 0644)
}

func release(version, fullVersion string) error {
	if err := gitCmd("add", "."); err != nil {
		return err
	}
	if err := gitCmd(
		"commit",
		"-sm", fmt.Sprintf("chore(%s): create release commit %s", version, fullVersion),
	); err != nil {
		return err
	}
	if err := gitCmd("push", "origin", fmt.Sprintf("release-%s", version)); err != nil {
		return err
	}
	if err := gitCmd("tag", fullVersion); err != nil {
		return err
	}
	return gitCmd("push", "--tags")
}

func main() {
	versionFlag := flag.String("version", "", "Tagged version to build")
	configFlag := flag.String("config", defaultConfigPath, "Path to CRD ref-docs config")
	imageFlag := flag.String("image-name", defaultImageName, "Default image name")
	newImageFlag := flag.String("image", defaultImageRef, "Default image reference")

	flag.Parse()

	resources := []string{"./config/crd", "./config/manager", "./config/rbac"}

	version, err := parseVersion(*versionFlag)
	if err != nil {
		log.Fatalf("Failed to parse version: %v", err)
	}

	if err := makeCmd("manifests", "api-docs"); err != nil {
		log.Fatalf("Failed to run make: %v", err)
	}

	kubeVersion, err := getLatestKubernetesRelease()
	if err != nil {
		log.Fatalf("Failed to fetch latest Kubernetes release: %v", err)
	}

	if err := replaceKubernetesVersion(*configFlag, kubeVersion); err != nil {
		log.Fatalf("Failed to replace Kubernetes version: %v", err)
	}

	if err := branchPrep(version, *versionFlag); err != nil {
		log.Fatalf("Failed to prepare branch: %v", err)
	}

	kustomization := createKustomization(resources, *imageFlag, *newImageFlag, *versionFlag)
	if err := writeKustomization(kustomization, "./kustomization.yaml"); err != nil {
		log.Fatalf("Failed to write kustomization: %v", err)
	}

	if err := release(version, *versionFlag); err != nil {
		log.Fatalf("Failed to release: %v", err)
	}
}
