# Pull

The `pull` subcommand pulls the default "core" images onto the node and
sets up initial `~/.helpernodectl.yaml` config file.

Example:

```shell
helpernodectl pull
```

Currently, the default "core" images being pulled are:

* quay.io/helpernode/pxe
* quay.io/helpernode/http
* quay.io/helpernode/loadbalancer
* quay.io/helpernode/dns
* quay.io/helpernode/dhcp

# Disconnected Pulling

If you are running the HelperNode in a disconnected environment (and
have pulled/tagged/pushed the core images into your registry), you can
use the `HELPERNODE_IMAGE_PREFIX` environment variable to indicate what
registry plus prefix you'd like to use.

Example:

```shell
export HELPERNODE_IMAGE_PREFIX=registry.example.com:5000/mystuff
```

This will result in the images having the prefix of
`registry.example.com:5000/mystuff` and suffix of `/helpernode/<service>`.

For example with the `HELPERNODE_IMAGE_PREFIX` set to
`registry.example.com:5000/mystuff` the `helpernodectl pull` command
will try and pull the following:

* registry.example.com:5000/mystuff/helpernode/pxe
* registry.example.com:5000/mystuff/helpernode/http
* registry.example.com:5000/mystuff/helpernode/loadbalancer
* registry.example.com:5000/mystuff/helpernode/dns
* registry.example.com:5000/mystuff/helpernode/dhcp
