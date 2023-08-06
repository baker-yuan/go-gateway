echo "生成 rpc 代码"

# 输出目录
OUT=./

protoc \
--go_out=${OUT} \
--go-grpc_out=${OUT} \
--go-grpc_opt=require_unimplemented_servers=false \
service.proto



