# BUG-009: Flutter centrifuge-dart Version May Be Incompatible with Pre-Null-Safety

## Severity: **HIGH** — Flutter app cannot build with current Dart SDK

## File
`ub-app-main/pubspec.yaml`

## Issue
The Flutter app requires Dart SDK `>=2.11.0 <3.0.0` (pre-null-safety), but:

1. The installed Dart SDK is 3.10.4 (null-safety only) — `flutter pub get` fails completely
2. The `centrifuge` package version needs to be compatible with pre-null-safety Dart (<3.0)
3. centrifuge-dart `^0.8.0` was selected for pre-null-safety compat, but this cannot be verified since `flutter pub get` fails

## Error
```
The lower bound of "sdk: '>=2.11.0 <3.0.0'" must be 2.12.0 or higher to enable null safety.
The current Dart SDK (3.10.4) only supports null safety.
```

## Pre-existing vs Migration
This is a **pre-existing environment issue** (Dart SDK too new for this project), NOT caused by the migration. However, the migration changed the dependency from `mqtt_client` to `centrifuge`, and correctness of this dependency swap cannot be verified until the SDK issue is resolved.

## Verification Needed
- Install Dart SDK 2.x to run `flutter pub get` and verify centrifuge resolves
- Check if centrifuge-dart 0.8.x actually supports pre-null-safety Dart
- Verify Centrifugo service imports and usage compile correctly

## Impact
- Flutter app build is blocked (pre-existing)
- centrifuge dependency correctness cannot be verified
