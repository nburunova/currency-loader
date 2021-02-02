FROM frolvlad/alpine-glibc

LABEL service="currency-loader"


COPY /bin/currency-loader/currency-loader /usr/lib/currency-loader/currency-loader



EXPOSE 5000

CMD /usr/lib/currency-loader/currency-loader --log-format=json --log-level=debug


