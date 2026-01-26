FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod ./
COPY main.go ./

RUN go build -o openshift-plugin .

FROM alpine:3.19

COPY --from=builder /app/openshift-plugin /usr/local/bin/openshift-plugin
COPY plugin.yaml /home/argocd/cmp-server/config/plugin.yaml

ENTRYPOINT ["/usr/local/bin/openshift-plugin"]
