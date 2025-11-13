#!/bin/bash

set -e

PROTO_ROOT="../proto"
GO_OUT="../generated/go"

echo "Generating Go code from protobuf..."

# Clean previous generated code
rm -rf $GO_OUT
mkdir -p $GO_OUT

# Generate Go code for all protobuf files
find $PROTO_ROOT -name "*.proto" | while read proto_file; do
    echo "Generating Go code for: $proto_file"

    protoc -I=$PROTO_ROOT \
        --go_out=$GO_OUT \
        --go_opt=paths=source_relative \
        --go-grpc_out=$GO_OUT \
        --go-grpc_opt=paths=source_relative \
        "$proto_file"
done

echo "Go code generation completed! Output: $GO_OUT"
