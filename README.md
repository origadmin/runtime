# Runtime Package

## Introduction

The current go package defines the general configuration of the service runtime, as well as the loading of the runtime
configuration.

### Available Packages

- **[proto](proto)**: For compatibility with other languages, the interface is defined using proto files and implemented
  using gRPC. All proto definition files used by the Runtime are placed in 'proto' directory.
- **[config](config)**: The files in this directory define the basic configuration of the service runtime, as well as
  the loading of the run configuration.
- **[registry](registry)**: This directory defines an alias for 'kratos/v2/registry', primarily for backward
  compatibility and for placing import error paths.
- **[transport](transport)**: The current directory currently defines only the transport implementation of gins, which
  is not complete. You can use [protoc-gen-go-gins](../cmd/protoc-gen-go-gins) generates the relevant code.

## Getting Started

To incorporate the Toolkit into your project, follow these steps:

1. **Add the dependency**: Add the Toolkit as a dependency in your `go.mod` file, specifying the latest version:

```bash
go get github.com/origadmin/toolkit/runtime@vX.Y.Z

```

Replace `vX.Y.Z` with the desired version or `latest` to fetch the most recent release.

2. **Import required packages**: In your Go source files, import the necessary packages from the Toolkit:

```go
import (
    "github.com/origadmin/toolkit/runtime"
    "github.com/origadmin/toolkit/runtime/config"
    "github.com/origadmin/toolkit/runtime/registry"
)

// NewDiscovery creates a new discovery.
func NewDiscovery(registryConfig *config.RegistryConfig) registry.Discovery {
	if registryConfig == nil {
		panic("no registry config")
	}
	discovery, err := runtime.NewDiscovery(registryConfig)
	if err != nil {
		panic(err)
	}
	return discovery
}

// NewRegistrar creates a new registrar.
func NewRegistrar(registryConfig *config.RegistryConfig) registry.Registrar {
	if registryConfig == nil {
        panic("no registry config")
	}
	registrar, err := runtime.NewRegistrar(registryConfig)
	if err != nil {
        panic(err)
	}
	return registrar
}

```

## Contributing

We welcome contributions from the community to improve and expand the Toolkit. To contribute, please follow these
guidelines:

1. **Familiarize yourself with the project**: Read the [CONTRIBUTING] file for details on the contribution process, code
   style, and Pull Request requirements.
2. **Submit an issue or proposal**: If you encounter any bugs, have feature suggestions, or want to discuss potential
   changes, create an issue in the [GitHub repository](https://github.com/origadmin/toolkit).
3. **Create a Pull Request**: After implementing your changes, submit a Pull Request following the guidelines outlined
   in [CONTRIBUTING].

## Contributors

### Code of Conduct

All contributors and participants are expected to abide by the [Contributor Covenant][ContributorHomepage],
version [2.1][v2.1]. This document outlines the expected behavior when interacting with the Toolkit community.

## License

The Toolkit is distributed under the terms of the [MIT]. This permissive license allows for free use, modification, and
distribution of the toolkit in both commercial and non-commercial contexts.

[CONTRIBUTING]: CONTRIBUTING.md

[ContributorHomepage]: https://www.contributor-covenant.org

[v2.1]: https://www.contributor-covenant.org/version/2/1/code_of_conduct.html

[MIT]: LICENSE
