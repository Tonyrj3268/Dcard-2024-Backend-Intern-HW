FROM golang:1.21-bullseye as builder

WORKDIR /app

COPY . .

RUN go mod download \
    && CGO_ENABLED=0 GOOS=linux go build -o myapp cmd/main.go

# 使用輕量級基礎鏡像
FROM gcr.io/distroless/base-debian10

COPY --from=builder /app/myapp /

EXPOSE 8080

CMD ["/myapp"]