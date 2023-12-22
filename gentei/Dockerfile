FROM golang:1.21-bullseye as build

WORKDIR /go/src/app
COPY . .
RUN go get -d -v ./...
RUN go vet -v
RUN go test -v ./...

RUN go build --tags "json1" -o /go/bin/app

FROM gcr.io/distroless/base-debian11
COPY --from=build /go/bin/app /
ENTRYPOINT ["/app"]