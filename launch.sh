#!/bin/bash -e

SERVICE_DIR=${PWD}
SERVICE_BINARY=scrapper-service.o
DEBUG_MODE=true
URL="http://localhost:8080"

case "$1" in
    "run")
        export DEBUG_MODE=${DEBUG_MODE}
        ./${SERVICE_BINARY}
        ;;
    "run-docker")
        docker run -it --rm --name scrapper-service \
            -e DEBUG_MODE=${DEBUG_MODE} \
            -p 8080:8080 \
            -v ${SERVICE_DIR}/cfg/config.yml:/etc/scrapper-service/config.yml \
            scrapper-service:latest
        ;;
    "test")
    curl -X POST ${URL} \
    -H 'Content-Type: application/json8' \
    -d @- << EOF
[
    "http://yandex.ru/",
    "http://google.com/"
]
EOF
        ;;
    "test-multi")
    curl -X POST ${URL} \
    -H 'Content-Type: application/json8' \
    -d @- << EOF
[
    "https://vc.ru/",
    "https://vc.ru/",
    "https://roem.ru/"
]
EOF
        ;;
    *)
        echo "Usage: $(basename $0) <run> | <run-docker> | <test> | <test-multi>"
        exit 1
       ;;
esac
