#!/usr/bin/env bash

# If the user passes '-h', '--help', or 'help' print out a little bit of help.
# text.
case "$1" in
	"-h" | "--help" | "help")
        printf 'Generate kaomojis on request.\n\n'
        printf 'Usage: %s [kind]\n' "$(basename "$0")"
        exit 1
        ;;
esac

# The user can pass an argument like "bear" or "angry" to specify the general
# kind of Kaomoji produced.
sentiment=""
if [[ $1 != "" ]]; then
	sentiment=" $1"
fi

# Ask mods to generate Kaomojis. Save the output in a variable.
kaomoji="$(mods "generate 10${sentiment} kaomojis. number them and put each one on its own line.")"
if [[ $kaomoji == "" ]]; then
	exit 1
fi

# Pipe mods output to gum so the user can choose the perfect kaomoji. Save that
# choice in a variable. Also note that we're using cut to drop the item number
# in front of the Kaomoji.
choice="$(echo "$kaomoji" | gum choose | cut -d ' ' -f 2)"
if [[ $choice == "" ]]; then
	exit 1
fi

# If xsel (X11) or pbcopy (macOS) exists, copy to the clipboard. If not, just
# print the Kaomoji.
if command -v xsel &> /dev/null; then
	printf '%s' "$choice" | xclip -sel clip # X11
elif command -v pbcopy &> /dev/null; then
	printf '%s' "$choice" | pbcopy # macOS
else
	# We can't copy, so just print it out.
	printf 'Here you go: %s\n' "$choice"
	exit 0
fi

# We're done!
printf 'Copied %s to the clipboard\n' "$choice"
