FROM debian:buster

RUN apt-get update && \
    apt-get install -y git curl make gcc apt-transport-https ca-certificates gnupg-agent software-properties-common
  
RUN curl -fsSL https://download.docker.com/linux/debian/gpg | apt-key add - && \
    apt-key fingerprint 0EBFCD88 && \
    add-apt-repository \
    "deb [arch=amd64] https://download.docker.com/linux/debian \
$(lsb_release -cs) \
stable" 

RUN    apt-get update && \
    apt-get install -y docker-ce-cli

ENV GALLY=0.0.37

RUN cd /usr/local/bin && \
  curl -LO https://github.com/missena-corp/gally/releases/download/v${GALLY}/gally_${GALLY}_linux_x86_64.tar.gz && \
  tar -xvf gally_${GALLY}_linux_x86_64.tar.gz gally && \
  rm -f gally_${GALLY}_linux_x86_64.tar.gz && \
  chmod +x gally

RUN cd /usr/local/bin && \
  curl -Lo operator-sdk https://github.com/operator-framework/operator-sdk/releases/download/v1.0.0/operator-sdk-v1.0.0-x86_64-linux-gnu && \
  chmod +x operator-sdk

ENV PATH="/usr/local/go/bin:${PATH}"

COPY --from=golang:1.15.0 /usr/local/go/ /usr/local/go/

