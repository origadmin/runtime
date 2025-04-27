# This is an upload intermediary service

## The features to be done now are as follows

- [ ] The client uploads the file and receives it from the server
- [ ] The client bridges the http upload and the server receives it.
- [ ] Both the client and the server should support HTTP and grpc.
- [ ] Add gRPC upload function
- [ ] Unified HTTP and gRPC upload interfaces
- [ ] Support multi-file upload (by calling uploader/receiver multiple times)
- [ ] Provides better error handling and resource management
- [ ] Added folder upload support
- [ ] Implemented the function of resumable transmission
- [ ] Buffer pools are used to optimize memory usage
- [ ] Maintain interface consistency
- [ ] Better error handling is provided
- [ ] Provides HTTP to gRPC bridge upload
- [ ] The bridge from HTTP to gRPC supports resumable upload
- [ ] HTTP to gRPC bridging optimizes memory usage using buffer pools
- [ ] Contextual support is provided for the HTTP to gRPC bridge upload

## Implementation of the underlying logic

- [ ] The user clicks the upload button to upload the file (HTTP).
- [ ] The server receives the uploader that bridges HTTP and gRPC) and forwards it to the grpc server.
- [ ] The grpc server receives the data and saves it to the file/OSS.

Examples:
```
func main() {

 // Create a gRPC Builder

  grpcBuilder := NewBuilder(
    WithServiceType(ServiceTypeGRPC),
    WithURI("grpc-server-address:port"),
  )

  // Create a bridge uploader
  bridgeUploader := NewBridgeUploader(grpcBuilder, &uploadBuilder{
    bufPool: &sync.Pool{
      New: func() interface{} {
        return make([]byte, 32*1024)
      },
    },
  })

  // Register the HTTP processor
  http.HandleFunc("/upload", bridgeUploader.ServeHTTP)
  http.HandleFunc("/upload/resume", bridgeUploader.ServeHTTPWithResume)
  // Start the HTTP server
  http.ListenAndServe(":8080", nil)

}
```

### Workflow:

- [ ] The client uploads the file via HTTP
- [ ] BridgeUploader receives HTTP requests
- [ ] Use HTTPReceiver to read the contents of the file
- [ ] Create a gRPC uploader and forward the file
- [ ] Return the gRPC response to the HTTP client