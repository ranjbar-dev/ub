#!/bin/bash
set -e

version=$(grep 'version: ' pubspec.yaml | sed 's/version: //')

flutter run -d chrome --web-renderer canvaskit --web-port 8000 --dart-define=ENV=DEV --dart-define=VERSION="$version"