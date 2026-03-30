# BUG-003: Class Name `MqttPublishManager` Still Contains "Mqtt"

## Severity: **MEDIUM** — Not a runtime bug, but contradicts the "remove all MQTT references" requirement

## File
`ub-server-main/src/Exchange/ExternalExchangeBundle/Services/MqttPublishManager.php`

## Issue
The class `MqttPublishManager` was not renamed during the migration. While its internals now use `CentrifugoManager` for publishing, the class name, file name, and all references still say "Mqtt".

## References (13 files affected by rename)
- `MqttPublishManager.php` — class definition + self-reference on line 157
- `RedisSubManager.php` — import (line 20), type hint (line 51), constructor param (line 57), docblock (line 29)
- `TradeBookUpdatedSubscriber.php` — import, type hint, constructor
- `PriceUpdatedSubscriber.php` — import, type hint, constructor
- `OhlcUpdatedSubscriber.php` — import, type hint, constructor
- `DepthUpdatedSubscriber.php` — import, type hint, constructor
- `ChangePricePercentageUpdatedSubscriber.php` — import, type hint, constructor
- `RedisSubManagerTest.php` — import + 5 mock references

## Fix
Rename class to `CentrifugoPublishManager` (file + class + all references). Since Symfony uses autowiring, service definitions will auto-resolve.

## Impact
- User requirement: "remove all MQTT references" — this class name violates that
- No functional bug — the class works correctly with its current Centrifugo internals
