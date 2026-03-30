#!/bin/bash
set -e

flutter pub get

version=$(grep 'version: ' pubspec.yaml | sed 's/version: //')

flutter clean

cp ./buildConfigs/mobile/flutter_native_splash.yaml ./

flutter pub run flutter_native_splash:create

rm -f ./flutter_native_splash.yaml

flutter build appbundle --dart-define=ENV=PRODUCTION --dart-define=VERSION="$version"
