#!/bin/sh

# This script is used to write a conventional commit message.
# It prompts the user to choose the type of commit as specified in the
# conventional commit spec. And then prompts for the summary and detailed
# description of the message and uses the values provided. as the summary and
# details of the message.
#
# If you want to add a simpler version of this script to your dotfiles, use:
#
# alias gcm='git commit -m "$(gum input)" -m "$(gum write)"'

TYPE=$(gum choose "fix" "feat" "docs" "style" "refactor" "test" "chore" "revert")
SCOPE=$(gum input --placeholder "scope")

# Since the scope is optional, wrap it in parentheses if it has a value.
test -n "$SCOPE" && SCOPE="($SCOPE)"

# Pre-populate the input with the type(scope): so that the user may change it
SUMMARY=$(gum input --value "$TYPE$SCOPE: " --placeholder "Summary of this change")
DESCRIPTION=$(gum write --placeholder "Details of this change")

# Check if any staged files are present. If there are non, ask the user if 
# they want to commit all files
if [ -z "$(git status -s -uno | grep -v '^ ' | awk '{print $2}')" ]; then
    gum confirm "No staged files. Commit all?" && git add .
fi

# Commit these changes if user confirms
gum confirm "Commit changes?" && git commit -m "$SUMMARY" -m "$DESCRIPTION"
