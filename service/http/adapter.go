package http

import (
	"time"

	transhttp "github.com/go-kratos/kratos/v2/transport/http"

	configv1 "github.com/origadmin/runtime/api/gen/go/config/v1"
	"github.com/origadmin/runtime/service/selector"
	"github.com/origadmin/runtime/service/tls"
	tkerrors "github.com/origadmin/toolkits/errors"
)

// adaptClientConfig Option to convert service configuration to a specific protocol
func adaptClientConfig(cfg *configv1.Service) ([]transhttp.ClientOption, error) {
	// 1. Validate the configuration
	if cfg == nil {
		return nil, tkerrors.Errorf("client config is required for creation")
	}

	httpCfg := cfg.GetHttp()
	if httpCfg == nil {
		return nil, tkerrors.Errorf("http client config is required for creation")
	}

	var opts []transhttp.ClientOption

	// 2. Processing endpoints
	if endpoint := httpCfg.GetEndpoint(); endpoint != "" {
		opts = append(opts, transhttp.WithEndpoint(endpoint))
	}

	// 3. Processing timeout
	if timeout := httpCfg.GetTimeout(); timeout > 0 {
		opts = append(opts, transhttp.WithTimeout(time.Duration(timeout)*time.Millisecond))
	}

	// 4. Handle TLS
	if tlsCfg := httpCfg.GetTlsConfig(); tlsCfg != nil {
		tlsConfig, err := tls.NewClientTLSConfig(tlsCfg)
		if err != nil {
			return nil, tkerrors.Wrapf(err, "invalid TLS config for client creation")
		}
		opts = append(opts, transhttp.WithTLSConfig(tlsConfig))
	}

	// 5. Process selectors
	if selectorCfg := cfg.GetSelector(); selectorCfg != nil {
		filter, err := selector.NewFilter(selectorCfg)
		if err != nil {
			return nil, tkerrors.Wrapf(err, "invalid selector config for client creation")
		}
		opts = append(opts, transhttp.WithNodeFilter(filter))
	}

	return opts, nil
}
