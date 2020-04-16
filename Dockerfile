FROM ubuntu

COPY extender_example /usr/bin/extender_example
RUN chmod +x /usr/bin/extender_example

ENTRYPOINT ["/usr/bin/extender_example"]