# Testing Documentation

## Overview

This document describes the testing strategy and setup for the iban.im project.

## Test Structure

The project uses Go's built-in testing framework with the following test organization:

### Unit Tests

1. **Model Tests** (`model/`)
   - `user_test.go`: Tests for User model
     - Password hashing and comparison
     - Struct field validation
   - `iban_test.go`: Tests for IBAN model
     - Password hashing and comparison
     - Handle validation and uniqueness checking
     - GORM callbacks (BeforeSave, Validate)
   - `group_test.go`: Tests for Group model
     - Struct field validation
     - IBAN associations

2. **Config Tests** (`config/`)
   - `config_test.go`: Tests for configuration management
     - Environment variable parsing
     - Default values
     - Singleton pattern

3. **Resolver Tests** (`resolvers/`)
   - `sign_up_test.go`: SignUp mutation tests
   - `sign_in_test.go`: SignIn/authentication tests  
   - `change_password_test.go`: Password change tests
   - `change_profile_test.go`: Profile update tests
   - `delete_profile_test.go`: Profile deletion tests
   - `iban_test.go`: IBAN CRUD operation tests

## Test Helpers

### `resolvers/test_helpers.go`

Provides common test utilities:

- `setupTestDB(t)`: Creates an in-memory SQLite database
- `createTestUser(t, db, ...)`: Creates a test user with hashed password
- `createTestIban(t, db, ...)`: Creates a test IBAN
- `contextWithUserID(userID)`: Creates authenticated context
- `setupTestResolverWithDB(t)`: Complete test setup with cleanup

## Running Tests

### Run all tests
```bash
go test ./...
```

### Run tests for a specific package
```bash
go test ./model
go test ./resolvers
go test ./config
```

### Run tests with verbose output
```bash
go test ./... -v
```

### Run specific test
```bash
go test ./resolvers -run TestSignUp
```

### Run tests with coverage
```bash
go test ./... -cover
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## Test Database

Tests use in-memory SQLite databases for speed and isolation:
- Each test gets a fresh database instance
- No external database dependencies required
- Automatic cleanup after each test
- GORM auto-migration ensures schema is up-to-date

## Test Patterns

### Table-Driven Tests

Most tests use table-driven patterns for comprehensive coverage:

```go
tests := []struct {
    name          string
    args          someArgs
    expectSuccess bool
    expectError   string
}{
    {
        name: "Success case",
        args: someArgs{...},
        expectSuccess: true,
    },
    {
        name: "Error case",
        args: someArgs{...},
        expectSuccess: false,
        expectError: "Expected error message",
    },
}

for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        // Test implementation
    })
}
```

### Testing Authenticated Endpoints

Use `contextWithUserID()` to create authenticated contexts:

```go
ctx := contextWithUserID(int(userID))
resp, err := resolver.SomeAuthenticatedMethod(ctx, args)
```

## Coverage Goals

- **Models**: Focus on business logic (password hashing, validations, callbacks)
- **Resolvers**: Test all mutations and queries with success/error cases
- **Config**: Test environment variable parsing and defaults
- **Edge Cases**: Empty values, nil pointers, boundary conditions
- **Authentication**: Authenticated and unauthenticated access

## Known Limitations

1. Some tests in `utils/` require a running server and are not unit tests
2. Schema tests depend on file system structure
3. GraphQL integration tests are not included (would require full server setup)

## Future Improvements

- [ ] Add integration tests with real database
- [ ] Add API-level tests using httptest
- [ ] Add performance/load tests
- [ ] Improve test coverage for handler package
- [ ] Add mutation testing
- [ ] Add GraphQL query/mutation integration tests

## Test Maintenance

- Run tests before committing changes
- Update tests when changing business logic
- Add tests for new features before implementation (TDD)
- Keep test data minimal and focused
- Use meaningful test names that describe what is being tested
