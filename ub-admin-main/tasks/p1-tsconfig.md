# Task: Enable noImplicitAny in tsconfig

**ID:** p1-tsconfig  
**Phase:** 1 — Type Safety Foundation  
**Severity:** 🔴 CRITICAL  
**Dependencies:** None (but must be coordinated with p1-generic-response, p1-request-params)  

## Problem

`noImplicitAny: false` in `tsconfig.json` undermines the entire `"strict": true` setting, allowing implicit `any` types everywhere and defeating TypeScript's purpose.

## File to Modify

**`tsconfig.json`** (project root)

### Current Content
```json
{
  "compilerOptions": {
    "noImplicitAny": false,
    "target": "es5",
    "lib": ["dom", "dom.iterable", "esnext"],
    "allowJs": true,
    "skipLibCheck": true,
    "esModuleInterop": true,
    "allowSyntheticDefaultImports": true,
    "strict": true,
    "forceConsistentCasingInFileNames": true,
    "module": "esnext",
    "moduleResolution": "node",
    "resolveJsonModule": true,
    "isolatedModules": true,
    "noEmit": true,
    "jsx": "react",
    "baseUrl": "./src"
  },
  "include": ["src", "internals/startingTemplate/**/*"]
}
```

### Target Content
```json
{
  "compilerOptions": {
    "noImplicitAny": true,
    "target": "es5",
    "lib": ["dom", "dom.iterable", "esnext"],
    "allowJs": true,
    "skipLibCheck": true,
    "esModuleInterop": true,
    "allowSyntheticDefaultImports": true,
    "strict": true,
    "forceConsistentCasingInFileNames": true,
    "noUnusedLocals": true,
    "noUnusedParameters": true,
    "noImplicitReturns": true,
    "noFallthroughCasesInSwitch": true,
    "module": "esnext",
    "moduleResolution": "node",
    "resolveJsonModule": true,
    "isolatedModules": true,
    "noEmit": true,
    "jsx": "react",
    "baseUrl": "./src"
  },
  "include": ["src", "internals/startingTemplate/**/*"]
}
```

## Execution Steps

1. Run `npm run checkTs` to get the current baseline of type errors
2. Change `noImplicitAny` from `false` to `true`
3. Add `noUnusedLocals`, `noUnusedParameters`, `noImplicitReturns`, `noFallthroughCasesInSwitch`
4. Run `npm run checkTs` again — expect 100+ new errors
5. Fix errors file-by-file. Recommended order:
   - `src/services/constants.ts` (StandardResponse, RequestParameters)
   - `src/services/*.ts` (all service functions)
   - `src/utils/*.ts` (utility functions)
   - `src/app/containers/*/saga.ts` (saga action payloads)
   - `src/app/containers/*/types.ts` (state shapes)
   - `src/app/components/**/*.tsx` (component props)
6. For temporary unresolvable cases, use `// @ts-expect-error — TODO: type this` (NOT `// @ts-ignore`)
7. Run `npm run checkTs` to verify zero errors
8. Run `npm test` to verify no regressions

## Validation

```bash
npm run checkTs   # Must pass with 0 errors
npm test          # Must pass all existing tests
npm run build     # Must build successfully
```
