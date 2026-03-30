# BUG-002: Incorrect Import in RedisSubManager.php

## Severity: **LOW** — Unused import, does not cause runtime errors but is incorrect

## File
`ub-server-main/src/Exchange/CoreBundle/Services/RedisSubManager.php`

## Issue
Line 12 was changed from:
```php
use Exchange\CommunicationBundle\Services\EmqttPublishManager;
```
to:
```php
use Exchange\CommunicationBundle\Services\CentrifugoManager;
```

However, `CentrifugoManager` is **never used** in this file. The class uses `Exchange\ExternalExchangeBundle\Services\MqttPublishManager` (line 20), not `CentrifugoManager`.

The original `EmqttPublishManager` import was already unused — this was a pre-existing dead import that the migration blindly renamed.

## Fix
Remove the unused import on line 12:
```php
// DELETE: use Exchange\CommunicationBundle\Services\CentrifugoManager;
```

## Impact
- No runtime impact (PHP tolerates unused imports)
- Code quality issue — dead import reference
