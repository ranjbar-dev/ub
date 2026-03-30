# BUG-004: React Shim Files Still Export Mqtt-Named Symbols

## Severity: **LOW** — Functional but retains MQTT naming

## Files
- `ub-client-cabinet-main/app/services/MqttService2.ts`
- `ub-client-cabinet-main/app/services/RegisteredMqttService.ts`

## Issue
These files re-export Centrifugo services under old MQTT names for backward compatibility:
```typescript
// MqttService2.ts
export { CentrifugoPublicService as MqttService, centrifugoPublicService as mqttService2 } from './CentrifugoPublicService';

// RegisteredMqttService.ts
export { CentrifugoAuthService as RegisteredMqttService, centrifugoAuthService as registeredMqttService } from './CentrifugoAuthService';
```

Currently **no files import from these shims** (verified via grep), making them dead code.

## Fix
Delete both files since nothing imports from them:
- `app/services/MqttService2.ts`
- `app/services/RegisteredMqttService.ts`

## Impact
- No functional impact (unused files)
- Files contain "Mqtt" in filenames and export names, contradicting the cleanup requirement
