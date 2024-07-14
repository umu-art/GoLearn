FROM golang:1.22
LABEL author="vikazeni"

WORKDIR /app

RUN apt update
RUN apt install -y protobuf-compiler
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

COPY . .

RUN go mod tidy

RUN export PATH=$PATH:$GOPATH/bin; \
    protoc --go_out=. --go-grpc_out=. ./homework_4/proto/protocol.proto

RUN CGO_ENABLED=0 GOOS=linux go build -o ./build/homework_4 ./homework_4/server

ENTRYPOINT ["./build/homework_4"]