#!/bin/bash

# A basic cli for handling the staging with git
#
# This is a simple cli that utilizes the output from several git commands, and
# pipes them into `gum filter` and `gum choose` to provide a simple cli wrapper
# for `git add` and `git reset`.
#
# For a more feature complete version of this script check out
# [Nickiel12's Git-Gum](https://github.com/Nickiel12/Git-Gum).

git_add_filter(){
    input=$1
    if [[ $input == "M "* ]]; then
        if [ ! -z "$(git ls-files . --exclude-standard --others -m | grep "${input:2}")" ]; then
            echo ${input:2}
        fi
    else
        if [ ! -z "$(git ls-files . --exclude-standard --others -m | grep "${input:3}")" ]; then
            echo ${input:3}
        fi
    fi
}

choice=$(echo -e "add\nreset" | gum filter)

case $choice in
    add )
        export -f git_add_filter
        selection=`git status --short | xargs -I{} bash -c 'git_add_filter "{}"' | sed "$ a None"  | gum choose --no-limit`
        if [ "$selection" = "None" ]; then
            echo "No files selected"
        else
            git add -- "$selection"
        fi
        ;;
    reset )
        selection=`git diff --staged --name-only | sed "$ a None" | gum choose --no-limit`
        if [ "$selection" = "None" ]; then
            echo "None selected - skipping"
        else
            echo "$selection" | git reset 
        fi
        ;;
esac
