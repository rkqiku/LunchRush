.PHONY: build run run-dapr test clean

build:
	/usr/local/bin/go build -o bin/lunchservice cmd/lunchservice/main.go

run:
	/usr/local/bin/go run cmd/lunchservice/main.go

run-dapr:
	dapr run \
		--app-id lunchservice \
		--app-port 8080 \
		--dapr-http-port 3500 \
		--dapr-grpc-port 50001 \
		--components-path ./config \
		/usr/local/bin/go run cmd/lunchservice/main.go

test:
	/usr/local/bin/go test ./...

clean:
	rm -rf bin/

deps:
	/usr/local/bin/go mod download
	/usr/local/bin/go mod tidy

docker-build:
	docker build -t lunchservice:latest .

docker-run:
	docker run -p 8080:8080 lunchservice:latest