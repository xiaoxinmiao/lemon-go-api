FROM pangpanglabs/golang:builder AS builder

RUN go get github.com/fatih/structs

ADD . /go/src/go-api
WORKDIR /go/src/go-api
ENV CGO_ENABLED=0
RUN go build -o go-api

FROM alpine
RUN apk --no-cache add ca-certificates
# FROM scratch
WORKDIR /go/src/go-api
COPY --from=builder /go/src/go-api/*.yml /go/src/go-api/
COPY --from=builder /go/src/go-api/go-api /go/src/go-api/
COPY --from=builder /swagger-ui/ /go/src/go-api/swagger-ui/
COPY --from=builder /go/src/go-api/index.html /go/src/go-api/swagger-ui/


EXPOSE 8080

CMD ["./go-api"]
