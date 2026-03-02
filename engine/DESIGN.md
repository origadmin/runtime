# Runtime Engine Design (v31 - 2026-03-02)

## 1. Core Philosophy: The Unified Starter and Bridge
The Runtime Engine is the unified entry point for both **Bootstrapping** (configuration loading) and **Bridging** (component instantiation). It manages the entire lifecycle from the first byte of config to the last initialized service.

## 2. Updated Project Structure

To avoid naming conflicts and improve logical grouping:
- `runtime/engine/metadata`: (Formerly `context`) Defines `Category`, `Scope`, and `Priority` constants.
- `runtime/engine/bootstrap`: (Migrated) Handles configuration loading and result encapsulation.
- `runtime/engine/container`: The IoC engine that performs dependency resolution.
- `runtime/engine/protocol`: Interfaces for config blocks and extractors.

## 3. The Bootstrap Integration

The Bootstrapping phase is now an integral part of the Engine's lifecycle:
1. **Load**: `engine/bootstrap` reads YAML/Env and parses the global business object.
2. **Bind**: The loaded object is automatically bound to the `engine/container`.
3. **Dispatch**: Categories are initialized based on their `Priority`.

## 4. Meta-Definitions (metadata package)

### 4.1 Scope & Category
- **Scope**: `Global`, `Server`, `Client`.
- **Category**: `Registry`, `Middleware`, `Database`, etc.

### 4.2 Priority
Standardized execution order (100 - 500+).

## 5. Interface Protocol (No Changes to Logic)
Registry and Handle remain the primary interfaces, utilizing `metadata` types.

## 6. Directory Layout

```
runtime/engine/
├── metadata/     # Category, Scope, Priority constants (No naming conflict with Go context)
├── bootstrap/    # Config loading and result implementation
├── container/    # IoC, State Machine, Navigation
├── protocol/     # ConfigBlock and Extractor interfaces
├── helpers.go    # Casting and retrieval helpers
└── engine.go     # Unified public API
```
