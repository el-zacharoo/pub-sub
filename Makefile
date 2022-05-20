# load env vars
# -include .env
# export mongoUri := $(value mongoUri)

# proto generates code from the most recent proto file(s)
.PHONY: proto
proto:
	cd proto && buf mod update
	buf generate
	buf build
	cd proto && buf push

.PHONY: run
run:
	dapr run \
		--dapr-http-port 3501 \
		--dapr-grpc-port 50002 \
		--app-id sub \
		--app-port 8081 \
		--app-protocol grpc \
		--components-path ./.dapr \
		go run .

.PHONY: test
test:
	go test -v ./...