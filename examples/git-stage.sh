#!/bin/bash

ADD="Add"
RESET="Reset"

ACTION=$(gum choose "$ADD" "$RESET")

if [ "$ACTION" == "$ADD" ]; then
    git status --short | cut -c 4- | gum choose --no-limit | xargs git add
else
    git status --short | cut -c 4- | gum choose --no-limit | xargs git restore
fi
