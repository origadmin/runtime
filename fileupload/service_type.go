package fileupload

// ServiceType represents the type of the service.
//
//go:generate stringer -type=ServiceType -output=service_type_string.go -trimprefix=ServiceType .
type ServiceType int

const (
	ServiceTypeGRPC    ServiceType = iota
	ServiceTypeHTTP    ServiceType = iota
	ServiceTypeUnknown ServiceType = iota
)
