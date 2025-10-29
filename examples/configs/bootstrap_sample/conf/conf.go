package conf

import (
	brokerpb "github.com/origadmin/runtime/api/gen/go/runtime/broker/v1"
	datapb "github.com/origadmin/runtime/api/gen/go/runtime/data/v1"
	discoverypb "github.com/origadmin/runtime/api/gen/go/runtime/discovery/v1"
	loggerpb "github.com/origadmin/runtime/api/gen/go/runtime/logger/v1"
	middlewarepb "github.com/origadmin/runtime/api/gen/go/runtime/middleware/v1"
	securitypb "github.com/origadmin/runtime/api/gen/go/runtime/security/v1"
	transportpb "github.com/origadmin/runtime/api/gen/go/runtime/transport/v1"
	websocketpb "github.com/origadmin/runtime/api/gen/go/runtime/websocket/v1"
)

// Bootstrap is the application bootstrap config.
type Bootstrap struct {
	Logger      *loggerpb.Logger
	Servers     *transportpb.Servers
	Data        *datapb.Data
	Security    *securitypb.Security
	Discoveries *discoverypb.Discoveries
	Middlewares *middlewarepb.Middlewares
	Broker      *brokerpb.Broker
	Websocket   *websocketpb.Websocket
}
