PROTO_DIR=.\proto
GO_OUT=.\gen

.PHONY: proto
proto:
	mkdir -p $(GO_OUT)
	protoc -I=$(PROTO_DIR) --go_out=$(GO_OUT) --go-grpc_out=$(GO_OUT) $(PROTO_DIR)/*.proto

.PHONY: tidy
tidy:
	go mod tidy