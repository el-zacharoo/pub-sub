FROM golang:alpine AS builder

ARG BUF_USER
ARG BUF_PAT
ARG GH_USER
ARG GH_PAT
ARG GH_ORG

# create an unprivileged user
RUN adduser \    
    --disabled-password \    
    --gecos "" \    
    --home "/nonexistent" \    
    --shell "/sbin/nologin" \    
    --no-create-home \    
    --uid 10001 \    
    appuser

# get CA certs required for HTTPS. Git is required for dependencies.
RUN apk update  && apk upgrade && apk add --no-cache git

# set working directory /src (default dir is /go)
WORKDIR /src
COPY . .

# access to private repos
RUN echo "machine github.com login $GH_USER password $GH_PAT" >> ~/.netrc
RUN echo "machine go.buf.build login $BUF_USER password $BUF_PAT" >> ~/.netrc
RUN go env -w GOPRIVATE="github.com/$GH_ORG/*"

# build as static-linked binary (no external dependencies).
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o /app

# build minimal image (800mb -> 15Mb)
FROM scratch

COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app /app

# perform any further action as an unprivileged user
USER appuser:appuser
ENTRYPOINT ["/app"]