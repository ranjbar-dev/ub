# Task: Document All MessageService Events

**ID:** p6-document-messages  
**Phase:** 6 — State Management Improvement  
**Severity:** 🔴 CRITICAL  
**Dependencies:** None  

## Problem

`MessageService` has 67 event names in the `MessageNames` enum. There is zero documentation about:
- What data each event carries
- Who sends each event (which saga)
- Who listens to each event (which component's useEffect)
- Whether the event is one-to-one or one-to-many

This makes the pub/sub system completely opaque to AI agents.

## File to Modify

**`src/services/message_service.ts`**

### Current Enum (lines 4–69)
```typescript
/*eslint-disable*/
enum MessageNames {
  SET_INPUT_ERROR = 'SET_INPUT_ERROR',
  SETLOADING = 'SETLOADING',
  SET_ROW_LOADING = 'SET_ROW_LOADING',
  SET_BUTTON_LOADING = 'SET_BUTTON_LOADING',
  CLOSE_REJECT_POPUP = 'CLOSE_REJECT_POPUP',
  CLOSE_POPUP = 'CLOSE_POPUP',
  // ... 61 more entries
  ALLPY_PARAMS_TO_GRID = 'ALLPY_PARAMS_TO_GRID',
}
```

### Target: Add JSDoc to Each Enum Value

```typescript
enum MessageNames {
  /**
   * Sets validation errors on form inputs.
   * @payload value: Record<string, string[]> — field name → error messages
   * @sender Any saga on 422 API response
   * @listener Form components
   */
  SET_INPUT_ERROR = 'SET_INPUT_ERROR',

  /**
   * Shows/hides the global loading spinner.
   * @payload value: boolean — true = show, false = hide
   * @payload loadingId?: string — optional specific loader ID
   * @sender All sagas (before/after API calls)
   * @listener App root layout, individual containers
   */
  SETLOADING = 'SETLOADING',

  /**
   * Shows/hides a row-level loading indicator in a grid.
   * @payload value: boolean
   * @payload rowId: number — the specific row ID
   * @payload gridName: GridNames — which grid to target
   * @sender Action sagas (approve, reject, etc.)
   * @listener SimpleGrid component
   */
  SET_ROW_LOADING = 'SET_ROW_LOADING',

  // ... document ALL 67 events with the same pattern
}
```

### Also Document BroadcastMessage

```typescript
/**
 * Standard message shape for all MessageService events.
 *
 * @property name - The event type from MessageNames enum
 * @property value - Primary payload (type varies by event)
 * @property payload - Alternative payload field (some events use this instead of value)
 * @property additional - Extra data for complex events
 * @property errorId - Identifies which input field has an error
 * @property userId - User context for user-specific events
 * @property loadingId - Identifies which loading indicator to control
 * @property rowId - Grid row identifier for row-level operations
 * @property gridName - Which AG Grid instance to target
 * @property type - Window type or toast severity
 * @property child - Child data for nested events
 */
interface BroadcastMessage {
  name: MessageNames;
  value?: any;
  payload?: any;
  // ...
}
```

### Also Fix Typos in Event Names

```
ALLPY_PARAMS_TO_GRID → should be APPLY_PARAMS_TO_GRID (typo: ALLPY)
SET_COMMITIONS_DATA → should be SET_COMMISSIONS_DATA (typo: COMMITIONS)
```

⚠️ Fixing the typo requires updating ALL senders and listeners. Search:
```bash
grep -rn "ALLPY_PARAMS_TO_GRID" src/
grep -rn "SET_COMMITIONS_DATA" src/
```

## Execution Steps

1. Remove the `/*eslint-disable*/` from line 3
2. Add JSDoc to each MessageNames enum value
3. For each event, search the codebase to identify:
   - Senders: `grep -rn "MessageNames.EVENT_NAME" src/ | grep "send"`
   - Listeners: `grep -rn "MessageNames.EVENT_NAME" src/ | grep -v "send"`
4. Document BroadcastMessage interface
5. Fix the ALLPY and COMMITIONS typos (search & replace across codebase)
6. Run `npm run checkTs` and `npm test`

## Validation

```bash
npm run checkTs   # Must pass
npm test          # Must pass
```
