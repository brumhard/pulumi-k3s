#! /bin/sh
# copied from https://docs.cilium.io/en/stable/gettingstarted/k3s/

curl -L --remote-name-all https://github.com/cilium/cilium-cli/releases/latest/download/cilium-linux-amd64.tar.gz{,.sha256sum}
sha256sum --check cilium-linux-amd64.tar.gz.sha256sum
tar xzvfC cilium-linux-amd64.tar.gz /usr/local/bin
rm cilium-linux-amd64.tar.gz cilium-linux-amd64.tar.gz.sha256sum
