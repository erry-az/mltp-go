FROM golang:1.22.2-alpine3.19 AS build

WORKDIR /go/src/app

COPY . .

#ARG ssh_prv_key

ARG service_type

#RUN mkdir -p /root/.ssh && \
#    chmod 0700 /root/.ssh && \
#    ssh-keyscan github.com > /root/.ssh/known_hosts
#
#RUN echo "$ssh_prv_key" | base64 --decode > /root/.ssh/id_rsa && \
#    chmod 600 /root/.ssh/id_rsa

RUN go env -w CGO_ENABLED=0 && \
    go env -w GO111MODULE="on" && \
    go env -w GOOS=linux && \
    go env -w GOARCH=amd64 && \
    go env -w GOPRIVATE=github.com/erry-az/*

#RUN git config --global --add url."git@github.com:".insteadOf "https://github.com/"

RUN go mod download

RUN go build -v -o $service_type cmd/$service_type/main.go

FROM alpine:3.19 AS final

ARG service_env

ARG service_type

ENV SERVICE_ENV=$service_env

EXPOSE 4848

COPY --from=build /go/src/app/$service_type /main
COPY --from=build /go/src/app/files/config/$service_env /etc/wallet/config/

CMD ["/main"]
