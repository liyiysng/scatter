::ECHO off

::项目目录
SET PROJECT_PATH=../../
::导出目录(由于GO是按照包名目录导出)
SET GO_EXPORT_PATH=%PROJECT_PATH%
::引用目录
SET REF_PATH=%PROJECT_PATH%proto
:: proto 文件所在目录
SET PROTO_PATH=%PROJECT_PATH%proto
:: 所有需编译的proto文件
SET CMP_PROTOS=%PROTO_PATH%/node/message/proto/*.proto %PROTO_PATH%/cluster/sessionpb/*.proto %PROTO_PATH%/cluster/subsrvpb/*.proto
:: 所有需编译的GRPC proto
SET CMP_GRPC_PROTOS=%PROTO_PATH%/cluster/sessionpb/*.proto %PROTO_PATH%/cluster/subsrvpb/*.proto

:: 编译proto
protoc --go_out=paths=source_relative:%GO_EXPORT_PATH% --proto_path=%REF_PATH% %CMP_PROTOS%
:: 编译GRPC
protoc --go-grpc_out=paths=source_relative:%GO_EXPORT_PATH% --proto_path=%REF_PATH% %CMP_GRPC_PROTOS%
pause