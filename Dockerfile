FROM golang:1.14


WORKDIR /go/src/app

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/kpack-build-backend .
#ENV GO111MODULE=on

FROM scratch

COPY --from=0 /app/kpack-build-backend /kpack-build-backend

EXPOSE 8080


CMD ["/kpack-build-backend"]