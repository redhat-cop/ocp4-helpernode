# Bootstrap Node Section

This section defines the bootstrap node configuration.

```yaml
bootstrap:
  name: "bootstrap"
  ipaddr: "192.168.7.20"
  macaddr: "52:54:00:60:72:67"
  disk: vda
```

The options are:

* `bootstrap.name` - The hostname (WITHOUT the fqdn) of the bootstrap node you want to set
* `bootstrap.ipaddr` - The IP address that you want set (this modifies the dhcp config file, the dns zonefile, and the reverse dns zonefile)
* `bootstrap.macaddr` - The mac address for dhcp reservation. This option is not needed if you're doing static ips.
* `bootstrap.disk` - The disk that will be used to install RHCOS onto.

After an install, you probably want to remove the bootstrap from the
loadbalancer. To do this, just omit the `bootstrap` section from the YAML.

> *NOTE* You need to `stop` then `start` the service after this change
> has been made.
