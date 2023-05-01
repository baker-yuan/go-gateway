```bash
https://blog.csdn.net/xmcy001122/article/details/126618680
https://baijiahao.baidu.com/s?id=1648713282754138283&wfr=spider&for=pc
```

```bash
curl --location --request POST 'http://localhost:8000/admin/httpRules' \
--data-raw '{
    "page": 1,
    "pageSIze": 100,
    "search": {
        "application": "baker-blog-blog"
    }
}'
```

```bash
create database github.com/baker-yuan/go-gateway/go-gateway-admin;
use github.com/baker-yuan/go-gateway/go-gateway-admin;
```

```bash
$ kratos new go_gateway -r https://gitee.com/go-kratos/kratos-layout.git

🚀 Creating service go_gateway, layout repo is https://gitee.com/go-kratos/kratos-layout.git, please wait a moment.
Already up to date.
CREATED go_gateway/.gitignore (552 bytes)
🍺 Project creation succeeded go_gateway
💻 Use the following command to start the project 👇:

$ cd go_gateway
$ go generate ./...
$ go build -o ./bin/ ./...
$ ./bin/go_gateway -conf ./configs
```bash

```
cd internal/data
go run -mod=mod entgo.io/ent/cmd/ent new HttpRule
```bash

```sql
CREATE TABLE `http_statement`
(
  `id`                bigint(11)    NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `application`       varchar(128)  NOT NULL DEFAULT '' COMMENT '应用名称',
  `interface_type`    varchar(256)  NOT NULL DEFAULT '' COMMENT '接口协议 1-http 2-gRPC 3-Double',
  `method_name`       varchar(128)  NOT NULL DEFAULT '' COMMENT '接口方法',
  `config`            varchar(2000) NOT NULL DEFAULT '' COMMENT '接口特殊配置json格式',
  `uri`               varchar(128)  NOT NULL DEFAULT '' COMMENT '网关接口',
  `http_command_type` varchar(32)   NOT NULL DEFAULT '' COMMENT '接口类型 GET、POST、PUT、DELETE',
  `create_date`       timestamp     NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time`       timestamp     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE = InnoDB
DEFAULT CHARSET = utf8mb4 COMMENT '网关接口映射信息';
```


# Kratos Project Template

## Install Kratos
```
go install github.com/go-kratos/kratos/cmd/kratos/v2@latest
```
## Create a service
```
# Create a template project
kratos new server

cd server
# Add a proto template
kratos proto add api/server/server.proto
# Generate the proto code
kratos proto client api/server/server.proto
# Generate the source code of service by proto file
kratos proto server api/server/server.proto -t internal/service

go generate ./...
go build -o ./bin/ ./...
./bin/server -conf ./configs
```
## Generate other auxiliary files by Makefile
```
# Download and update dependencies
make init
# Generate API files (include: pb.go, http, grpc, validate, swagger) by proto file
make api
# Generate all files
make all
```
## Automated Initialization (wire)
```
# install wire
go get github.com/google/wire/cmd/wire

# generate wire
cd cmd/server
wire
```

## Docker
```bash
# build
docker build -t <your-docker-image-name> .

# run
docker run --rm -p 8000:8000 -p 9000:9000 -v </path/to/your/configs>:/data/conf <your-docker-image-name>
```

