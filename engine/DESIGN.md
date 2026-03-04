# Runtime Engine Design (v34 - 2026-03-04)

## 1. Core Philosophy: Option-Based IoC Bootstrapping
The Runtime Engine is a business-agnostic IoC container that decouples **Capability Declaration** (init phase) from **Instance Configuration** (creation phase). The container is now "self-configuring" via a unified Options pattern.

## 2. Decoupled Lifecycle (The Timeline)

### 2.1 Capability Accumulation (`init` phase)
Components register their intent to be produced during Go's `init()` phase.
- **Global Pool**: `engine.Register(Category, Provider, ...opts)` stores registration metadata in a private global shadow pool.
- **Zero Side-Effect**: This phase only records metadata; no instantiation or logic is executed.

### 2.2 Container Bootstrapping (`NewContainer` phase)
When a new container instance is created, its initial state is defined via `RegistryOption`.
- **`WithResolver(res)`**: Injects the "Global Resolver" which defines how to route configuration sources to component categories.
- **`WithGlobalRegistrations()`**: Instructs the container to fetch a snapshot of the global shadow pool and load those capabilities into its private registry.
- **Benefits**: Ensures testability by allowing "clean" containers and provides an atomic transition from static metadata to an active registry.

### 2.3 Configuration Injection (`Load` phase)
Configuration data is injected via `Load(ctx, source, ...opts)` after the application has loaded its config files.
- **Directional Routing**: Supports targeting specific categories or instances.
- **Cascading Extraction**: Uses local `WithExtractor` if present, falling back to the injected Global Resolver.

## 3. Namespacing & Reserved Symbols

Identifiers starting with an underscore (`_`) are **Reserved for System Use**.
- **`_global`**: The system fallback scope for cascading lookups.
- **Policy**: The engine forbids users from registering any `Category` or `Scope` starting with `_`.

## 4. Context & Lifecycle (App Level)

The `App` instance manages the base lifecycle context following a strict chain:
1. **Default**: `context.Background()`.
2. **Override**: Replaced if `WithContext(ctx)` option is provided.
3. **Encapsulation**: Automatically wrapped with `context.WithCancel()` to ensure the application can signal a graceful shutdown.

## 5. Metadata Types (Contract-Driven)

All core types (`Scope`, `Category`, `Priority`, `Registration`) are defined in `runtime/contracts/component` to ensure engine purity and prevent circular dependencies.

## 6. Project Structure

```
runtime/
├── contracts/component/ # Core Types, Interfaces, and Registration struct.
├── engine/
│   ├── container/       # Core IoC logic, Shadow pool, and Perspective isolation.
│   ├── helpers.go       # Casting and retrieval shortcuts.
│   └── engine.go        # Unified NewContainer and Option definitions.
└── types.go             # Application-level aliases and Option mapping.
```
