FROM golang:1.16-alpine as builder
WORKDIR /app
RUN apk add --no-cache make
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN make build


FROM alpine
WORKDIR /app
COPY --from=builder /app/main /app/main
CMD ["/app/main"]
