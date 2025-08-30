FROM golang:1.22-alpine AS build
WORKDIR /app
COPY ./main.go .
RUN go build -ldflags="-s -w" -o sendmail main.go

FROM scratch
WORKDIR /app
COPY --from=build /etc/passwd /etc/passwd
COPY --from=build /etc/group /etc/group
COPY --from=build --chown=1000:1000 --chmod=100 /app/sendmail .

USER 1000:1000
ENTRYPOINT ["/app/sendmail"]
