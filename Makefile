test:
	docker-compose up goutils_test

coverage:
	go test ./helper/... ./service/... -cover

vendor:
	glide install
