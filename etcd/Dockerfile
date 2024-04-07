FROM ubuntu:22.04

RUN apt-get update
RUN apt-get install -y \
    htop \
    ranger \
    etcd \
    vim

# setup etcd RBAC
# TODO

# TODO
# figure out how to get the config file to work

EXPOSE 2379

#ENV HostIP="192.168.2.148"

CMD etcd --advertise-client-urls 'http://0.0.0.0:2379' --listen-client-urls 'http://0.0.0.0:2379'
