FROM alpine:3.23.2
RUN apk update && apk add --upgrade libcrypto3 libssl3
RUN adduser -u 10000 -D -g '' sbomreport-to-dependencytrack sbomreport-to-dependencytrack

COPY sbomreport-to-dependencytrack /usr/local/bin/sbomreport-to-dependencytrack

USER 10000
EXPOSE 8080
ENTRYPOINT ["sbomreport-to-dependencytrack"]
