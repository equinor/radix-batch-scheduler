FROM golang:1.17.6-alpine3.15 as builder
ENV GO111MODULE=on

RUN addgroup -S -g 1000 batch-scheduler \
 && adduser -S -u 1000 -G batch-scheduler batch-scheduler

RUN apk update && apk upgrade
RUN apk add bash jq alpine-sdk sed gawk git ca-certificates curl && \
    apk add --no-cache gcc musl-dev && \
    go get -u golang.org/x/lint/golint

#TODO For dev only - remove
WORKDIR /home/user1/go/src/github.com/equinor
RUN git clone -b 183988-batch-job-scheduler https://github.com/equinor/radix-job-scheduler.git

WORKDIR /go/src/github.com/equinor/radix-batch-scheduler/

# get dependencies
COPY go.mod go.sum ./
RUN go mod download

# copy api code
COPY . .

#TODO restore after dev
## lint and unit tests
#RUN golint `go list ./...` && \
#    go vet `go list ./...` && \
#    CGO_ENABLED=0 GOOS=linux go test `go list ./...`

# Build radix api go project
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w" -a -installsuffix cgo -o /usr/local/bin/radix-batch-scheduler

#TODO restore after dev
#FROM scratch
#COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
#COPY --from=builder /etc/passwd /etc/passwd
#COPY --from=builder /usr/local/bin/radix-batch-scheduler /usr/local/bin/radix-batch-scheduler

USER 1000

ENTRYPOINT ["/usr/local/bin/radix-batch-scheduler"]
