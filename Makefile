build:
	 GOPATH=`pwd` GO111MODULE=off go build ./

up:
	export PORT=8889 AUTH=false DUMP_PATH=/home/xed/go/bin/grpc.dev-dash.dump && ./read_dump

start:
	./grpc-dump --port=12345 --key=/home/xed/project/soter/ssl/grpc.dev-dash.soteranalytics.com/privkey6.pem --cert=/home/xed/project/soter/ssl/grpc.dev-dash.soteranalytics.com/fullchain6.pem --proto_roots=/home/xed/go/src/soter-api/proto/ > grpc.dev-dash.dump

.PHONY: up build