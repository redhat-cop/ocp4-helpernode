# DHCP Service Section

This section sets up the DHCP server.

```yaml
dhcp:
  router: "192.168.7.1"
  bcast: "192.168.7.255"
  netmask: "255.255.255.0"
  poolstart: "192.168.7.10"
  poolend: "192.168.7.30"
  ipid: "192.168.7.0"
  netmaskid: "255.255.255.0"
```

Explanation of the options you can set:

* `dhcp.router` - This is the default gateway of your network you're going to assign to the masters/workers
* `dhcp.bcast` - This is the broadcast address for your network
* `dhcp.netmask` - This is the netmask that gets assigned to your masters/workers
* `dhcp.poolstart` - This is the first address in your dhcp address pool
* `dhcp.poolend` - This is the last address in your dhcp address pool
* `dhcp.ipid` - This is the ip network id for the range
* `dhcp.netmaskid` - This is the networkmask id for the range.

> DHCP is OPTIONAL if doing static ips.

These variables are used to set up the dhcp config file. Note, that you need to also set up the `DNS` section if using DHCP, even if you're not using the HelperNode DNS server. This is because the DHCP service uses that configuration fo hand out the DNS information via DHCP.

Just set the `DNS` section to whichever DNS server(s) you're using.
