#!/usr/bin/bash
# Author: zenobit
# Description: Uses gum to provide a simple VMs runner for quickemu and quickget
# License MIT

_define_variables() {
	progname="${progname:="${0##*/}"}"
	version='0.1'
	vms=(*.conf)
	color=$(( RANDOM % 255 + 1 ))
	#export BORDER="rounded"
	export BORDERS_FOREGROUND="$color"
	export GUM_FILTER_INDICATOR_FOREGROUND="$color"
	export GUM_CHOOSE_CURSOR_FOREGROUND="$color"
	export GUM_CHOOSE_SELECTED_FOREGROUND="$color"
	export GUM_FILTER_SELECTED_PREFIX_FOREGROUND="$color"
	export GUM_FILTER_SELECTED_PREFIX_BORDER_FOREGROUND="$color"
	export GUM_FILTER_MATCH_FOREGROUND="$color"
	export GUM_FILTER_PROMPT_FOREGROUND="$color"

	if ! command -v quickemu >/dev/null 2>&1; then
		echo 'You are missing quickemu!'
	fi
	QUICKGET=$(command -v quickget)
	if ! command -v gum >/dev/null 2>&1; then
		echo 'You are missing gum! Exiting...' && exit 1
	fi
}

_generate_supported(){
	echo "Extracting OS Editions and Releases..."
	rm -rf /tmp/distros
	mkdir -p /tmp/distros
	"$QUICKGET" | awk 'NR==2,/zorin/' | cut -d':' -f2 | grep -o '[^ ]*' > /tmp/supported
	while read -r get_name; do
		supported=$($QUICKGET $get_name | sed 1d)
		echo ${get_name}
		echo "$supported"
		echo "$supported" > "/tmp/distros/${get_name}"
	done < /tmp/supported
}

_if_needed() {
	if [ ! -f /tmp/supported ]; then
		_generate_supported
	fi
}

show_vms() {
	if [ ${#vms[@]} -eq 0 ]; then
		echo 'No VMs found.'
	else
	echo "${vms[@]%.*}" | tr " " "\n"
	fi
}

show_headers() {
	distro=$(shuf -n 1 /tmp/supported)
	homepage=$("$QUICKGET" -s "${distro}")

	header_logo=$(gum style --padding "0 1" --border=rounded --border-foreground $color " Simple VMs runner
▄▄▄▄ ▄▄▄  ▄  ▄ ▄   ▄
█  █ █  █ █  █ █▄  █
█  █ █▄▀  █  █ █ ▀▄█
█▄▀▄ █ ▀▄ █▄▄█ █   █
v$version▀  for quickemu")

	tip_header=$(gum style --bold "Tip: ")
	tip_distro=$(gum style --align left "$distro")
	tip_temp=$(gum join --align top "$tip_header" "$tip_distro")
	tip_homepage=$(gum style --align left "$homepage")
	tip=$(gum join --vertical --align top "$tip_temp" "$tip_homepage")
	tip_border=$(gum style --padding "0 1" --border=rounded --border-foreground $color "$tip")

	vms_header=$(gum style --align center --bold "virtual machines:" && gum style $(echo "${vms[@]%.*}" | tr " " "\n"))
	vms_border=$(gum style --padding "0 1" --border=rounded --border-foreground $color "$vms_header")

	gum join --vertical --align top "$header_logo" "$tip_border" "$vms_border"
}

gum_choose_os() {
	os=$(gum filter < /tmp/supported)
	choices=$(cat "/tmp/distros/$os")
}

gum_choose_release() {
	release=$(echo "$choices" | grep 'Releases:' | cut -d':' -f2 | grep -o '[^ ]*' | gum filter --sort)
}

gum_choose_edition() {
	edition=$(echo "$choices" | grep 'Editions:' | cut -d':' -f2 | grep -o '[^ ]*' | gum filter --sort)
}

gum_choose_VM() {
	if ls | grep ".conf" ; then
		chosen=$(ls -1 | grep ".conf" | rev | cut -d'.' -f2- | rev | gum filter)
	else
		echo "No VMs to run."
	fi
}

create_VM() {
	gum_choose_os
	if [ -z "$os" ]; then exit 100
	elif [ "$(echo "$choices" | wc -l)" = 1 ]; then
		clear
		gum_choose_release
		clear
		"$QUICKGET" "$os" "$release"
	else
		clear
		gum_choose_release
		clear
		gum_choose_edition
		clear
		"$QUICKGET" "$os" "$release" "$edition"
	fi
	show_headers
}

run_VM() {
	quickemu -vm "$chosen.conf"
}

open_distro_homepage(){
	gum_choose_os
	"$QUICKGET" -o "${os}" >/dev/null 2>&1 &
}

# MENU
_show_menu() {
	while true
	do
	start=$(echo "create
run
homepage
EXIT $progname" | gum choose --selected run)
	case $start in
		create ) create_VM;;
		run ) gum_choose_VM && run_VM;;
		homepage ) open_distro_homepage;;
		"EXIT $progname" ) exit 0;;
	esac
	done
}

# run
clear
_define_variables
_if_needed
show_headers
_show_menu
