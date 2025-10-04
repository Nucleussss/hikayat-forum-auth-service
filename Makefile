
PROTO_PATH = ../hikayat-proto
PROTO_FILE = auth/v1/auth.proto
GEN_DIR = ./api

generate:
	protoc \
	  --go_out=$(GEN_DIR) \
	  --go-grpc_out=$(GEN_DIR) \
	  --go_opt=paths=source_relative \
	  --go-grpc_opt=paths=source_relative \
	  -I $(PROTO_PATH) \
	  $(PROTO_FILE)

.PHONY: generate