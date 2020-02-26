FROM golang:1.14


WORKDIR /go/src/app

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/kpack-build-backend .
#ENV GO111MODULE=on

FROM debian:stretch-slim

WORKDIR /app

COPY --from=0 /app/kpack-build-backend /app/kpack-build-backend

EXPOSE 8080


CMD ["/app/kpack-build-backend"]