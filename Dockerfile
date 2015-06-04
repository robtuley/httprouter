# requires statically linked go binary to be compiled
FROM scratch

COPY httprouter /httprouter

ENTRYPOINT ["/httprouter"]

EXPOSE 8080