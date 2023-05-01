echo "生成 grpc 代码"

OUT=./

#protoc \
#--go_out=${OUT} \
#--go-grpc_out=${OUT} \
#--go-grpc_opt=require_unimplemented_servers=false \
#helloworld.proto


protoc \
--proto_path=./ \
--go_out=paths=source_relative:./ \
--go-http_out=paths=source_relative:./ \
--go-grpc_out=paths=source_relative:./ \
--openapi_out=fq_schema_naming=true,default_response=false:. \
helloworld.proto