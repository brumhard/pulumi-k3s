FROM ubuntu

RUN apt-get update && apt-get install -y openssh-server
RUN useradd -m -s /bin/bash sshuser
COPY ./testdata/id_rsa.pub /home/sshuser/.ssh/authorized_keys
EXPOSE 22
ENTRYPOINT service ssh start && while true; do sleep 30; done;