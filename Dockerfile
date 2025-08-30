FROM golang:1.22-alpine AS build
WORKDIR /app
COPY ./main.go .
RUN go build -ldflags="-s -w" -o sendmail main.go
RUN apk add --no-cache ca-certificates

FROM scratch
WORKDIR /app
COPY --from=build /etc/passwd /etc/passwd
COPY --from=build /etc/group /etc/group
COPY --from=build --chown=1000:1000 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build --chown=1000:1000 --chmod=100 /app/sendmail .

USER 1000:1000
ENTRYPOINT ["/app/sendmail"]
