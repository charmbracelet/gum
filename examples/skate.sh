#!/bin/sh

# Building a simple `skate` TUI with gum to allow you to select a database and
# pick a value from skate.

DATABASE=$(skate list-dbs | gum choose)
skate list --keys-only "$DATABASE" | gum filter | xargs -I {} skate get {}"$DATABASE"
