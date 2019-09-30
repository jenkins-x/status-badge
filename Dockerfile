FROM scratch
EXPOSE 8080
ENTRYPOINT ["/kaniko-test"]
COPY ./bin/ /