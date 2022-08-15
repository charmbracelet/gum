#! /bin/sh

# This script is used to manage git branches such as delete, update, and rebase
# them. It prompts the user to choose the branches and the action they want to
# perform.
#
# For an explanation on the script and tutorial on how to create it, watch:
# https://www.youtube.com/watch?v=tnikefEuArQ

GIT_COLOR="#f14e32"

git_color_text () {
  gum style --foreground "$GIT_COLOR" "$1"
}

get_branches () {
  if [ ${1+x} ]; then
    gum choose --selected.foreground="$GIT_COLOR" --limit="$1" $(git branch --format="%(refname:short)")
  else
    gum choose --selected.foreground="$GIT_COLOR" --no-limit $(git branch --format="%(refname:short)")
  fi
}

git rev-parse --git-dir > /dev/null 2>&1

if [ $? -ne 0 ];
then
  echo "$(git_color_text "!!") Must be run in a $(git_color_text "git") repo" 
  exit 1
fi

gum style \
  --border normal \
  --margin "1" \
  --padding "1" \
  --border-foreground "$GIT_COLOR" \
  "$(git_color_text ' Git') Branch Manager"

echo "Choose $(git_color_text 'branches') to operate on:"
branches=$(get_branches)

echo ""
echo "Choose a $(git_color_text "command"):"
command=$(gum choose --cursor.foreground="$GIT_COLOR" rebase delete update)
echo ""

echo $branches | tr " " "\n" | while read -r branch
do
  case $command in
    rebase)
      base_branch=$(get_branches 1)
      git fetch origin
      git checkout "$branch"
      git rebase "origin/$base_branch"
      ;;
    delete)
      git branch -D "$branch"
      ;;
    update)
      git checkout "$branch"
      git pull --ff-only
      ;;
  esac
done
