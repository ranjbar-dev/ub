#!/bin/bash
set -e

# Commit and tag this change.
version=$(grep 'version: ' pubspec.yaml | sed 's/version: //')
git commit -m "Bump version to $version" pubspec.yaml
git tag "$version"
git push