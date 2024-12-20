# Runtime Package

## Introduction

The current go package defines the general configuration of the service runtime, as well as the loading of the runtime
configuration.

The `Runtime` controls the resources required when a project is started, including configuration files, logs,
monitoring,
caches, and databases.

### Before You Start

Before you start using the Runtime package, ensure that you have the following prerequisites:
In order to prevent import conflicts caused by packages with the same name as `kratos`, packages with the same name in
this database will import the export content from `kratos`.
All type definitions will be prefixed with the `K` fixed prefix.
Note: Only `type` declarations are prefixed, not functions.

### Available Packages

- **[bootstrap](bootstrap)**: The bootstrap package contains Configuration file reading and writing, initialization
  variable declaration, etc
- **[config](config)**: The files in this directory define the basic configuration of the service runtime, as well as
  the loading of the run configuration.
- **[context](context)**: The context directory defines the context interface and the context implementation.
- **[data](data)**: The data directory defines the data interface, caching, database, and other storage implementation.
- **[gen](gen)**: The protobuf directory contains the definition of the protobuf protocol.
- **[proto](proto)**: For compatibility with other languages, the interface is defined using proto files and implemented
  using gRPC. All proto definition files used by the Runtime are placed in 'proto' directory.
- **[mail](mail)**: The mail directory defines the email interface and the email implementation.
- **[middleware](middleware)**: The middleware directory defines the middleware interface and the middleware
-
- **[registry](registry)**: This directory defines an alias for 'kratos/v2/registry', primarily for backward
  compatibility and for placing import error paths.
- **[service](service)**: The service directory contains the definition of the service interface, which is used to
  define the interface of the service and the implementation of the service.

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
