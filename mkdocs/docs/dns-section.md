# DNS Service Section

This section sets up your DNS server.

```yaml
dns:
  domain: "example.com"
  clusterid: "ocp4"
  forwarder1: "8.8.8.8"
  forwarder2: "8.8.4.4"
```

Explanation of the DNS variables:

* `dns.domain` - This is what domain the installed DNS server will have. This needs to match what you will put for the `baseDomain` inside the `install-config.yaml` OpenShift installer configuration file.
* `dns.clusterid` - This is what your `clusterid` will be named and needs to match what is in `metadata.name` inside the `install-config.yaml` file.
* `dns.forwarder1` - This will be set up as the DNS forwarder. This is usually one of the corporate (or "upstream") DNS servers.
* `dns.forwarder2` - This will be set up as the second DNS forwarder. This is usually one of the corporate (or "upstream") DNS servers.

The DNS server will be set up using `dns.clusterid` + `dns.domain` as the domain it's serving. In the above example, the helper will be setup to be the SOA for `ocp4.example.com`.

> NOTE: Although you CAN use the helper as your dns server. It's best to have your DNS server delegate the `dns.clusterid + dns.domain domain to` the helper (i.e. Delegate `ocp4.example.com` to the helper)

The DNS section is optional, and only needed if you need to run the DNS server.
