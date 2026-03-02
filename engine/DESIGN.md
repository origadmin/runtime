# Runtime Engine Design (v30 - 2026-03-02)

## 1. Core Philosophy: Interface-Driven Chained Container
The Runtime Engine serves as the core component manager for `runtime.App`. It bridges the gap between global business configurations and functional module factories (Registry, Middleware, Database, etc.) using a zero-reflection, order-sensitive, and scope-aware architecture.

## 2. Key Principles

### 2.1 Interface-Based Config Sniffing
To remain decoupled from specific business configuration structures (e.g., `*conf.Config`), the engine relies on **interface assertions** during the extraction phase. 
- **Extractor Logic**: Instead of casting to a concrete struct, the `Extractor` asserts the root configuration against standard interfaces like `RegistryConfigGetter` or `MiddlewareConfigGetter`.
- **Decoupling**: This allows the business layer to define its own configuration layout while the engine remains generic.

### 2.2 Navigation-Based Handle
The `Handle` interface uses a "filesystem-like" navigation pattern:
- **Category-First**: Users navigate into a category context using `h.In(Category)`.
- **Scope-Optional**: Environment isolation (Server/Client) is an optional modifier via `WithScope`.
- **Chained Retrieval**: `h.In(Category).Get(ctx, name)` provides a clean, semantic way to access dependencies.

### 2.3 Ordered Block Protocol
Engine strictly follows the ordering defined in the configuration (especially critical for Middleware chains):
- **Normalization**: Handled by `configutil` logic (`Active -> Default -> Ordered Configs`).
- **Sequential Init**: The `Init` phase iterates through categories by **Priority** and through instances by their original sequence.

## 3. Architecture Layers

### 3.1 Registry Interface (Behavior Definition)
```go
type Registry interface {
    Handle
    // Register defines the factory behavior and extraction logic for a category.
    Register(c Category, e Extractor, p Provider, opts ...RegisterOption)
    // BindRoot injects the global business configuration object.
    BindRoot(root any)
    // Init performs sorted, sequential initialization of all registered components.
    Init(ctx context.Context) error
}
```

### 3.2 Handle Interface (Context Provider)
```go
type Handle interface {
    // Get retrieves a component in the current Category context. "" means Default.
    Get(ctx context.Context, name string) (any, error)
    // In switches navigation to a different category/scope.
    In(category Category, opts ...RegisterOption) Handle
    // BindConfig performs high-performance type-safe assignment of current config.
    BindConfig(target any) error
}
```

### 3.3 Protocol Definitions
- **Extractor**: `func(root any) (*ModuleConfig, error)` - Extracts structured data using interface sniffing.
- **ModuleConfig**: Standardized container for `[]ConfigEntry` and an `Active` winner.

## 4. Execution Flow: The Lifecycle

1.  **Creation**: `App` initializes `NewContainer(nil)`.
2.  **Registration**: Core factories are registered with **Standard Priorities** (100-500).
3.  **Loading**: Configuration is loaded; `BindRoot(bizConfig)` is called.
4.  **Warm-Up**: `Init(ctx)` triggers:
    -   Sort categories by Priority.
    -   Call `Extractor` for each category (Interface Sniffing).
    -   Sequentially instantiate each `ConfigEntry` using `Provider`.
5.  **Runtime**: Components are retrieved via `Cast[T](ctx, h.In(Category), name)`.

## 5. Directory Structure

```
runtime/engine/
├── context/      # Static Category, Scope, and Priority constants.
├── container/    # Core IoC, state machine, and navigation implementation.
├── protocol/     # Interface-based configuration block protocols.
├── helpers.go    # Type-safe generic helpers (Cast, GetDefault, BindConfig).
└── DESIGN.md     # This authoritative design document.
```
