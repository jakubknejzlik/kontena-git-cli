FROM kontena/cli

COPY bin/kontena-git-alpine /usr/local/bin/kontena-git

RUN chmod +x /usr/local/bin/kontena-git

ENTRYPOINT []
