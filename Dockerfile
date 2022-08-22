# syntax=docker/dockerfile:1

##
## Build
##
FROM golang:1.18-buster AS build

WORKDIR /go/src/

COPY /app ./app
COPY go.mod ./
COPY go.sum ./

RUN ls -la

RUN go mod download

RUN go build -o app ./app
RUN ls -la

FROM gcr.io/distroless/base-debian10
WORKDIR /
EXPOSE 8080
COPY --from=build /go/src/app /app
CMD ["./app/app"]
