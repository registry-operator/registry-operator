# E2E Tests

This section outlines how to set up and run end-to-end (E2E) tests for registry operator. E2E tests ensure that the entire application flow, from start to finish, functions as expected.

## Quickstart

To quickly set up and run the E2E tests, follow these steps:

1. **Install Dependencies**: Begin by installing the necessary dependencies using Poetry. This will set up a virtual environment and install the required packages specified in your `pyproject.toml` file:

    ```sh
    poetry install
    ```

2. **Run E2E Tests**: After installing the dependencies, you can run the E2E tests using the `pytest` framework. This command will execute all the E2E tests in your project:

    ```sh
    poetry run pytest
    ```

    You can also execute the tests from the level of root of project by running:

    ```sh
    make test-e2e
    ```

## Feature Files

In E2E testing, feature files define the test scenarios and steps that need to be executed. These files typically use a language like Gherkin to describe the behavior of the application in a human-readable format.

- **Location**: Feature files are located in a `tests/e2e/features`.
- **Structure**: Each feature file contains one or more scenarios, which outline a series of steps that represent a specific behavior or use case.
- **Writing Feature Files**: When writing feature files, follow the [Gherkin][gherkin-syntax] syntax to define `Feature`, `Scenario`, `Given`, `When`, and `Then` statements. These statements provide a clear, natural-language description of each test case.

## Dependency Management

Poetry is used to manage dependencies for your E2E tests. It ensures that your project has consistent and isolated environments for running the tests. 

- **Adding Dependencies**: If you need to add new dependencies for your E2E tests, use the `poetry add` command:

    ```sh
    poetry add <package-name>
    ```

- **Virtual Environment**: Poetry creates a virtual environment within your project, isolating the dependencies and ensuring compatibility.

- **Dependency Groups**: You can use groups to categorize dependencies (e.g., development and production). For example, to add a dependency to a specific group:

    ```sh
    poetry add <package-name> --group <group-name>
    ```

## Code Quality

### black, isort, mypy

- **isort** is a Python code formatting tool that automatically sorts and organizes import statements in your Python code. By using isort, you ensure that imports are consistently organized, which improves code readability and maintainability. It helps avoid confusion and errors that can arise from duplicate or misordered imports.

    ```sh
    poetry run isort .
    ```

- **black** is another Python code formatting tool that enforces a consistent coding style by reformatting Python code according to a set of predefined rules. By using black, you can maintain a consistent and standardized code style across your project, making it easier to read and collaborate with other developers. It can save time by automatically handling code formatting tasks.

    ```sh
    poetry run black .
    ```

- **mypy** is a static type-checking tool for Python code. It helps catch type errors and inconsistencies at compile-time, providing better code quality and preventing potential runtime errors. By using mypy, you can improve the robustness and reliability of your code.

    ```sh
    poetry run mypy .
    ```

<!-- Resources -->

[gherkin-syntax]: https://cucumber.io/docs/gherkin/reference/
