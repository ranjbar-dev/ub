# BUG-007: Centrifugo Publish Data Format Mismatch Between PHP and Go

## Severity: **MEDIUM** — May cause client-side parsing errors

## Issue
The PHP `CentrifugoManager` publishes data as PHP arrays (converted to JSON objects by phpcent):
```php
$this->centrifugoClient->publish($channel, ['price' => '100', 'amount' => '5']);
// Centrifugo receives: {"price": "100", "amount": "5"}
```

The Go `CentrifugoManager` publishes raw `[]byte` payloads that are first unmarshaled then re-marshaled:
```go
func (m *centrifugoManager) PublishTrades(ctx context.Context, pairName string, payload []byte) {
    data := m.unmarshalPayload(payload)  // json.Unmarshal into interface{}
    m.client.Publish(channel, data)       // re-marshals to JSON
}
```

## Potential Issues
1. **Double-encoding**: If the Go payload is already JSON-encoded bytes and gets marshaled again, clients may receive escaped JSON strings instead of objects
2. **Format inconsistency**: PHP sends structured arrays; Go sends pre-serialized payloads — the Centrifugo `data` field format may differ between the two publishers
3. **Client parsing**: Frontend clients (React/Flutter) may expect a specific format and break when receiving data from the other publisher

## Verification Needed
- Compare actual Centrifugo message payloads from PHP vs Go publishers
- Verify client-side parsing handles both formats
- Check if `unmarshalPayload` properly converts `[]byte` → `interface{}` without double-encoding

## Impact
- Clients may receive malformed JSON from one publisher but not the other
- Cannot be fully verified without running the full stack
