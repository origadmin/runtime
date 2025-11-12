package declarative

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-kratos/kratos/v2/transport"
	kratoshttp "github.com/go-kratos/kratos/v2/transport/http"
	"google.golang.org/grpc/metadata"

	"github.com/origadmin/runtime/interfaces/security/declarative"
)

// HTTPHeaderProvider adapts an http.Header to the ValueProvider interface.
type HTTPHeaderProvider struct {
	Header http.Header
}

func (p HTTPHeaderProvider) Values(key string) []string {
	return p.Header[key]
}

func (p HTTPHeaderProvider) Get(key string) string {
	return p.Header.Get(key)
}

func (p HTTPHeaderProvider) GetAll() map[string][]string {
	return p.Header
}

// GRPCCredentialsProvider adapts a grpc.Metadata to the ValueProvider interface.
type GRPCCredentialsProvider struct {
	MD metadata.MD
}

func (p GRPCCredentialsProvider) Values(key string) []string {
	return p.MD.Get(key)
}

func (p GRPCCredentialsProvider) Get(key string) string {
	return p.MD.Get(key)[0]
}

func (p GRPCCredentialsProvider) GetAll() map[string][]string {
	return p.MD
}

// FromHTTPHeader 从 http.Header 创建 HTTPHeaderProvider
func FromHTTPHeader(header http.Header) *HTTPHeaderProvider {
	return &HTTPHeaderProvider{Header: header}
}

// FromHTTPRequest 从 http.Request 创建 HTTPHeaderProvider
func FromHTTPRequest(r *http.Request) *HTTPHeaderProvider {
	return &HTTPHeaderProvider{Header: r.Header}
}

// FromGRPCMetadata 从 metadata.MD 创建 GRPCCredentialsProvider
func FromGRPCMetadata(md metadata.MD) *GRPCCredentialsProvider {
	return &GRPCCredentialsProvider{MD: md}
}

// FromGRPCContext 从 context.Context 创建 GRPCCredentialsProvider
// 如果上下文中没有元数据，则返回错误
func FromGRPCContext(ctx context.Context) (*GRPCCredentialsProvider, error) {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		return &GRPCCredentialsProvider{MD: md}, nil
	}
	return nil, errors.New("no metadata found in context")
}

func FromServerContext(ctx context.Context) (declarative.ValueProvider, error) {
	if tr, ok := transport.FromServerContext(ctx); ok {
		switch tr.Kind() {
		case transport.KindHTTP:
			if ht, ok := tr.(kratoshttp.Transporter); ok {
				return FromHTTPRequest(ht.Request()), nil
			}
			return nil, errors.New("invalid HTTP transport type")
		case transport.KindGRPC:
			md, _ := metadata.FromIncomingContext(ctx)
			return FromGRPCMetadata(md), nil
		default:
			return nil, fmt.Errorf("unsupported transport type: %v", tr.Kind())
		}
	}
	return nil, errors.New("no transport found in context")
}
