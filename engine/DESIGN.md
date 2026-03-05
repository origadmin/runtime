# Runtime Engine Design (v37 - 2026-03-05)

## 1. Core Philosophy: Zero-Reflection & Metadata-Driven
The Runtime Engine is a business-agnostic IoC container. It strictly avoids reflection in its core logic, relying on explicit metadata registration and high-performance Go type assertions for configuration handling.

## 2. Decoupled Lifecycle (The Timeline)

The engine follows a strict state machine to ensure configuration consistency.

### 2.1 Configuring Phase (Mutable)
- **Container Initialization**: Created via `engine.NewContainer(opts...)`.
- **Knowledge Injection**: Use `WithCategoryResolvers(map)` to provide extraction rules.
- **Registration**: Components register via `Register(Category, Provider, ...opts)`.
- **State**: Mutable. `Register` is allowed.

### 2.2 Running Phase (Immutable/Locked)
- **Locking**: Triggered by `Load(ctx, source, ...opts)`.
- **State**: Once `Load` is called, the container is **Locked**. Further `Register` calls will panic.
- **Metadata-Driven**: `Load` iterates through registered providers only.

## 3. Configuration & Data Flow

### 3.1 Resolution Priority
1.  Local Resolver (from `Load` options).
2.  Provider Resolver (from `Register` options).
3.  Category Resolver (from `NewContainer` options).
4.  **Pass-through Mode**: If no resolver matches, the `source` (Root) is passed directly to the provider.

### 3.2 Configuration Consumption
Providers should use the recommended `AsConfig[T](h)` helper to consume data.
- **Performance**: Near-zero overhead using type assertions.

```go
func MyProvider(ctx context.Context, h component.Handle) (any, error) {
    cfg, err := engine.AsConfig[MyConfig](h)
    if err != nil {
        return nil, err
    }
    return NewInstance(cfg), nil
}
```

## 4. Instance Retrieval Helpers

The engine provides a set of generic helpers in the `engine` package for type-safe instance retrieval:

- **`Get[T](ctx, h, name)`**: Retrieves and casts a component.
- **`GetOr[T](ctx, h, name)`**: Retrieves by name, or falls back to `_default`.
- **`ToMap[T](ctx, h)`**: Collects all instances of type T into a map.
- **`Iter[T](ctx, h)`**: Returns a type-safe iterator.

## 5. Implementation Purity
- **Engine Core**: `engine/container/container.go` contains **zero reflection**.
- **Isolation**: Business logic is kept in `runtime/provider.go`.
