protoc --go_out=.. --proto_path=../proto ../proto/node/node.proto
protoc --go-grpc_out=.. --proto_path=../proto ../proto/node/node.proto