# IBAN Public Route Feature

This document describes the public IBAN route feature that allows accessing IBAN addresses via clean, user-friendly URLs.

## Overview

The application now supports serving IBAN addresses at the route pattern `/:userHandle/:ibanHandle`, making it easy to share and access IBAN information through memorable links.

## Usage

### Accessing IBAN via URL

Navigate to:
```
https://your-domain.com/username/ibanhandle
```

For example:
```
https://iban.im/testuser/primary
https://iban.im/testuser/savings
```

### Response Formats

#### HTML (Default)
By default, the route returns a user-friendly HTML page displaying:
- User information (name and handle)
- IBAN handle/alias
- Description (if available)
- The full IBAN number
- A "Copy to Clipboard" button with visual feedback

#### JSON
To get a JSON response, use one of these methods:

1. **Accept Header**:
```bash
curl -H "Accept: application/json" https://iban.im/username/ibanhandle
```

2. **Query Parameter**:
```bash
curl "https://iban.im/username/ibanhandle?format=json"
```

JSON response format:
```json
{
  "userHandle": "testuser",
  "ibanHandle": "primary",
  "iban": "TR320010009999901234567890",
  "description": "My primary bank account",
  "firstName": "Test",
  "lastName": "User"
}
```

## Privacy & Security

- **Only public IBANs are accessible** via this route
- Private IBANs (marked with `is_private = true`) will return a 404 error
- No authentication is required for accessing public IBANs
- Invalid user handles or IBAN handles return user-friendly error messages

## Error Handling

### User Not Found
If the user handle doesn't exist, the response will be:
- **HTML**: Error page with "User not found" message
- **JSON**: `{"error": "User not found"}` with 404 status

### IBAN Not Found or Private
If the IBAN handle doesn't exist or is marked as private:
- **HTML**: Error page with "IBAN not found or is private" message
- **JSON**: `{"error": "IBAN not found or is private"}` with 404 status

## Implementation Details

### Handler
The main handler is `RenderIbanPage` in `handler/iban_handler.go`:
- Validates user handle exists
- Validates IBAN handle exists and is public
- Supports content negotiation for HTML/JSON
- Uses Gin framework for routing

### Templates
- `templates/iban.tmpl.html`: Main IBAN display page
- `templates/error.tmpl.html`: Error display page
- Both templates use Tailwind CSS for styling

### Route Registration
The route is registered in `main.go`:
```go
router.GET("/:userHandle/:ibanHandle", handler.RenderIbanPage)
```

## Testing

### Creating Test Data
Use the seed script to create test data:
```bash
go run ./tools/seed/main.go
```

This creates a test user with several IBANs including:
- `testuser/primary` - Public IBAN
- `testuser/savings` - Public IBAN
- `testuser/private` - Private IBAN (not accessible)

### Running Tests
```bash
go test ./handler -v
```

Tests cover:
- Valid IBAN requests (HTML and JSON)
- Private IBAN access (should return 404)
- Non-existent user/IBAN handling
- Route validation logic

## Examples

### Curl Examples

```bash
# Get IBAN as HTML
curl http://localhost:8080/testuser/primary

# Get IBAN as JSON
curl -H "Accept: application/json" http://localhost:8080/testuser/primary

# Get IBAN as JSON (alternative)
curl "http://localhost:8080/testuser/primary?format=json"

# Try to access private IBAN (will fail)
curl http://localhost:8080/testuser/private
```

### Expected Behavior

1. **Valid Public IBAN**: Returns IBAN details (HTML or JSON)
2. **Private IBAN**: Returns 404 with error message
3. **Non-existent User**: Returns 404 with "User not found"
4. **Non-existent IBAN**: Returns 404 with "IBAN not found or is private"

## Future Enhancements

Possible improvements for the future:
- QR code generation for IBAN
- Share buttons (Twitter, LinkedIn, etc.)
- Analytics/tracking for IBAN views
- Custom themes per user
- Rate limiting to prevent abuse
