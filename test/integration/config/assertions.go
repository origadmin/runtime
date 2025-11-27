package config

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/testing/protocmp"

	"github.com/origadmin/runtime"
	discoveryv1 "github.com/origadmin/runtime/api/gen/go/config/discovery/v1"
	loggerv1 "github.com/origadmin/runtime/api/gen/go/config/logger/v1"
	middlewarev1 "github.com/origadmin/runtime/api/gen/go/config/middleware/v1"
	tracev1 "github.com/origadmin/runtime/api/gen/go/config/trace/v1"
	transportv1 "github.com/origadmin/runtime/api/gen/go/config/transport/v1"
	"github.com/origadmin/runtime/interfaces"
	configproto "github.com/origadmin/runtime/test/integration/config/proto"
)

// isNilConcreteValue checks if an interface holds a nil concrete value.
// This handles the Go gotcha where a non-nil interface can still wrap a nil pointer.
func isNilConcreteValue(i interface{}) bool {
	if i == nil {
		return true // Interface itself is nil
	}
	val := reflect.ValueOf(i)
	// Check if it's a pointer and if that pointer is nil
	return val.Kind() == reflect.Ptr && val.IsNil()
}

// AssertTestConfig performs a detailed, field-by-field assertion on the TestConfig struct.
func AssertTestConfig(t *testing.T, expected, actual *configproto.TestConfig) {
	require.NotNil(t, expected, "Expected config should not be nil")
	require.NotNil(t, actual, "Actual config should not be nil")

	AssertAppConfig(t, runtime.ConvertToAppInfo(expected.App), runtime.ConvertToAppInfo(actual.App))
	AssertServersConfig(t, expected.Servers, actual.Servers)
	AssertClientConfig(t, expected.Client, actual.Client)
	AssertLoggerConfig(t, expected.Logger, actual.Logger)
	AssertDiscoveriesConfig(t, expected.Discoveries, actual.Discoveries)
	AssertTraceConfig(t, expected.Trace, actual.Trace)
	AssertMiddlewaresConfig(t, expected.Middlewares, actual.Middlewares)

	require.Equal(t, expected.RegistrationDiscoveryName, actual.RegistrationDiscoveryName, "RegistrationDiscoveryName does not match")
}

// AssertAppConfig asserts that two App configurations are semantically equal by comparing their interface methods.
func AssertAppConfig(t *testing.T, expected, actual interfaces.AppInfo) {
	expectedIsNil := isNilConcreteValue(expected)
	actualIsNil := isNilConcreteValue(actual)

	if expectedIsNil {
		require.True(t, actualIsNil, "Actual App config should be nil, but was not")
		return
	}
	require.False(t, actualIsNil, "Expected App config to be non-nil, but it was nil")

	// Perform a field-by-field comparison using the interface's getter methods.
	// This is the correct way to test an interface's implementation.
	require.Equal(t, expected.ID(), actual.ID(), "App ID does not match")
	require.Equal(t, expected.Name(), actual.Name(), "App Name does not match")
	require.Equal(t, expected.Version(), actual.Version(), "App Version does not match")
	require.Equal(t, expected.Metadata(), actual.Metadata(), "App Metadata does not match")
}

// AssertServersConfig asserts that two Servers configurations are semantically equal.
func AssertServersConfig(t *testing.T, expected, actual *transportv1.Servers) {
	if expected == nil {
		require.Nil(t, actual, "Actual Servers config should be nil, but was not")
		return
	}
	require.NotNil(t, actual, "Expected Servers config to be non-nil, but it was nil")
	if diff := cmp.Diff(expected, actual, protocmp.Transform()); diff != "" {
		require.Fail(t, "Servers configuration does not match", "Diff (-expected +actual):\n%s", diff)
	}
}

// AssertClientConfig asserts that two Client configurations are semantically equal.
func AssertClientConfig(t *testing.T, expected, actual *transportv1.Client) {
	if expected == nil {
		require.Nil(t, actual, "Actual Client config should be nil, but was not")
		return
	}
	require.NotNil(t, actual, "Expected Client config to be non-nil, but it was nil")
	if diff := cmp.Diff(expected, actual, protocmp.Transform()); diff != "" {
		require.Fail(t, "Client configuration does not match", "Diff (-expected +actual):\n%s", diff)
	}
}

// AssertLoggerConfig asserts that two Logger configurations are semantically equal.
func AssertLoggerConfig(t *testing.T, expected, actual *loggerv1.Logger) {
	if expected == nil {
		require.Nil(t, actual, "Actual Logger config should be nil, but was not")
		return
	}
	require.NotNil(t, actual, "Expected Logger config to be non-nil, but it was nil")
	if diff := cmp.Diff(expected, actual, protocmp.Transform()); diff != "" {
		require.Fail(t, "Logger configuration does not match", "Diff (-expected +actual):\n%s", diff)
	}
}

// AssertDiscoveriesConfig asserts that two Discoveries configurations are semantically equal.
func AssertDiscoveriesConfig(t *testing.T, expected, actual *discoveryv1.Discoveries) {
	if expected == nil {
		require.Nil(t, actual, "Actual Discoveries config should be nil, but was not")
		return
	}
	require.NotNil(t, actual, "Expected Discoveries config to be non-nil, but it was nil")
	if diff := cmp.Diff(expected, actual, protocmp.Transform()); diff != "" {
		require.Fail(t, "Discoveries configuration does not match", "Diff (-expected +actual):\n%s", diff)
	}
}

// AssertTraceConfig asserts that two Trace configurations are semantically equal.
func AssertTraceConfig(t *testing.T, expected, actual *tracev1.Trace) {
	if expected == nil {
		require.Nil(t, actual, "Actual Trace config should be nil, but was not")
		return
	}
	require.NotNil(t, actual, "Expected Trace config to be non-nil, but it was nil")
	if diff := cmp.Diff(expected, actual, protocmp.Transform()); diff != "" {
		require.Fail(t, "Trace configuration does not match", "Diff (-expected +actual):\n%s", diff)
	}
}

// AssertMiddlewaresConfig asserts that two Middlewares configurations are semantically equal.
func AssertMiddlewaresConfig(t *testing.T, expected, actual *middlewarev1.Middlewares) {
	if expected == nil {
		require.Nil(t, actual, "Actual Middlewares config should be nil, but was not")
		return
	}
	require.NotNil(t, actual, "Expected Middlewares config to be non-nil, but it was nil")
	if diff := cmp.Diff(expected, actual, protocmp.Transform()); diff != "" {
		require.Fail(t, "Middlewares configuration does not match", "Diff (-expected +actual):\n%s", diff)
	}
}
