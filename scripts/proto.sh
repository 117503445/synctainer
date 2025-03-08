#!/usr/bin/env sh

set -ev

# protoc --go_out=/workspace --go-grpc_out=/workspace --proto_path /workspace/proto /workspace/proto/synctainer.proto
protoc --go_out=/workspace --twirp_out=/workspace --proto_path /workspace/proto /workspace/proto/synctainer.proto

cd /workspace/proto
twirpscript synctainer.proto
mv synctainer.pb.ts /workspace/fe/src/rpc