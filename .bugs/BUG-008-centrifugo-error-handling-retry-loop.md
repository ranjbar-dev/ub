# BUG-008: PHP CentrifugoManager Retry Logic Can Throw Unhandled Exception

## Severity: **LOW** — Edge case in error handling

## File
`ub-server-main/src/Exchange/CommunicationBundle/Services/CentrifugoManager.php`

## Issue
The `publishMessage()` method has a retry-once pattern:
```php
private function publishMessage(string $channel, array $data): void
{
    if ($this->parameterBag->get('kernel.environment') == 'test') {
        return;
    }

    try {
        $this->centrifugoClient->publish($channel, $data);
    } catch (\Exception $ex) {
        // Retry once on failure
        $this->centrifugoClient->publish($channel, $data);
    }
}
```

If the first `publish()` throws an exception and the retry also throws, the second exception is **not caught** and will propagate up the call stack. This could crash event subscribers processing orders/trades.

## Fix
Wrap the retry in a try-catch too, or log and swallow:
```php
try {
    $this->centrifugoClient->publish($channel, $data);
} catch (\Exception $ex) {
    try {
        $this->centrifugoClient->publish($channel, $data);
    } catch (\Exception $retryEx) {
        // Log error, don't crash the caller
    }
}
```

## Impact
- If Centrifugo is temporarily unreachable, trade/order processing may fail
- With EMQX, the old client likely had similar behavior, but HTTP API failures are more likely than MQTT TCP failures
