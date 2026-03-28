# Task 006: Expand Validation Error Messages

## Priority: MEDIUM
## Risk: LOW (improves error messages — no logic change)
## Estimated Scope: 1 file

---

## Problem

`internal/api/handler/error.go:21` has a TODO:
```go
//todo we should write this more comprehensive to include all tags
```

The `customizeValidationError` function only handles 4 of 20+ validator tags:
- `required`
- `oneof`
- `gt`
- default fallback

Missing tags that are likely used in the codebase (based on go-playground/validator):
- `email` — email format
- `min` / `max` — string length or numeric bounds
- `len` — exact length
- `gte` / `lte` — greater/less than or equal
- `numeric` — numeric string
- `alphanum` — alphanumeric
- `uuid` — UUID format
- `url` — URL format
- `eqfield` — field equality (password confirmation)
- `nefield` — field inequality

## Goal

Expand the validation tag handling to produce user-friendly error messages for all tags used in the codebase.

## Implementation Plan

### Step 1: Find all validation tags used in the codebase

```bash
grep -rn 'validate:"' internal/ --include="*.go" -o | sort -u
```

This will show every struct tag like `validate:"required,email,min=6"`. Collect the unique tag names.

### Step 2: Expand `customizeValidationError` in `internal/api/handler/error.go`

Current function (around line 21):
```go
func customizeValidationError(errs validator.ValidationErrors) []customError {
    //todo we should write this more comprehensive to include all tags
    var customErrors []customError
    for _, err := range errs {
        switch err.Tag() {
        case "required":
            customErrors = append(customErrors, customError{
                field:   err.Field(),
                message: fmt.Sprintf("%s is required", err.Field()),
            })
        case "oneof":
            customErrors = append(customErrors, customError{
                field:   err.Field(),
                message: fmt.Sprintf("%s must be one of: %s", err.Field(), err.Param()),
            })
        case "gt":
            customErrors = append(customErrors, customError{
                field:   err.Field(),
                message: fmt.Sprintf("%s must be greater than %s", err.Field(), err.Param()),
            })
        default:
            customErrors = append(customErrors, customError{
                field:   err.Field(),
                message: fmt.Sprintf("%s is not valid", err.Field()),
            })
        }
    }
    return customErrors
}
```

**Add these cases before the `default`:**

```go
case "email":
    customErrors = append(customErrors, customError{
        field:   err.Field(),
        message: fmt.Sprintf("%s must be a valid email address", err.Field()),
    })
case "min":
    customErrors = append(customErrors, customError{
        field:   err.Field(),
        message: fmt.Sprintf("%s must be at least %s characters", err.Field(), err.Param()),
    })
case "max":
    customErrors = append(customErrors, customError{
        field:   err.Field(),
        message: fmt.Sprintf("%s must be at most %s characters", err.Field(), err.Param()),
    })
case "len":
    customErrors = append(customErrors, customError{
        field:   err.Field(),
        message: fmt.Sprintf("%s must be exactly %s characters", err.Field(), err.Param()),
    })
case "gte":
    customErrors = append(customErrors, customError{
        field:   err.Field(),
        message: fmt.Sprintf("%s must be greater than or equal to %s", err.Field(), err.Param()),
    })
case "lte":
    customErrors = append(customErrors, customError{
        field:   err.Field(),
        message: fmt.Sprintf("%s must be less than or equal to %s", err.Field(), err.Param()),
    })
case "numeric":
    customErrors = append(customErrors, customError{
        field:   err.Field(),
        message: fmt.Sprintf("%s must be a number", err.Field()),
    })
case "alphanum":
    customErrors = append(customErrors, customError{
        field:   err.Field(),
        message: fmt.Sprintf("%s must contain only letters and numbers", err.Field()),
    })
case "uuid":
    customErrors = append(customErrors, customError{
        field:   err.Field(),
        message: fmt.Sprintf("%s must be a valid UUID", err.Field()),
    })
case "url":
    customErrors = append(customErrors, customError{
        field:   err.Field(),
        message: fmt.Sprintf("%s must be a valid URL", err.Field()),
    })
case "eqfield":
    customErrors = append(customErrors, customError{
        field:   err.Field(),
        message: fmt.Sprintf("%s must match %s", err.Field(), err.Param()),
    })
```

### Step 3: Remove the TODO comment

Delete the line:
```go
//todo we should write this more comprehensive to include all tags
```

## Verification

```bash
go build ./...
# Optional: write a test in error_test.go that validates each tag produces the expected message
```

## Notes

- The `err.Param()` method returns the tag parameter (e.g., for `min=6`, Param() returns `"6"`)
- The `err.Field()` returns the struct field name
- Only add cases for tags actually used in the codebase (grep first)
