PROTO_FILES=$PWD/pb/*.proto
protoc --go_out=gen --go_opt=paths=source_relative \
  -I=$PWD/pb $PROTO_FILES