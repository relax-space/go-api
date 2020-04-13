FROM pangpanglabs/golang:builder-beta AS builder
RUN go env -w GOPROXY=https://goproxy.cn,direct
WORKDIR /go/src/go-api
COPY . .
RUN go mod download
# disable cgo
ENV CGO_ENABLED=0

# build steps
RUN go build -o go-api

# make application docker image use alpine
FROM pangpanglabs/alpine-ssl
WORKDIR /go/bin/
# copy config files to image, if nomni use: replace WORKDIR with /go/src/nomni/go-api
COPY --from=builder /go/src/go-api/*.yml ./
# copy execute file to image
COPY --from=builder /go/src/go-api/go-api ./
EXPOSE 8080
CMD ["./go-api"]

