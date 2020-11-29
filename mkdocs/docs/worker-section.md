# Worker Nodes Section

Similar to the master section; this sets up worker node
configuration. Please note that this is an array.

> 🚨 This section is optional if you're installing a compact cluster.

```yaml
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

* `workers.name` - The hostname (WITHOUT the fqdn) of the worker node you want to set
* `workers.ipaddr` - The IP address that you want set (this modifies the dhcp config file, the dns zonefile, and the reverse dns zonefile)
* `workers.macaddr` - The mac address for dhcp reservation. This option is not needed if you're doing static ips.
* `workers.disk` - The name of the disk to install RHCOS onto.

> NOTE: At LEAST 2 workers is needed if you're installing a standard version of OpenShift 4
