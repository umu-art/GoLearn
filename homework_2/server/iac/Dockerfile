FROM golang:1.22
LABEL author="vikazeni"

WORKDIR /app

COPY . .

RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o ./build/homework_2 ./homework_2/server

ENTRYPOINT ["./build/homework_2"]