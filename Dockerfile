FROM golang:1.12 as build

WORKDIR /tmp/app/

# Golang dependencies to WORKDIR
COPY ./glide.yaml .
COPY ./glide.lock .
# Install package manager "Glide" and install needed dependencies for our exporter
RUN curl https://glide.sh/get | sh \
    && glide install

# Copy Golang src files from context into WORKDIR
WORKDIR /go/src/ms-appcenter-exporter/src
COPY ./src .
# Copy Golang dependencies and build the Golang app
RUN mkdir /app/ \
    && cp -a /tmp/app/vendor ./vendor/ \
    && go build -o /app/ms-appcenter-exporter

#############################################
# FINAL IMAGE
#############################################
FROM alpine
RUN apk add --no-cache \
        libc6-compat \
    	ca-certificates \
    	wget \
    	curl
COPY --from=build /app/ /app/
# Use normal 'non-root' user to execute our app. (Linux starts creating "normal" users at UID 1000)
USER 1000
ENTRYPOINT ["/app/ms-appcenter-exporter"]
