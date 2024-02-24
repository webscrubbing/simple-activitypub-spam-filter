FROM golang:latest as gobuilder
WORKDIR /go/src/app
COPY go.mod ./
RUN go mod download
COPY main.go ./
RUN CGO_ENABLED=0 go build -o /go/bin/app

FROM gcr.io/distroless/base
COPY --from=gobuilder /go/bin/app /
CMD ["/app"]