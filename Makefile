test:
	go test ./helper/... ./service/...

coverage:
	go test ./helper/... ./service/... -cover

vendor:
	glide install
