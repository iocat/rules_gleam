#!/bin/bash

# Find the location of the runfiles directory
RUNFILES_DIR="$0.runfiles/_main"

# Use the runfiles location to call the helper script
cd "$RUNFILES_DIR" && "$RUNFILES_DIR/{main_js}" $@