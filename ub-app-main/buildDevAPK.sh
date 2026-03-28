#!/bin/bash
set -e

flutter pub get

version=$(grep 'version: ' pubspec.yaml | sed 's/version: //')

flutter clean

cp ./buildConfigs/mobile/flutter_native_splash.yaml ./

flutter pub run flutter_native_splash:create

rm -f ./flutter_native_splash.yaml

flutter build apk --release --split-per-abi --dart-define=ENV=DEV --dart-define=VERSION="$version"

rm -rf ./z_androidApps/*

cp  ./build/app/outputs/apk/release/* ./z_androidApps

# send files to telegram bot
echo "Sending 64bit file to telegram bot"

telegram-send --file ./z_androidApps/app-arm64-v8a-release.apk --caption "dev-64bit" --timeout 1000000000.0

echo "Sending 32bit file to telegram bot"

telegram-send --file ./z_androidApps/app-armeabi-v7a-release.apk --caption "dev-32bit" --timeout 1000000000.0