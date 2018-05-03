FROM kontena/cli

COPY bin/binary-alpine /usr/local/bin/kontena-git

RUN chmod +x /usr/local/bin/kontena-git

ENTRYPOINT []
