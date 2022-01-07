FROM chaunceyshannon/cicd-tools:1.0.0 as statik-builder

WORKDIR /app 

COPY files/ files/

RUN statik -src=files '-include=*' -f

# --- 

FROM golang:1.17.3-buster as golang-builder-base

WORKDIR /src

COPY go.* ./
RUN go mod download 

RUN apt update; \
    apt install libpcap-dev -y

# --- 

FROM golang-builder-base as golang-builder

WORKDIR /app

ENV CGO_ENABLED=1

COPY go.* ./
COPY *.go ./
COPY --from=statik-builder /app/statik ./statik/

RUN --mount=type=cache,target=/root/.cache/go-build go build -tags pcap -o run -ldflags " -a -s -w -extldflags '-static'"

# --- 

FROM chaunceyshannon/cicd-tools:1.0.0 as upx-builder

ARG BIN_NAME=run

WORKDIR /app

COPY --from=golang-builder /app/${BIN_NAME} ./

RUN upx -9 ${BIN_NAME}

# --- 

FROM gcr.io/distroless/static:nonroot

WORKDIR /app

COPY --from=upx-builder /app/run /bin/notification-webhook

CMD ["/bin/notification-webhook", "-c", "/app/notification-webhook.ini"]
