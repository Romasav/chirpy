FROM debian:stable-slim

COPY chirpy /bin/chirpy

ENV CHIRPY_DB_PATH=/app/data/database.json

CMD ["/bin/chirpy"]