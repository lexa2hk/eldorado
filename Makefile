up:
	docker compose up -d

down:
	docker compose down --rmi all

lint:
	golangci-lint run ./...

test:
	go test -v -race -coverprofile=c.out ./... \
	&& go tool cover -html=c.out \
	&& rm c.out

proto:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		internal/services/auth/proto/auth.proto

mock:
	go generate ./...
