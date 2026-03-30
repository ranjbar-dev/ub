# BUG-001: PHP Syntax Error in MqttPublishManager.php

## Severity: **CRITICAL** — Breaks PHP parsing, service will not load

## File
`ub-server-main/src/Exchange/ExternalExchangeBundle/Services/MqttPublishManager.php`

## Error
```
PHP Parse error: syntax error, unexpected token "*", expecting "function" on line 166
```

## Root Cause
During the EMQX cleanup, a commented-out method call was removed between lines 163-165, but the removal also deleted the `/**` opening of the docblock for the `publishTradesToTradeBookFromExternalExchange()` method. The orphaned `* @param` lines (lines 166-169) are now treated as PHP code rather than a comment block.

## Current (Broken)
```php
// line 162:
    }


     * @param $pairCurrencyName
     * @param $trades
     * //TODO unit test
     */
    public function publishTradesToTradeBookFromExternalExchange($pairCurrencyName, $trades): void
```

## Fix
Add back the `/**` docblock opener:
```php
    }

    /**
     * @param $pairCurrencyName
     * @param $trades
     */
    public function publishTradesToTradeBookFromExternalExchange($pairCurrencyName, $trades): void
```

## Impact
- The entire `MqttPublishManager` class fails to load
- All 5 event subscribers that depend on it will fail (TradeBookUpdated, PriceUpdated, OhlcUpdated, DepthUpdated, ChangePricePercentageUpdated)
- `RedisSubManager` will also fail to construct
- **All real-time data publishing from the PHP backend is broken**
