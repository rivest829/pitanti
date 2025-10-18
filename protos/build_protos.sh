#!/bin/bash

echo "正在编译当前目录下的所有 .proto 文件..."

# 设置输出目录为当前目录
OUTPUT_DIR="."

# 查找并编译所有 .proto 文件
for proto_file in *.proto; do
    if [ -f "$proto_file" ]; then
        echo "编译: $proto_file"
        protoc --go_out=$OUTPUT_DIR \
               --go-grpc_out=$OUTPUT_DIR \
               "$proto_file"
    fi
done

echo "编译完成！"