FROM golang:1.12-alpine AS build_base

RUN apk add --no-cache git

WORKDIR /tmp/grpcapp

# we want the populate the model cached based on the go {mod,sum} file

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

# build the go app

RUN go build -o  ./out/grpcapp .


# start the fresh from smaller image 

FROM alpine:3.9 

RUN apk add ca-certificates

COPY --from=build_base /tmp/grpcapp/out/grpcapp /app/grpcapp

# This container exposes port 8080 to the outside world

EXPOSE 8080

# Run the binary program produced by `go install`

CMD ["/app/grpcapp"]   

