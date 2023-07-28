```
$ protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    routeguide/route_guide.proto
```

- route_guide.pb.go, which contains all the protocol buffer code to populate, serialize, and retrieve request and response message types.
- route_guide_grpc.pb.go, which contains the following:
  - An interface type (or stub) for clients to call with the methods defined in the RouteGuide service.
  - An interface type for servers to implement, also with the methods defined in the RouteGuide service.

route_guide.proto là một tệp proto được định nghĩa trong ngôn ngữ proto3. Nó định nghĩa cấu trúc dữ liệu và các thông điệp cho ứng dụng và cũng định nghĩa các API mà ứng dụng sử dụng để giao tiếp với nhau.

Sau khi định nghĩa route_guide.proto, nó có thể được sử dụng để tạo ra mã cho các ngôn ngữ khác nhau. Ví dụ, route_guide.pb.go được tạo ra từ route_guide.proto và là một tập tin mã Go chứa các cấu trúc dữ liệu và mã để truy cập và sử dụng các thông điệp trong route_guide.proto.

route_guide_grpc.pb.go cũng được tạo ra từ route_guide.proto, nhưng nó là một tập tin mã Go chứa các định nghĩa mã RPC gRPC cho các API được định nghĩa trong route_guide.proto. Nó cung cấp các hàm để sử dụng các API trong route_guide.proto thông qua giao thức gRPC.
