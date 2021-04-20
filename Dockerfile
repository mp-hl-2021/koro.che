FROM golang:1.16.2-alpine3.13 as builder
RUN mkdir /build
ADD . /build/
WORKDIR /build
RUN CGO_ENABLED=0 GOOS=linux go build -a -o koro.che main.go


# generate clean, final image for end users
FROM alpine:3.13
COPY --from=builder /build/koro.che .

# executable
ENTRYPOINT [ "./koro.che" ]