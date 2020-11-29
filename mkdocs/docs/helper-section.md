# Helper Section

This section sets up the variables for the server/vm/node where the HelperNode services will run.

```yaml
helper:
  name: "helper"
  ipaddr: "192.168.7.77"
  networkifacename: "ens3"
```

This is how it breaks down

* `helper.name` - REQUIRED: This needs to be set to the hostname you want your helper to be (some people leave it as "helper" others change it to "bastion"). This will create an entry in the DNS service.
* `helper.ipaddr` - REQUIRED Set this to the current IP address of the helper. This is used to set up the reverse dns definition of not only the HelperNode itself, but also the in-addr.arpa  configuration in DNS.
* `helper.networkifacename` - REQUIRED: This is set to the interface that has the `helper.ipaddr` ip address.

> NOTE: The `helper.networkifacename` is the ACTUAL name of the interface, NOT the NetworkManager name (you should NEVER need to set it to something like `System eth0`. Set it to what you see when you run the `ip addr` command)
