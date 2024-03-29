FROM golang:1.21-alpine

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

# run as unprivileged user from now on
RUN chown -R 1000 /go
USER 1000

ENV CGO_ENABLED=0

# an internal docker volume that will contain the go building cache to ensure we don't have to build from scratch
# every time the container starts, which increases the startup time
ENV GOCACHE=/go/
ENV GOTMPDIR=/go/tmp/
RUN mkdir -p /go /go/cache /go/pkg /go/tmp
VOLUME /go/cache
VOLUME /go/pkg
VOLUME /go/tmp

# set default working directory
WORKDIR /go/src/app/

# already download and cache dependencies
COPY --chown=1000 go.mod go.sum /go/src/app/
RUN go mod download

# compile the program once to fill the build cache
COPY --chown=1000 . .
RUN go build -o /dev/null ./cmd/cli/...

# 8080 is used to host health endpoint, 40000 is used to run the go debugger (dlv)
EXPOSE 8080 40000

# run reflex that will re-launch the go program when changes are detected
ENTRYPOINT ["/go/bin/reflex", "--sequential", "--config=/go/src/app/kubernetes/docker/reflex.conf"]
