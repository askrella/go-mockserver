FROM golang:1.21 as build

WORKDIR /go/src/app
COPY go.mod go.mod
#COPY go.sum go.sum
COPY cmd cmd
COPY internal internal

RUN go mod download

RUN CGO_ENABLED=0 go build -o /go/bin/app ./cmd/main.go

FROM gcr.io/distroless/static-debian12

COPY --from=build /go/bin/app /
CMD ["/app"]