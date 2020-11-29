# Master Nodes Section

Similar to the bootstrap section; this sets up master node
configuration. Please note that this is an array.

```yaml
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
```

* `masters.name` - The hostname (WITHOUT the fqdn) of the master node you want to set (x of 3).
* `masters.ipaddr` - The IP address (x of 3) that you want set (this modifies the dhcp config file, the dns zonefile, and the reverse dns zonefile)
* `masters.macaddr` - The mac address for dhcp reservation. This option is not needed if you're doing static ips.
* `masters.disk` - The name of the disk to install RHCOS onto.

> NOTE: 3 Masters are MANDATORY for installation of OpenShift 4


