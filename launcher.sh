#!/bin/bash

# CHMOD All files (.sh)
find . -type f -iname "*.sh" -exec chmod +x {} \;

# Main choice function
function choice () {
        MAI=$(gum input --placeholder "[G]um, [Q]uit")

        if [ "$MAI" = "G" ]; then
                ./gum
        fi
        if [ "$MAI" = "Q" ]; then
                exit
        fi
}

choice
