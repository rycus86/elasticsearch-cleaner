FROM golang:1.10 as builder

ARG CC=""
ARG CC_PKG=""
ARG CC_GOARCH=""

ADD . /go/src/github.com/rycus86/elasticsearch-cleaner
WORKDIR /go/src/github.com/rycus86/elasticsearch-cleaner

RUN if [ -n "$CC_PKG" ]; then \
      apt-get update && apt-get install -y $CC_PKG; \
    fi \
    && export CC=$CC \
    && export GOOS=linux \
    && export GOARCH=$CC_GOARCH \
    && export CGO_ENABLED=0 \
    && go build -o /var/tmp/app -v .

FROM scratch

LABEL maintainer "Viktor Adam <rycus86@gmail.com>"

COPY --from=builder /var/tmp/app /elasticsearch-cleaner

CMD [ "/elasticsearch-cleaner" ]
