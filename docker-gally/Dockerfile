FROM docker:19.06

ENV GALLY=0.0.37

RUN apk add  --no-cache git curl && \
  cd /usr/local/bin && \
  curl -LO https://github.com/missena-corp/gally/releases/download/v${GALLY}/gally_${GALLY}_linux_x86_64.tar.gz && \
  tar -xvf gally_${GALLY}_linux_x86_64.tar.gz gally && \
  chmod +x gally

