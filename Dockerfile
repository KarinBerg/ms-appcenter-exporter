FROM golang:1.12 as build

# golang deps
WORKDIR /tmp/app/
COPY ./src/glide.yaml /tmp/app/
COPY ./src/glide.lock /tmp/app/
RUN curl https://glide.sh/get | sh \
    && glide install

WORKDIR /go/src/ms-appcenter-exporter/src
COPY ./src /go/src/ms-appcenter-exporter/src
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
USER 1000

CMD ["/app/ms-appcenter-exporter"]
