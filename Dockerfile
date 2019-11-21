FROM registry-v2.revinate.net/common/go:go-1.9.2

EXPOSE 3001
WORKDIR /go/src/github.com/revinate/go-utils
CMD go test ./helper/... ./service/...
