# Task: Standardize Import Ordering

**ID:** p5-import-ordering  
**Phase:** 5 — Code Organization  
**Severity:** 🟢 LOW  
**Dependencies:** None  

## Problem

Import statements have no consistent ordering. Some files mix node_modules, absolute paths, and relative paths randomly.

## Target Convention

```typescript
// 1. React & framework
import React, { useEffect, useState } from 'react';
import { useDispatch, useSelector } from 'react-redux';

// 2. Third-party libraries
import { PayloadAction } from '@reduxjs/toolkit';
import { call, put, takeLatest } from 'redux-saga/effects';

// 3. Absolute imports (services, utils, store)
import { apiService } from 'services/api_service';
import { MessageService, MessageNames } from 'services/message_service';
import { StandardResponse } from 'services/constants';

// 4. Relative imports (same container/feature)
import { actions } from './slice';
import { selectUserData } from './selectors';
import { UserAccountsState } from './types';
```

Separate each group with a blank line.

## Implementation: ESLint Rule

Add to `.eslintrc` (or create/update):

```json
{
  "rules": {
    "import/order": [
      "error",
      {
        "groups": [
          ["builtin", "external"],
          "internal",
          ["parent", "sibling", "index"]
        ],
        "newlines-between": "always",
        "alphabetize": { "order": "asc", "caseInsensitive": true }
      }
    ]
  }
}
```

Then run autofix:
```bash
npx eslint --fix --rule 'import/order: error' src/
```

## Execution Steps

1. Add the ESLint import/order rule to config
2. Run `npx eslint --fix src/` to auto-sort imports
3. Manually review any files that couldn't be auto-fixed
4. Run `npm run checkTs` and `npm test`

## Validation

```bash
npm run lint      # Must pass with no import order warnings
npm run checkTs   # Must pass
npm test          # Must pass
```
