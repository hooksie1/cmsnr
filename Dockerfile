FROM golang:alpine as builder
WORKDIR /app
ENV IMAGE_TAG=dev
RUN apk update && apk upgrade && apk add --no-cache ca-certificates git
RUN update-ca-certificates
ADD . /app/
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags="-s -w -X 'github.com/hooksie1/cmsnr/cmd.Version=$(printf $(git describe --tags | cut -d '-' -f 1)-$(git rev-parse --short HEAD))'" -installsuffix cgo -o cmsnrctl .


FROM scratch

COPY --from=builder /app/cmsnrctl .
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENTRYPOINT ["./cmsnrctl"]
