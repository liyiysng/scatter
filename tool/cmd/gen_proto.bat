ECHO off

::项目目录
SET PROJECT_PATH=../../
::导出目录(由于GO是按照包名目录导出)
SET GO_EXPORT_PATH=%PROJECT_PATH%
::引用目录
SET REF_PATH=%PROJECT_PATH%proto
:: node 相关proto文件
SET CMP_NODE_PROTO_PATH=%PROJECT_PATH%proto/node/*.proto
:: testing 相关proto文件
SET CMP_TESTING_PROTO_PATH=%PROJECT_PATH%proto/testing/*.proto

:: 编译proto
protoc --go_out=%GO_EXPORT_PATH% --proto_path=%REF_PATH% %CMP_NODE_PROTO_PATH% %CMP_TESTING_PROTO_PATH%
:: 编译GRPC
protoc --go-grpc_out=%GO_EXPORT_PATH% --proto_path=../../proto ../../proto/node/node.proto

pause