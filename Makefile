.PHONY: run_postgres
run_postgres: |
	docker compose -f docker-compose-postgres.yml up -d

.PHONY: run_in_memory
run_in_memory: |
	docker compose -f docker-compose-in-memory.yml up -d

.PHONY: stop_postgres
stop_postgres: |
	docker compose -f docker-compose-postgres.yml down

.PHONY: stop_in_memory
stop_in_memory: |
	docker compose -f docker-compose-in-memory.yml down

.PHONY: generate_proto_rpc
generate_proto_rpc: |
	protoc --go-grpc_opt=require_unimplemented_servers=false --go-grpc_out=internal/generated protobuf/*_rpc.proto

.PHONY: generate_proto
generate_proto: |
	find protobuf -type f -name '*.proto' ! -name '*_rpc.proto' -exec protoc --go_out=internal/generated {} +

.PHONY: cover_out
cover_out: |
	go test -coverprofile=c.out ./... -coverpkg=./...
	cat c.out | grep -v "cmd" | grep -v "_mock.go" | grep -v ".pb" > tmp.out
	go tool cover -func=tmp.out

.PHONY: cover_html
cover_html: |
	go test -v ./... -coverprofile=c.out ./... -coverpkg=./...
	cat c.out | grep -v "cmd" | grep -v "_mock.go" | grep -v ".pb" > tmp.out
	go tool cover -html=tmp.out
