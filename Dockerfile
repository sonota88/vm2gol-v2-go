FROM ubuntu:18.04

RUN apt update \
  && apt install -y --no-install-recommends \
    ca-certificates \
    ruby \
    wget \
  && apt-get clean \
  && rm -rf /var/lib/apt/lists/*

ARG USER
ARG GROUP

RUN groupadd ${USER} \
  && useradd ${USER} -g ${GROUP} -m

USER ${USER}

WORKDIR /tmp

ARG ARCHIVE_FILE=go1.17.1.linux-amd64.tar.gz

RUN wget https://golang.org/dl/${ARCHIVE_FILE} \
  && tar --directory=/home/${USER} -xzf $ARCHIVE_FILE \
  && rm $ARCHIVE_FILE

RUN mkdir /home/${USER}/work

ENV PATH /home/${USER}/go/bin:${PATH}
ENV USER ${USER}

WORKDIR /home/${USER}/work
