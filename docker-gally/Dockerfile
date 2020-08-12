FROM docker:19.03

ENV GALLY=0.0.37

RUN apk add  --no-cache git curl make && \
  cd /usr/local && \
  curl -LO https://golang.org/dl/go1.15.linux-amd64.tar.gz && \
  tar -xzf go1.15.linux-amd64.tar.gz && \
  rm -f go1.15.linux-amd64.tar.gz && \
  cd /usr/local/bin && \
  curl -LO https://github.com/missena-corp/gally/releases/download/v${GALLY}/gally_${GALLY}_linux_x86_64.tar.gz && \
  tar -xvf gally_${GALLY}_linux_x86_64.tar.gz gally && \
  rm -f gally_${GALLY}_linux_x86_64.tar.gz && \
  chmod +x gally

ENV PATH="/usr/local/go/bin:${PATH}"

