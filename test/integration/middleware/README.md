# Middleware Integration Test

This directory contains integration tests for middleware configuration and instance conversion.

## Test Files

- `middleware_test.go`: Main test file containing test cases for middleware configuration to instance conversion.
- `helper.go`: Helper functions for loading YAML configuration files.
- `configs/config.yaml`: Test configuration file for middlewares.

## Test Cases

1. **TestMiddlewareConfigToInstance**:
   - Tests loading of middleware configuration from YAML file.
   - Verifies creation of middleware instances using NewClient and NewServer methods.
   - Checks middleware chain construction using BuildClient method.
   - Validates selector middleware functionality.
   - Tests handling of disabled and unknown type middlewares.

2. **TestBuildServerMiddlewares**:
   - Tests middleware chain construction for server-side using BuildServer method.

## Running Tests

To run the tests, execute the following command from the project root directory:

```bash
go test ./runtime/test/integration/middleware/... -v
```

## Middleware Configuration

The test configuration includes the following middlewares:

1. **Metadata Middleware**:
   - Adds custom metadata headers to requests/responses.
   - Configured with prefix "x-origadmin" and sample data.

2. **Logging Middleware**:
   - Enables logging of requests and responses.
   - Configured with info log level and text format.

3. **Selector Middleware**:
   - Applies specific middlewares to requests matching a pattern.
   - Configured to match "/api/v1/users" path and apply metadata middleware.