build:
	protoc -I. --go_out=plugins=micro:$(GOPATH)/src/github.com/jcorry/bigbookwords/dictionary-service/proto \
		--proto_path=$(GOPATH)/src/github.com/jcorry/bigbookwords/proto/ \
		dictionary/dictionary.proto
	GOOS=linux GOARCH=amd64
	docker build -t dictionary-service .

run:
	docker run -p 50051:50051 \
		-e MICRO_SERVER_ADDRESS=:50051 \
		-e MICRO_REGISTRY=mdns dictionary-service