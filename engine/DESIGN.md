# Runtime Engine Design (v37 - 2026-03-06)

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

## 3. Dimensional Management: Scope & Tags

The engine manages components through two orthogonal dimensions to achieve strict isolation and flexible aggregation.

### 3.1 Scope: Vertical Isolation (Physical/Logical Environment)
Scope defines **where** a component exists. It represents a hard boundary between environments.

- **Registration (Multi-Selection)**: A component can be registered to multiple scopes (e.g., a Logger belongs to both `server` and `client` scopes).
- **Retrieval (Single Perspective)**: A `Handle` must reside in exactly one scope at any time. There is no implicit fallback between non-global scopes.
- **Matching Rule**: `RequestedScope` must exist in the component's `RegisteredScopes` list.
- **Global Scope**: `_global` acts as a fallback. Components with no scope specified are considered global and visible to all.

### 3.2 Tags: Horizontal Filtering (Capabilities & Identities)
Tags define **what** a component is or what it can do. It represents a soft filter for capability aggregation.

- **Registration (Single Identity)**: Each component registration must have at most **one** unique `Tag`. This ensures a clear identity for the component instance.
- **Retrieval (Multi-Capability)**: A `Handle` (Perspective) can request **multiple** tags simultaneously to aggregate various capabilities.
- **Matching Rule (Membership)**: A component is visible to a perspective if:
    1.  The component is **Common** (no tag registered).
    2.  The component's **Identity Tag** is present in the perspective's requested tag set.
- **Full Perspective**: A perspective with no requested tags can see all component identities.

### 3.3 Rule Summary (S-Multi-Single, T-Single-Multi)

| Dimension | Registration (Provider) | Retrieval (Handle/In) | Logic |
| :--- | :--- | :--- | :--- |
| **Scope** | **Multiple** (belong to many) | **Single** (reside in one) | Vertical Isolation |
| **Tags** | **Single** (unique identity) | **Multiple** (possess many) | Horizontal Capability |

## 4. Configuration & Data Flow

### 4.1 Resolution Priority
1.  Local Resolver (from `Load` options).
2.  Provider Resolver (from `Register` options).
3.  Category Resolver (from `NewContainer` options).
4.  **Pass-through Mode**: If no resolver matches, the `source` (Root) is passed directly to the provider.

### 4.2 Configuration Consumption
Providers should use the recommended `comp.AsConfig[T](h)` helper to consume data.
- **Performance**: Near-zero overhead using type assertions.

```go
func MyProvider(ctx context.Context, h component.Handle) (any, error) {
    cfg, err := comp.AsConfig[MyConfig](h)
    if err != nil {
        return nil, err
    }
    return NewInstance(cfg), nil
}
```

## 5. Instance Retrieval Helpers

The engine leverages the `comp` helper package for type-safe instance retrieval:

- **`comp.Get[T](ctx, h, name)`**: Retrieves and casts a component.
- **`comp.GetDefault[T](ctx, h)`**: Retrieves the active/default instance.
- **`comp.GetMap[T](ctx, h)`**: Collects all instances of type T into a map.
- **`comp.Iter[T](ctx, h)`**: Returns a type-safe iterator.

## 6. Implementation Integrity
- **Engine Core**: `engine/container/container.go` contains **zero reflection**.
- **Perspective Consistency**: `In()` resets the scope to `_global` by default but preserves Tags to maintain capability context.
- **Identity Safety**: Cached instances are only reused if the requester's Provider identity (Tag) matches the instance's birth identity.
