# PRE-EXISTING: Environment Issues (Not Migration-Caused)

These issues existed before the Centrifugo migration and block building/testing of some sub-projects.

---

## 1. Flutter — Dart SDK Version Incompatibility
- **Project constraint**: `>=2.11.0 <3.0.0` (pre-null-safety)
- **Installed SDK**: Dart 3.10.4 (null-safety only)
- **Error**: `The lower bound of "sdk: '>=2.11.0 <3.0.0'" must be 2.12.0 or higher to enable null safety.`
- **Fix**: Install Dart 2.x SDK, or migrate the entire app to null safety (major effort)

## 2. React Client Cabinet — node-sass Build Failure
- **Error**: `node-sass` uses `node-gyp` which requires Python ≤3.12
- **Installed**: Python 3.13 (incompatible with `node-gyp`)
- **Error**: `ModuleNotFoundError: No module named 'distutils'` (removed in Python 3.12)
- **Fix**: Use Python 3.11 or switch from `node-sass` to `sass` (Dart Sass)

## 3. React Admin — TypeScript Compilation Error
- **Error**: `node_modules/@types/minimatch/index.d.ts(29,48): error TS1005: ',' expected.`
- **Cause**: Type definition incompatibility with current TypeScript version
- **Fix**: Pin `@types/minimatch` to a compatible version or update

## 4. Go Exchange CLI — Integration Tests Need Running Database
- **Error**: `dial tcp: lookup db: no such host` in 1 test package
- **Affected**: Integration tests that connect to MySQL
- **Fix**: Run tests inside Docker with database service available
