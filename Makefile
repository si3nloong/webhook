generate-protobuf:
	@rm -rf ./grpc/proto/* && \
	protoc --proto_path=./grpc/api \
	--go_out=./grpc/proto --go_opt=paths=source_relative \
	--go-grpc_out=./grpc/proto --go-grpc_opt=paths=source_relative \
	./grpc/api/*.proto && \
	protoc-go-inject-tag -input="./grpc/proto/*.pb.go" && \
	echo "proto code generation successful"

build:
	go build -o rwhook main.go	