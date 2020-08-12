FROM docker:19.03

ENV GALLY=0.0.37

RUN apk add  --no-cache git curl make && \
  cd /usr/local/bin && \
  curl -LO https://github.com/missena-corp/gally/releases/download/v${GALLY}/gally_${GALLY}_linux_x86_64.tar.gz && \
  tar -xvf gally_${GALLY}_linux_x86_64.tar.gz gally && \
  rm -f gally_${GALLY}_linux_x86_64.tar.gz && \
  chmod +x gally

ENV PATH="/usr/local/go/bin:${PATH}"

COPY --from=golang:1.15.0-alpine3.12 /usr/local/go/ /usr/local/go/

