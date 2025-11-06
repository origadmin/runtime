package config

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/testing/protocmp"

	appv1 "github.com/origadmin/runtime/api/gen/go/config/app/v1"
	discoveryv1 "github.com/origadmin/runtime/api/gen/go/config/discovery/v1"
	loggerv1 "github.com/origadmin/runtime/api/gen/go/config/logger/v1"
	middlewarev1 "github.com/origadmin/runtime/api/gen/go/config/middleware/v1"
	tracev1 "github.com/origadmin/runtime/api/gen/go/config/trace/v1"
	transportv1 "github.com/origadmin/runtime/api/gen/go/config/transport/v1"
	configproto "github.com/origadmin/runtime/test/integration/config/proto"
)

// AssertTestConfig performs a detailed, field-by-field assertion on the TestConfig struct.
func AssertTestConfig(t *testing.T, expected, actual *configproto.TestConfig) {
	require.NotNil(t, expected, "Expected config should not be nil")
	require.NotNil(t, actual, "Actual config should not be nil")

	AssertAppConfig(t, expected.App, actual.App)
	AssertServersConfig(t, expected.Servers, actual.Servers)
	AssertClientConfig(t, expected.Client, actual.Client)
	AssertLoggerConfig(t, expected.Logger, actual.Logger)
	AssertDiscoveriesConfig(t, expected.Discoveries, actual.Discoveries)
	AssertTraceConfig(t, expected.Trace, actual.Trace)
	AssertMiddlewaresConfig(t, expected.Middlewares, actual.Middlewares)

	require.Equal(t, expected.RegistrationDiscoveryName, actual.RegistrationDiscoveryName, "RegistrationDiscoveryName does not match")
}

// AssertAppConfig asserts that two App configurations are semantically equal.
func AssertAppConfig(t *testing.T, expected, actual *appv1.App) {
	if expected == nil {
		require.Nil(t, actual, "Actual App config should be nil, but was not")
		return
	}
	require.NotNil(t, actual, "Expected App config to be non-nil, but it was nil")
	if diff := cmp.Diff(expected, actual, protocmp.Transform()); diff != "" {
		require.Fail(t, "App configuration does not match", "Diff (-expected +actual):\n%s", diff)
	}
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
