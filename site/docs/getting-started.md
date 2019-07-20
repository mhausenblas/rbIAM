# Getting Started

!!! warning
    This tool is heavy WIP and currently the CLI binaries are available only for macOS and Linux platforms. Please [raise an issue on GitHub](https://github.com/mhausenblas/rbIAM/issues?q=is%3Aissue+is%3Aopen+sort%3Aupdated-desc) if you're experiencing problems or something does not quite work like described in here.

This guide walks you through the set up and usage of `rbIAM`.

## Prerequisites

In order for you to use `rbIAM`, the following must be true:

- You have credentials for AWS configured.
- You have access to an EKS cluster or in general an Kubernetes-on-AWS cluster.
- You have `kubectl` [installed](https://kubernetes.io/docs/tasks/tools/install-kubectl/).

## Install

To install `rbIAM` execute the following two commands. First, download
the respective binary (here shown for macOS) like so:

```sh
curl -L https://github.com/mhausenblas/rbIAM/releases/latest/download/rbiam-macos -o /usr/local/bin/rbiam
```

And then make it executable:

```sh
chmod +x /usr/local/bin/rbiam
```

!!! tip
    For Linux install, simply replace the `-macos` part with `-linux`

## Usage

Walkthrough: