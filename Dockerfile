FROM scratch
EXPOSE 8080
ENTRYPOINT ["/status-badge"]
COPY ./bin/ /