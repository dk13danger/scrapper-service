all: build

.PHONY: glide
glide:
	glide install

.PHONY: build
build: glide
	CGO_ENABLED=0 go build -a -installsuffix cgo -o ./scrapper-service.o .

.PHONY: recompile
recompile:
	CGO_ENABLED=0 go build -i -installsuffix cgo -o ./scrapper-service.o .

.PHONY: docker
docker:
	docker build -t scrapper-service:latest .
