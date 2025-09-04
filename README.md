# Runtime Package

## Introduction

The `@/runtime` package serves as the foundational implementation for building **distributed-ready services** within our system, **leveraging the Kratos microservice framework**. It provides a set of core abstractions and utilities that **simplify the development and management of Kratos-based services**, enabling developers to write business logic independent of the underlying deployment model.

This design ensures that services developed using this package are inherently **loosely coupled** and **independently deployable**. Whether initially deployed as part of a monolithic application or as standalone microservices, the transition between these deployment modes can be achieved seamlessly by simply changing the underlying infrastructure implementations, without requiring modifications to the core business code.

The `Runtime` component centralizes the management of essential resources required for service operation, including configuration, logging, monitoring, caching, and data storage, **all integrated within the Kratos ecosystem.**

## Core Philosophy / Design Principles

1.  **Kratos-Native Development**: The `@/runtime` package is built upon the Kratos framework, providing a streamlined and opinionated way to develop Kratos-based microservices. It aims to reduce boilerplate and enforce best practices within the Kratos ecosystem.
2.  **Deployment Agnosticism**: Business logic should be decoupled from deployment concerns. By leveraging Kratos's robust capabilities, the `@/runtime` package provides abstractions that allow services to run equally well in a single process (monolithic deployment) or across multiple nodes (distributed deployment).
3.  **Loose Coupling by Design**: Through standardized interfaces and explicit communication mechanisms (e.g., API contracts, gRPC), the package, in conjunction with Kratos's service discovery and communication patterns, encourages and facilitates the development of services that are inherently loosely coupled, making them easier to test, maintain, and evolve independently.
4.  **Infrastructure Simplification & Swappability**: It abstracts away the complexities of underlying infrastructure components (e.g., configuration sources, logging backends, service discovery mechanisms). This allows for easy swapping of different infrastructure implementations (e.g., from local file-based config to a distributed config center) as the system scales, without impacting the service's core logic, **all while integrating seamlessly with Kratos's extension points.**
5.  **Standardization & Consistency**: The package provides a consistent set of patterns and tools for common cross-cutting concerns, ensuring a standardized development and operational experience across all services built on this architecture, **adhering to Kratos's architectural principles.**

## Naming Conventions

To ensure clarity and consistency within the `@/runtime` package, especially given its close integration with the Kratos framework, we adhere to the following naming conventions:

1.  **Kratos-Specific Elements**: Any type, interface, or struct that is directly from the Kratos framework, or is a direct wrapper/alias/bridge implementation of a Kratos element within `@/runtime`, **must be prefixed with `K`**.
    *   **Purpose**: This clearly indicates that the element is Kratos-native or directly interacts with Kratos.
    *   **Examples**: `KRegistrar`, `KDiscovery`, `KApp`, `KServer`.

2.  **@/runtime Abstractions**: Any type, interface, or struct defined within `@/runtime` that represents a general-purpose abstraction or utility, even if its underlying implementation uses Kratos, **should NOT be prefixed with `K`**.
    *   **Purpose**: These are the public-facing APIs of the `@/runtime` package, designed to be used by business logic without direct Kratos dependency at the interface level.
    *   **Examples**: `Sender` (for mail), `StorageProvider`, `ConfigLoader`.

This distinction helps developers quickly understand the nature and scope of each component, promoting better code readability and maintainability.

### Before You Start

Before you start using the Runtime package, ensure that you have the following prerequisites:
In order to prevent import conflicts caused by packages with the same name as `kratos`, packages with the same name in
this database will import the export content from `kratos`.
All type definitions will be prefixed with the `K` fixed prefix.
Note: Only `type` declarations are prefixed, not functions.

### Available Packages

- **[bootstrap](bootstrap)**: The bootstrap package contains Configuration file reading and writing, initialization
  variable declaration, etc
- **[cmd](cmd)**: Contains command-line utilities or example main packages for the runtime.
- **[config](config)**: The files in this directory define the basic configuration of the service runtime, as well as
  the loading of the run configuration.
- **[context](context)**: The context directory defines the context interface and the context implementation.
- **[gen](gen)**: This directory contains **generated code** (e.g., Go structs, gRPC service stubs) derived from the Protocol Buffer (`.proto`) definition files.
- **[internal](internal)**: Contains internal packages and helper utilities not intended for external consumption.
- **[log](log)**: Provides logging interfaces and implementations for the runtime, integrated with Kratos's logging system.
- **[mail](mail)**: The mail directory defines the email interface and the email implementation.
- **[middleware](middleware)**: The middleware directory defines the middleware interface and the middleware **implementation for common cross-cutting concerns like authentication, logging, and rate limiting.**
- **[proto](proto)**: This directory contains the **original Protocol Buffer (`.proto`) definition files** that define the service interfaces and data structures. These files are used to generate code for various languages, ensuring compatibility and strong typing across different services.
- **[registry](registry)**: This directory defines an alias for 'kratos/v2/registry', primarily for backward
  compatibility and for placing import error paths.
- **[service](service)**: The service directory contains the definition of the service interface, which is used to
  define the interface of the service and the implementation of the service.
- **[storage](storage)**: This directory provides **abstractions and implementations for various data storage mechanisms**, including caching, databases, and file storage. It centralizes data access concerns, allowing services to easily swap underlying storage technologies.
- **[third_party](third_party)**: Contains vendored or third-party code dependencies.

### Top-Level Files

- **`runtime.go`**: The main entry point and core logic for the runtime package, orchestrating the initialization and lifecycle of services.
- **`generate.go`**: Defines `go generate` commands used for automated code generation tasks, such as generating protobuf code.
- **`buf.lock`**: A lock file generated by Buf, ensuring reproducible builds of Protocol Buffers.
- **`buf.yaml`**: The main configuration file for Buf, defining linting rules, formatting, and other settings for `.proto` files.
- **`buf.gen.yaml`**: Configuration for Buf's code generation, specifying how `.proto` files are compiled into code for different languages (e.g., Go).

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
