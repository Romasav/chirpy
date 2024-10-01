FROM debian:stable-slim

COPY chirpy /bin/goserver

CMD ["/bin/goserver"]