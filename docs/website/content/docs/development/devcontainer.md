---
title: Dev Container
weight: 2
---

A dev container is a lightweight, portable development environment defined by a `devcontainer.json` file inside the .devcontainer folder, typically used with containerization technologies like Docker.
It allows developers to quickly set up their environment, with containerization technology being the only prerequisite, as the container defines the necessary tools, dependencies, and settings for 
a consistent development environment across any platform.

## Visual Studio Code

Get started with [vsccode](https://code.visualstudio.com/docs/devcontainers/containers).

To debug from the container:
1. Build the debug recipya: `task build-debug`
2. Start recipya: `./bin/recipya_debug serve`
3. Press F5 to start the VS Code debugger
4. Select the `recipya_debug` process from the list


