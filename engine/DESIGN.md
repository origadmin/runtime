# Engine Container Design Principles

## 1. Naming & Redirection Logic

### _default Marker Principle
`_default` (defined as `component.DefaultName`) is a **system marker**, not a concrete component name. It acts as a redirection target when a user requests a component without specifying a name.

### Default Selection Priority
The target for `_default` is determined by this strict priority:
1.  **Named "default":** If an entry is explicitly named `"default"`.
2.  **Explicit Active:** If the `ModuleConfig.Active` field is set.
3.  **Single Entry Fallback:** If and only if there is exactly one entry in the `ModuleConfig`.

### Dynamic Naming
If a single source is bound without a specific name (e.g., when `Resolver` returns `nil`), the entry should be named after its `Category` (e.g., "logger"), and then mapped as the `_default`.

### Instance Integrity
A component must always be stored under its actual name. The `_default` entry in the instance map should be a reference/link to the actual named entry's metadata to ensure consistency.

## 2. Advanced Features

### Inject Protocol
Injected instances (via `Inject`) are treated as pre-instantiated (`StatusReady`). 
- If no explicit name is provided for an injected instance, it is given a system-generated name (e.g., `_injected_logger`).
- Injected instances follow the same redirection rules as loaded components.

### DefaultEntries (Unconditional Seeding)
The `DefaultEntries` field in `RegistrationOptions` allows a provider to declare names that should always exist in the container's registry for that category, even if the external configuration source is empty or missing those names. 
- During `Load`, these entries are seeded as `StatusNone` if they don't already exist.
- This ensures that a container can still instantiate a "default" set of components without explicit configuration.
