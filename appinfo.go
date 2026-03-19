package runtime

import (
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/goexts/generic/must"
	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	appv1 "github.com/origadmin/runtime/api/gen/go/config/app/v1"
	"github.com/origadmin/toolkits/identifier"
)

const (
	defaultProject = "origadmin"
	defaultAppName = "Unknown Service"
	defaultVersion = "1.0.0"
	defaultEnv     = "dev"
)

// NewAppInfo creates a new application information instance.
func NewAppInfo(name, version string) *appv1.App {
	ai := &appv1.App{
		Name:     name,
		Version:  version,
		Metadata: make(map[string]string),
	}
	AdjustAppInfo(ai)
	return ai
}

// NewAppInfoBuilder returns a new, blank App instance for building.
func NewAppInfoBuilder() *appv1.App {
	ai := &appv1.App{
		Metadata: make(map[string]string),
	}
	AdjustAppInfo(ai)
	return ai
}

// AdjustAppInfo adjusts the application info, setting default values.
func AdjustAppInfo(ai *appv1.App) {
	if ai == nil {
		return
	}

	if ai.Project == "" {
		ai.Project = defaultProject
	}

	if ai.Id == "" {
		ai.Id = uuid.New().String()
	}

	if ai.Version == "" {
		ai.Version = defaultVersion
	}

	if ai.Name == "" {
		ai.Name = defaultAppName
	}

	if ai.Env == "" {
		ai.Env = defaultEnv
	}

	if ai.Hostname == "" {
		ai.Hostname = ResolveHost()
	}

	if ai.Metadata == nil {
		ai.Metadata = make(map[string]string)
	}

	if ai.InstanceId == "" {
		ai.InstanceId = NewInstanceID(ai.Project, ai.Id, ai.Version, ai.Hostname)
	}

	if ai.StartTime == nil {
		ai.StartTime = timestamppb.Now()
	}
}

// NewInstanceID generates an instance ID.
func NewInstanceID(project, appID, version, host string) string {
	return fmt.Sprintf("%s.%s.%s@%s#%s", project, appID, version, host, must.Do(identifier.GenerateString()))
}

// ResolveHost returns the host identifier.
func ResolveHost() string {
	if v := os.Getenv("POD_NAME"); v != "" {
		return v
	}
	if v := os.Getenv("HOSTNAME"); v != "" {
		return v
	}
	if h, err := os.Hostname(); err == nil && h != "" {
		return h
	}
	if ip := firstNonLoopbackIP(); ip != "" {
		return ip
	}
	return "unknown-host"
}

// UpdateAppInfo merges application information from a source App into a destination App.
// Only non-empty fields from the source will overwrite the destination fields.
func UpdateAppInfo(dst, src *appv1.App) {
	if dst == nil || src == nil {
		return
	}

	if src.Project != "" {
		dst.Project = src.Project
	}
	if src.Id != "" {
		dst.Id = src.Id
	}
	if src.Name != "" {
		dst.Name = src.Name
	}
	if src.Version != "" {
		dst.Version = src.Version
	}
	if src.Env != "" {
		dst.Env = src.Env
	}
	if src.Hostname != "" {
		dst.Hostname = src.Hostname
	}
	if src.InstanceId != "" {
		dst.InstanceId = src.InstanceId
	}
	if src.Metadata != nil {
		if dst.Metadata == nil {
			dst.Metadata = make(map[string]string)
		}
		for k, v := range src.Metadata {
			dst.Metadata[k] = v
		}
	}
	if src.StartTime != nil {
		dst.StartTime = src.StartTime
	}

	// Always ensure the final state is valid
	AdjustAppInfo(dst)
}

// CloneAppInfo clones the application info.
func CloneAppInfo(src *appv1.App) *appv1.App {
	if src == nil {
		return nil
	}
	cloned := proto.Clone(src)
	if app, ok := cloned.(*appv1.App); ok {
		return app
	}
	return nil
}

// firstNonLoopbackIP returns the first non-loopback IPv4 address.
func firstNonLoopbackIP() string {
	ifaces, err := net.Interfaces()
	if err != nil {
		return ""
	}
	for _, iface := range ifaces {
		if (iface.Flags&net.FlagUp) == 0 || (iface.Flags&net.FlagLoopback) != 0 {
			continue
		}
		if isContainerLikeInterface(iface.Name) {
			continue
		}
		addrs, _ := iface.Addrs()
		for _, a := range addrs {
			var ip net.IP
			switch v := a.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue
			}
			if ip[0] == 169 && ip[1] == 254 {
				continue
			}
			if ip[0] == 172 && ip[1] == 17 {
				continue
			}
			return ip.String()
		}
	}
	return ""
}

// isContainerLikeInterface checks if the interface name is container-like.
func isContainerLikeInterface(name string) bool {
	n := strings.ToLower(name)
	return n == "docker0" ||
		strings.HasPrefix(n, "veth") ||
		strings.HasPrefix(n, "br-") ||
		strings.HasPrefix(n, "cni0") ||
		strings.HasPrefix(n, "flannel") ||
		strings.HasPrefix(n, "weave") ||
		strings.HasPrefix(n, "virbr") ||
		strings.HasPrefix(n, "cbr0")
}
