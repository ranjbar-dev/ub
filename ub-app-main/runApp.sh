#!/bin/bash
set -e

version=$(grep 'version: ' pubspec.yaml | sed 's/version: //')

flutter run --dart-define=ENV=DEV --dart-define=VERSION="$version"