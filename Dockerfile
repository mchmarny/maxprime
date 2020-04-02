# BUILD
FROM golang:latest as builder

# copy
WORKDIR /src/
COPY . /src/

# build
ENV GO111MODULE=on
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -a -tags netgo \
    -ldflags '-w -extldflags "-static"' \
    -mod vendor \
    -o maxprime

# RUN
FROM gcr.io/distroless/static:nonroot
COPY --from=builder /src/maxprime /app/
COPY --from=builder /src/templates /app/templates/
COPY --from=builder /src/static /app/static/

WORKDIR /app
ENTRYPOINT ["./maxprime"]
