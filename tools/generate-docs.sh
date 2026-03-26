#!/usr/bin/env bash
set -euo pipefail

# This script prepares examples for tfplugindocs by concatenating
# a full, runnable example module (variables, outputs, etc.) into a single
# file that the tool can parse.

ORIG_EXAMPLES_DIR="./examples"
TMP_EXAMPLES_DIR="./examples_tmp"
GEN_EXAMPLES_DIR="./examples" # The tool looks for this directory by default

# Cleanup function to be called on script exit
cleanup() {
  echo "Cleaning up..."
  # Check if the temp dir exists before trying to move it
  if [ -d "$TMP_EXAMPLES_DIR" ]; then
    rm -rf "$GEN_EXAMPLES_DIR"
    mv "$TMP_EXAMPLES_DIR" "$ORIG_EXAMPLES_DIR"
  fi
}

# 1. Rename original examples to a temporary location.
#    If the temp dir exists, a previous run may have failed. Clean it up.
if [ -d "$TMP_EXAMPLES_DIR" ]; then
    rm -rf "$TMP_EXAMPLES_DIR"
fi
mv "$ORIG_EXAMPLES_DIR" "$TMP_EXAMPLES_DIR"

# Ensure cleanup happens on script exit or interruption
trap cleanup EXIT

# 2. Create the directory structure tfplugindocs expects
mkdir -p "$GEN_EXAMPLES_DIR/resources"
mkdir -p "$GEN_EXAMPLES_DIR/data-sources"

# 3. Process resources
echo "Preparing resource examples..."
for r_dir in "$TMP_EXAMPLES_DIR"/resources/*; do
  if [ ! -d "$r_dir" ]; then continue; fi

  resource_name=$(basename "$r_dir")
  target_dir="$GEN_EXAMPLES_DIR/resources/$resource_name"
  target_file="$target_dir/resource.tf"
  mkdir -p "$target_dir"

  # Find all the leaf directories containing a main.tf, which represent a runnable example
  find "$r_dir" -type f -name "main.tf" -exec dirname {} \; | sort | while read -r example_sub_dir; do
    # Add a comment header to delineate multiple examples
    sub_dir_path=${example_sub_dir#$r_dir/}
    if [ -n "$sub_dir_path" ] && [ "$sub_dir_path" != "." ]; then
        echo "# ------ Example: ${sub_dir_path} ------" >> "$target_file"
    fi

    # Concatenate all .tf files that make up the example, in order
    for tf_file in versions.tf variables.tf main.tf outputs.tf; do
        if [ -f "$example_sub_dir/$tf_file" ]; then
            cat "$example_sub_dir/$tf_file" >> "$target_file"
            echo -e "
" >> "$target_file"
        fi
    done
  done
done

# 4. Process data sources
echo "Preparing data source examples..."
for d_dir in "$TMP_EXAMPLES_DIR"/data-sources/*; do
  if [ ! -d "$d_dir" ]; then continue; fi

  ds_name=$(basename "$d_dir")
  target_dir="$GEN_EXAMPLES_DIR/data-sources/$ds_name"
  target_file="$target_dir/data-source.tf"
  mkdir -p "$target_dir"

  # Concatenate all .tf files that make up the example, in order
  for tf_file in versions.tf variables.tf main.tf outputs.tf; do
      if [ -f "$d_dir/$tf_file" ]; then
          cat "$d_dir/$tf_file" >> "$target_file"
          echo -e "
" >> "$target_file"
      fi
  done
done

# 5. Run the documentation generator
echo "Running tfplugindocs..."
go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs

echo "Documentation generation complete."
