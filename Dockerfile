FROM registry.revinate.net/common/go:go-1.8

EXPOSE 3001
WORKDIR /go/src/github.com/revinate/go-utils
CMD go test ./helper/... ./service/...
