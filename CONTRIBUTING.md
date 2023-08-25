# Contributing to Go POS

First and foremost, thank you for your interest in contributing to Go POS project! We appreciate your time and effort in helping us improve the project. To ensure a smooth and collaborative development process, please read the following guidelines before contributing to the project.

## Setting up the development environment

Before you start contributing to the project, you will need to set up your development environment. To get started, you should have the following software installed:

- [Go](https://golang.org/) 1.21 or higher
- [Task](https://taskfile.dev/)
- [Docker](https://www.docker.com/) or [Podman](https://podman.io/)
- [Docker Compose](https://docs.docker.com/compose/) or [Podman Compose](https://github.com/containers/podman-compose)
- [PostgreSQL](https://hub.docker.com/_/postgres) container

You should also have [Git](https://git-scm.com/) installed to clone the repository and submit merge requests.

To get started with the project, you can follow these steps:

1. Fork this repository.
2. Clone your forked repository to your local machine.
3. Install the project dependencies: `task install`
4. Create a copy of the `.env.example` file and rename it to `.env`. Update configuration values as needed.
5. Run PostgreSQL database container: `task db:up && task db:create && task migrate:up`
6. Start the development server: `task dev`

## Submitting bug reports

If you encounter any issues with the project, please submit a bug report. To do so, please follow these guidelines:

1. Check the existing issues to see if your bug has already been reported.
2. If your bug has not been reported, create a new issue with a descriptive title and detailed steps to reproduce the bug.
3. Include any relevant error messages or screenshots.
4. Assign the `bug` label to the issue.

## Making feature requests

If you have an idea for a new feature, you can submit a feature request. To do so, please follow these guidelines:

1. Check the existing issues and merge requests to see if your feature request has already been submitted.
2. If your feature request has not been submitted, create a new issue with a descriptive title and detailed description of the feature.
3. Assign the `enhancement` label to the issue.

## Submitting code changes

If you would like to submit a code change, please follow these guidelines:

1. Create a new branch for your changes: `git checkout -b {feat,fix,refactor,docs,ci,chore,test}/my-new-branch-name`
2. Make your changes and once you are satisfied with your changes, commit them with a descriptive message using the [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0) format.
3. Push your changes to your fork: `git push origin {feat,fix,refactor,docs,ci,chore,test}/my-new-branch-name`
4. Create a merge request against the `main` branch of the original repository.
5. Assign the appropriate labels to the merge request.

## Code standards

When submitting code changes, please adhere to the following standards:

1. **Formatting and linting**: Format your code using `gofmt` and lint your code using `golangci-lint`. Make sure your code passes all linting checks before submitting your changes. If you are on VS Code, you can install the [Go](https://marketplace.visualstudio.com/items?itemName=golang.Go) extension then use `golangci-lint` as a linter.
2. **Packages and imports**: Use goimports to organize imports and remove unused imports.
3. **Naming conventions**: Follow the Go naming conventions for variables, functions, and types. Use descriptive and meaningful names for variables, functions, and types. Avoid using abbreviations or acronyms unless they are well-known and widely used. Refer to the [Go Naming Slides](https://go.dev/talks/2014/names.slide) for more information.
4. **Comments**: Include comments to explain complex logic or provide clarity where needed. Use complete sentences and follow the Go commenting style (e.g., `//` for single-line comments and `/** ... */` for multi-line comments).

## Conclusion

Thank you for your interest in contributing to Go POS project. We appreciate your time and effort in helping us improve the project. We look forward to your contributions and hope to see you in the community soon!
