FROM golang:1.22 AS build
LABEL author="vikazeni"

WORKDIR /app

COPY . .

RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ./client-win.exe ./homework_2/client
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./client ./homework_2/client

FROM nginx:stable-alpine3.17
LABEL author="vikazeni"

COPY --from=build /app/homework_2/README.md /usr/share/nginx/html/index.html
COPY --from=build /app/client-win.exe /usr/share/nginx/html/client-win.exe
COPY --from=build /app/client /usr/share/nginx/html/client
COPY --from=build /app/homework_2/iac/image/go-homework-2-nginx/nginx.conf /etc/nginx/conf.d/default.conf