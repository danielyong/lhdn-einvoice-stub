FROM golang:alpine as build
WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN go build -o ./app_bin
FROM alpine:latest as runner
COPY --from=build /usr/src/app/app_bin /usr/local/bin/app
CMD ["app"]
