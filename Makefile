pb:
	protoc --proto_path=./protobuf \
	--go_out=./protobuf --go_opt=paths=source_relative \
	--go-grpc_out=./protobuf --go-grpc_opt=paths=source_relative \
	./protobuf/*.proto && \
	protoc-go-inject-tag -input="./protobuf/*.pb.go" && \
	echo "proto code generation successful"

build:
	go build -o webhook main.go	

	# @rm -rf ./protobuf/*.proto && \
	#