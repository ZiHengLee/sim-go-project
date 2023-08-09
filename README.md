## golang+gin+grpc(golang快速启动项目)

#### 安装protoc，protoc-go-inject-tag，protoc-gen-go-grpc:
```
https://github.com/favadi/protoc-go-inject-tag
go install google.golang.org/protobuf/cmd/protoc-gen-go
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc
```

#### 生成pb文件
```
./gen_proto.sh
```

#### 启动
```go
docker-compose up -d
go build main.go
main config.toml
```