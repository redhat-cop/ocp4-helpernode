# Sample YAML Configurations

Sample configuration files can be found here. These are samples of most
common use cases.

# Full Stack

A "fullstack" install of the HelperNode services.

```yaml
version: v2
arch: "x86_64"
helper:
  name: "helper"
  ipaddr: "192.168.7.77"
  networkifacename: "ens3"
dns:
  domain: "example.com"
  clusterid: "ocp4"
  forwarder1: "8.8.8.8"
  forwarder2: "8.8.4.4"
dhcp:
  router: "192.168.7.1"
  bcast: "192.168.7.255"
  netmask: "255.255.255.0"
  poolstart: "192.168.7.10"
  poolend: "192.168.7.30"
  ipid: "192.168.7.0"
  netmaskid: "255.255.255.0"
bootstrap:
  name: "bootstrap"
  ipaddr: "192.168.7.20"
  macaddr: "52:54:00:60:72:67"
  disk: vda
masters:
  - name: "master0"
    ipaddr: "192.168.7.21"
    macaddr: "52:54:00:e7:9d:67"
    disk: vda
  - name: "master1"
    ipaddr: "192.168.7.22"
    macaddr: "52:54:00:80:16:23"
    disk: vda
  - name: "master2"
    ipaddr: "192.168.7.23"
    macaddr: "52:54:00:d5:1c:39"
    disk: vda
workers:
  - name: "worker0"
    ipaddr: "192.168.7.11"
    macaddr: "52:54:00:f4:26:a1"
    disk: vda
  - name: "worker1"
    ipaddr: "192.168.7.12"
    macaddr: "52:54:00:82:90:00"
    disk: vda
```
# Static IPs

Here is an example if you're doing a "static ip" install and don't need
all the services.

```yaml
version: v2
arch: "x86_64"
helper:
  name: "helper"
  ipaddr: "192.168.7.77"
  networkifacename: "ens3"
disabledServices:
  - dhcp
  - pxe
dns:
  domain: "example.com"
  clusterid: "ocp4"
  forwarder1: "8.8.8.8"
  forwarder2: "8.8.4.4"
bootstrap:
  name: "bootstrap"
  ipaddr: "192.168.7.20"
masters:
  - name: "master0"
    ipaddr: "192.168.7.21"
  - name: "master1"
    ipaddr: "192.168.7.22"
  - name: "master2"
    ipaddr: "192.168.7.23"
workers:
  - name: "worker0"
    ipaddr: "192.168.7.11"
  - name: "worker1"
    ipaddr: "192.168.7.12"
```

# Compact Static

Here is an example of setting up a "compact" cluster (where you have 3 nodes that act as a master and a worker), not using `pxe` or `dhcp` since a static ip install will be performed.

```yaml
version: v2
arch: "x86_64"
helper:
  name: "helper"
  ipaddr: "192.168.7.77"
  networkifacename: "ens3"
disabledServices:
  - dhcp
  - pxe
dns:
  domain: "example.com"
  clusterid: "ocp4"
  forwarder1: "8.8.8.8"
  forwarder2: "8.8.4.4"
bootstrap:
  name: "bootstrap"
  ipaddr: "192.168.7.20"
masters:
  - name: "master0"
    ipaddr: "192.168.7.21"
  - name: "master1"
    ipaddr: "192.168.7.22"
  - name: "master2"
    ipaddr: "192.168.7.23"
```

# Pluggable Services

A "fullstack" install of the HelperNode services with Pluggable Services.

```yaml
version: v2
arch: "x86_64"
helper:
  name: "helper"
  ipaddr: "192.168.7.77"
  networkifacename: "ens3"
dns:
  domain: "example.com"
  clusterid: "ocp4"
  forwarder1: "8.8.8.8"
  forwarder2: "8.8.4.4"
dhcp:
  router: "192.168.7.1"
  bcast: "192.168.7.255"
  netmask: "255.255.255.0"
  poolstart: "192.168.7.10"
  poolend: "192.168.7.30"
  ipid: "192.168.7.0"
  netmaskid: "255.255.255.0"
bootstrap:
  name: "bootstrap"
  ipaddr: "192.168.7.20"
  macaddr: "52:54:00:60:72:67"
  disk: vda
masters:
  - name: "master0"
    ipaddr: "192.168.7.21"
    macaddr: "52:54:00:e7:9d:67"
    disk: vda
  - name: "master1"
    ipaddr: "192.168.7.22"
    macaddr: "52:54:00:80:16:23"
    disk: vda
  - name: "master2"
    ipaddr: "192.168.7.23"
    macaddr: "52:54:00:d5:1c:39"
    disk: vda
workers:
  - name: "worker0"
    ipaddr: "192.168.7.11"
    macaddr: "52:54:00:f4:26:a1"
    disk: vda
  - name: "worker1"
    ipaddr: "192.168.7.12"
    macaddr: "52:54:00:82:90:00"
    disk: vda
pluggableServices:
  anotherweb:
    image: quay.io/christianh814/test-webserver:latest
    ports:
      - 8899/tcp
      - 8899/udp
  myweb:
    image: quay.io/christianh814/webserver-test:latest
    ports:
      - 8888/tcp
```
