#!/bin/sh

# List all pull requests and checkout the branch

gh pr list | gum choose | cut -f1 | xargs gh pr checkout
