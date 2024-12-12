#!/bin/sh
# generates completions for each shell
set -e
rm -rf completions
mkdir completions
for sh in bash zsh fish; do
	go run main.go completion "$sh" >"completions/repocheck.$sh"
done
