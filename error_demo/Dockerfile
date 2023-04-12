FROM golang:alpine AS builder
ENV CGO_ENABLED 0
ENV GOOS linux
ENV GOPROXY https://goproxy.cn,direct
WORKDIR /build
COPY . .
RUN go mod tidy
RUN go build -ldflags="-s -w" -o errdemo ./main.go

FROM alpine
WORKDIR /app
ENV TZ Asia/Shanghai
COPY --from=builder /build/errdemo /app/errdemo
CMD ["./errdemo"]