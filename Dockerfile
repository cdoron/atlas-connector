FROM registry.access.redhat.com/ubi8/ubi-minimal:8.6

ENV HOME=/tmp
WORKDIR /tmp

COPY atlas-connector /atlas-connector
USER 10001
EXPOSE 8080

ENTRYPOINT ["/atlas-connector"]
CMD [ "run" ]
