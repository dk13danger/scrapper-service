# ===== GOLANG BUILDER IMAGE
FROM golang:1.9 as builder

ADD . ${GOPATH}/src/github.com/dk13danger/scrapper-service
WORKDIR ${GOPATH}/src/github.com/dk13danger/scrapper-service

# Install glide
RUN curl https://glide.sh/get | sh

# Install dependencies and build executable file
RUN glide install \
 && CGO_ENABLED=0 go build -a -installsuffix cgo -o ./scrapper-service.o . \
 && mv ${GOPATH}/src/github.com/dk13danger/scrapper-service/scrapper-service.o /scrapper-service.o

# ===== FINAL IMAGE
FROM alpine:3.6
COPY --from=builder /scrapper-service.o /scrapper-service
RUN apk update && apk add ca-certificates
ENTRYPOINT ["/scrapper-service", "-config", "/etc/scrapper-service/config.yml"]
