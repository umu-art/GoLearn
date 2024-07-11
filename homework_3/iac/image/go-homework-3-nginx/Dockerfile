FROM golang:1.22 AS build
LABEL author="vikazeni"

WORKDIR /app

RUN apt update
RUN apt install -y protobuf-compiler
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

COPY . .

RUN go mod tidy

RUN export PATH=$PATH:$GOPATH/bin; \
    protoc --go_out=. --go-grpc_out=. ./homework_3/proto/protocol.proto

RUN CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ./client-win.exe ./homework_3/client
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./client ./homework_3/client

FROM nginx:stable-alpine3.17
LABEL author="vikazeni"

COPY --from=build /app/homework_3/README.md /usr/share/nginx/html/index.html
COPY --from=build /app/client-win.exe /usr/share/nginx/html/client-win.exe
COPY --from=build /app/client /usr/share/nginx/html/client
COPY --from=build /app/homework_3/iac/image/go-homework-3-nginx/nginx.conf /etc/nginx/conf.d/default.conf