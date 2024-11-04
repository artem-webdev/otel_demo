generate_proto_user:
	protoc --experimental_allow_proto3_optional --go_out=. --go-grpc_out=. ./backend/pkg/api/grpc/user.proto
