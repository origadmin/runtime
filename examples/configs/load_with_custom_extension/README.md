Custom Proto Extension Example

This example shows how to define a custom proto configuration that integrates with the framework's runtime extension APIs so the runtime can decode directly into your custom proto messages.

What it demonstrates
- Define your own custom config proto that imports the framework's runtime extension types.
- Compose the unified application config by embedding framework config messages alongside your custom section.
- Provide a sample YAML that matches the generated messages and can be decoded by the runtime into your custom proto type.

Files
- protos/custom_config.proto: Custom configuration proto definition that imports runtime extension APIs.
- config/example.yaml: Example configuration YAML that includes framework sections and a custom section.
- config/bootstrap.yaml: Bootstrap file to load the example config through runtime.
- main.go: Minimal example showing how to decode the full config into your custom top-level proto message using the runtime's ConfigDecoder.

How to run
1) Generate Go code for the example proto (if not already generated):
   - Ensure buf/protoc is set up. The repo contains buf.yaml at runtime root.
   - You can place generation settings in a local buf.gen.yaml or use protoc directly.
2) Run the example:
   cd runtime/examples/protos/custom_extension_example
   go run ./

You should see the runtime load config/example.yaml and decode the custom section into your custom proto message.
