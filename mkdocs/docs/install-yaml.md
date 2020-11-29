# Install Bootstrapping

The `install` subcommand pulls the default "core" images onto the node
and sets up initial `~/.helpernodectl.yaml` config file.

Example:

```shell
helpernodectl install
```

Currently, the default "core" images being pulled are:

* quay.io/helpernode/pxe
* quay.io/helpernode/http
* quay.io/helpernode/loadbalancer
* quay.io/helpernode/dns
* quay.io/helpernode/dhcp
