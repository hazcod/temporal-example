FROM golang:1.21-alpine AS builder

# always statically compile
ENV CGO_ENABLED=0

# add unprivileged user
RUN adduser -s /bin/true -u 1000 -D -h /app app \
  && sed -i -r "/^(app|root)/!d" /etc/group /etc/passwd \
  && sed -i -r 's#^(.*):[^:]*$#\1:/sbin/nologin#' /etc/passwd

# add ca certificates and timezone data files
# hadolint ignore=DL3018
RUN apk add --upgrade --no-cache ca-certificates tzdata git \
    && go install github.com/cespare/reflex@latest \
    && go install github.com/go-delve/delve/cmd/dlv@latest

# copy in our app source code
COPY . /go/src/app
WORKDIR /go/src/app

# compile our application
RUN go build -o /app ./cmd/cli/...

#
# ---
#

# empty base image
FROM scratch

# add-in our timezone data file
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# add-in our unprivileged user
COPY --from=builder /etc/passwd /etc/group /etc/shadow /etc/

# add-in our ca certificates
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# from now on, run as the unprivileged user
USER 1000

# copy over the compiled binary
COPY --from=builder /app/cli /

# health endpoint listens on 8080
EXPOSE 8080

# app binary
ENTRYPOINT ["/cli", "--config=/config/config.yaml"]