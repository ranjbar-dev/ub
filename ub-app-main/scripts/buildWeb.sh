#!/bin/bash
set -e

flutter pub get

version=$(grep 'version: ' pubspec.yaml | sed 's/version: //')

cp ./buildConfigs/web/flutter_native_splash.yaml ./

rm -f ./web/manifest.json
rm -f ./web/index.html

cp ./buildConfigs/web/manifest.json ./web

flutter pub run flutter_native_splash:create

rm -f ./web/index.html
cp ./buildConfigs/web/index.html ./web

rm -f ./flutter_native_splash.yaml

flutter build web --web-renderer canvaskit --release --dart-define=ENV=PRODUCTION --dart-define=VERSION="$version"
