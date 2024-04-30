# Contributing to registry-operator

**Table of Contents**

- [Contributing to registry-operator](#contributing-to-registry-operator)
  - [Issues](#issues)
    - [Reporting an Issue](#reporting-an-issue)
    - [Issue Lifecycle](#issue-lifecycle)
  - [Pull Requests](#pull-requests)
  - [Developing](#developing)
    - [Go Environment and Go Modules](#go-environment-and-go-modules)
    - [Code Linting with golangci-lint](#code-linting-with-golangci-lint)
    - [Testing](#testing)
      - [Writing Tests](#writing-tests)
        - [Unit tests](#unit-tests)
        - [Integration tests](#integration-tests)
        - [End-to-end tests](#end-to-end-tests)
    - [Reducing Third-Party Libraries](#reducing-third-party-libraries)
      - [Guidelines](#guidelines)
  - [Releasing](#releasing)
    - [Tagging a release](#tagging-a-release)
      - [Prerequisites](#prerequisites)
      - [Tagging the release](#tagging-the-release)
    - [If a release fails](#if-a-release-fails)
      - [Github Releases](#github-releases)
        - [Prerequisites](#prerequisites-1)

**First:** if you're unsure or afraid of _anything_, just ask or submit the issue or pull request anyways. You won't be yelled at for giving your best effort. The worst that can happen is that you'll be politely asked to change something. We appreciate all contributions!

For those folks who want a bit more guidance on the best way to contribute to the project, read on. Addressing the points below lets us merge or address your contributions quickly.

## Issues

### Reporting an Issue

* Make sure you test against the latest released version. It is possible we already fixed the bug you're experiencing.
* If you experienced a panic, please create a [gist](https://gist.github.com) of the *entire* generated crash log for us to look at. Double check no sensitive items were in the log.
* Respond as promptly as possible to any questions made by the _registry-operator_ team to your issue. Stale issues will be closed.

### Issue Lifecycle

1. The issue is reported.
2. The issue is verified and categorized by a _registry-operator_ collaborator. Categorization is done via labels. For example, bugs are marked as "bugs".
3. Unless it is critical, the issue is left for a period of time (sometimes many weeks), giving outside contributors a chance to address the issue.
4. The issue is addressed in a pull request. The issue will be referenced in commit message(s) so that the code that fixes it is clearly linked.
5. The issue is closed. Sometimes, valid issues will be closed to keep the issue tracker clean. The issue is still indexed and available for future viewers, or can be re-opened if necessary.

## Pull Requests

Pull requests must always be opened from a fork of `registry-operator`, even if you have
commit rights to the repository so that all contributors follow the same process.

## Developing

### Go Environment and Go Modules

To contribute to registry-operator, you need to have Go installed on your system and set up with Go modules. Follow these steps to get started:

1. Install Go:
   - For macOS users, the recommended way is to use Homebrew:
     ```
     $ brew install go
     ```
   - For other platforms or manual installation, you can download and install Go from the [official website](https://golang.org/dl/).

1. Clone the `registry-operator` repository to your local machine:
   ```
   $ git clone https://github.com/$YOUR_USERNAME/registry-operator.git
   ```

2. Change into the `registry-operator` directory:
   ```
   $ cd registry-operator
   ```

3. Now you're all set with the Go environment and Go modules!

### Code Linting with golangci-lint

To ensure consistent code quality, we use `golangci-lint` as a single point for code linting. You can install `golangci-lint` via Homebrew (for macOS users) or using the `go install` command (for all platforms).

With `golangci-lint` installed, you can now run it against the registry-operator codebase to check for any linting issues:

```sh
$ make lint
```

Fix any linting issues reported by `golangci-lint` before submitting your changes.

Remember, we encourage contributions to be well-formatted and follow the project's coding conventions. Happy coding!

### Testing

#### Writing Tests

When adding new features or fixing bugs, it's essential to write tests to ensure the stability and correctness of the code changes. `registry-operator` uses both unit tests and integration tests.

##### Unit tests

Unit tests focus on testing individual functions and components in isolation. To write a unit test, create a new file in the `*_test.go` format alongside the code you want to test. Use the Go testing framework to create test functions that cover different scenarios and edge cases.

Example unit test:

[![Example](https://img.shields.io/badge/Run-go.dev%2Fplay-29BEB0?logo=go)](https://go.dev/play/p/Gm6M3vNVQOw)

```go
package main

import (
	"testing"
)

func Add(a, b int) int {
	return a + b
}

func TestAdd(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name       string
		arg1, arg2 int
		expected   int
	}{
		{arg1: 1, arg2: 2, expected: 3},
		{arg1: 2, arg2: 3, expected: 5},
		{arg1: -2, arg2: 3, expected: 1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			actual := Add(tc.arg1, tc.arg2)
			if tc.expected != actual {
				t.Errorf("Expected the sum of %d and %d to be %d, got %d",
					tc.arg1, tc.arg2,
					tc.expected,
					actual)
			}
		})
	}
}
```

##### Integration tests

Integration tests check the interaction between different parts of the system and may involve external dependencies like databases or APIs. Write integration tests in separate test files with appropriate names to differentiate them from unit tests.

Example integration test:

[![Example](https://img.shields.io/badge/Run-go.dev%2Fplay-29BEB0?logo=go)](https://go.dev/play/p/CquI-e2tXW0)

```go
// go:build integration
package main

import (
	"os"
	"testing"
)

type Client struct {
	apiKey string
}

func (c *Client) APIAdd(a, b int) int {
	return a + b
}

func TestAdd(t *testing.T) {
	t.Parallel()

	apiKey, ok := os.LookupEnv("API_KEY")
	if !ok {
		t.Fatalf("API_KEY must be set for test to run")
	}

	defaultClient := &Client{apiKey: apiKey}

	for _, tc := range []struct {
		name       string
		client     *Client
		arg1, arg2 int
		expected   int
	}{
		{client: defaultClient, arg1: 1, arg2: 2, expected: 3},
		{client: defaultClient, arg1: 2, arg2: 3, expected: 5},
		{client: defaultClient, arg1: -2, arg2: 3, expected: 1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			actual := tc.client.APIAdd(tc.arg1, tc.arg2)
			if tc.expected != actual {
				t.Errorf("Expected the sum of %d and %d to be %d, got %d",
					tc.arg1, tc.arg2,
					tc.expected,
					actual)
			}
		})
	}
}
```

##### End-to-end tests

[Kyverno Chainsaw](http://kyverno.github.io/chainsaw/latest/) is a low-code, declarative tool for writing end-to-end (E2E) tests for Kyverno policies. It streamlines the process of creating and maintaining E2E tests by providing a user-friendly syntax and a powerful assertion model.

For more info refer to the [Kyverno Chainsaw documentation](https://kyverno.github.io/chainsaw/latest/writing-tests/).

### Reducing Third-Party Libraries

We strive to minimize dependencies on third-party libraries in `registry-operator` for the following reasons:

- **Complexity:** Each additional library introduces potential vulnerabilities and increases the overall complexity of the codebase.
- **Maintenance:** Maintaining external dependencies can be time-consuming, especially when libraries are no longer actively supported.
- **Flexibility:** Limiting third-party libraries allows us to have more control over the functionality and behavior of the operator.

#### Guidelines

- Before introducing a new third-party library, carefully consider:
    - **Alternatives:** Are there existing implementations within the Go standard library that could be used instead?
    - **Scope:** Does the library offer exactly the functionality needed, or does it include unnecessary features?
    - **Activity:** Is the library actively maintained and receiving updates?
    - **License:** Is the license compatible with the `registry-operator` project?
- **Prioritize well-established and actively maintained libraries.**
- **Document the rationale for including each third-party library.**
- **Explore opportunities to remove unused or outdated libraries.**

**Remember:** While minimizing third-party libraries is important, it should not come at the expense of functionality or maintainability. Always strive for a balance between these factors.

**By following these guidelines, you can help us keep `registry-operator` lean, efficient, and easy to maintain.**

## Releasing

### Tagging a release

When it's time to make a new release of `registry-operator`, follow these steps to tag the release:

#### Prerequisites

Before tagging a release, make sure:
- All changes for the release are merged into the `main` branch.


#### Tagging the release

Merge PR opened by [release-please-action](https://github.com/google-github-actions/release-please-action) integration.

### If a release fails

If a release fails for any reason, follow these steps to handle the situation:

#### Github Releases

##### Prerequisites

Before attempting to create a new release, make sure:
- You have the necessary permissions to create a release on the repository.
- Check the release workflow and ensure it's properly configured to handle the release process.

If the release workflow fails:
1. Investigate the cause of the failure by reviewing the logs and error messages.
2. Make necessary fixes to the release workflow or the repository configuration.
3. Retry the release workflow.

Remember, creating a new release is a critical process, so always double-check everything before proceeding.

Remember, this document is a starting point for contributors to understand how to work with registry-operator and contribute effectively. It's important to keep it up to date and include any changes in the development and contribution processes over time.

Feel free to extend and modify this document to reflect any new practices or guidelines for contributing to registry-operator. Happy contributing!
