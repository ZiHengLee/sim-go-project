#!/bin/bash
protoc -I . *.proto --go_out=. --go-grpc_out=.

SEARCH_DIR="../proto"

# 使用 find 命令来查找所有的 .pb.go 文件
find "$SEARCH_DIR" -name "*.pb.go" | while read -r file; do
  if [ -f "$file" ]; then
    protoc-go-inject-tag -input=$file
    echo "finished: $file"
  fi
done