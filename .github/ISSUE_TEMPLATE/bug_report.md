---
name: "🐛 Bug Report"
description: Create a report to help us improve
title: "(short issue description)"
labels: [bug]
assignees: []
body:
  - type: textarea
    id: description
    attributes:
      label: Describe the bug
      description: What is the problem? A clear and concise description of the bug.
    validations:
      required: true

  - type: textarea
    id: current
    attributes:
      label: Current behavior
      description: |
        Please include full errors, uncaught exceptions, screenshots, and relevant logs.
        Using environment variable JFROG_CLI_LOG_LEVEL="DEBUG" upon running the command will provide more log information.
    validations:
      required: true

  - type: textarea
    id: reproduction
    attributes:
      label: Reproduction steps
      description: |
        Provide steps to reproduce the behavior.
    validations:
      required: false

  - type: textarea
    id: expected
    attributes:
      label: Expected behavior
      description: |
        What did you expect to happen?
    validations:
      required: false

  - type: input
    id: client-go-version
    attributes:
      label: JFrog Client-Go version
    validations:
      required: true

  - type: input
    id: cli-version
    attributes:
      label: JFrog CLI version (if applicable)
      description: using "jf --version"
    validations:
      required: false

  - type: input
    id: os-version
    attributes:
      label: Operating system type and version
    validations:
      required: true

  - type: input
    id: rt-version
    attributes:
      label: JFrog Artifactory version
    validations:
      required: false

  - type: input
    id: xr-version
    attributes:
      label: JFrog Xray version
    validations:
      required: false
