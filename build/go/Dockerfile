FROM alpine:3.12

LABEL os="ubuntu-16.04" \
      service="currency-loader"

RUN apk add --no-cache ca-certificates

COPY bin/currency-loader/currency-loader /usr/lib/currency-loader/currency-loader

EXPOSE 5000

CMD /usr/lib/currency-loader/currency-loader --log-format=json -- log-level=error


