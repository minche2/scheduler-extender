FROM ubuntu

COPY first_extender /usr/bin/first_extender
RUN chmod +x /usr/bin/first_extender

ENTRYPOINT ["/usr/bin/first_extender"]
