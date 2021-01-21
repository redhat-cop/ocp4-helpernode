# Preflight

The `preflight` subcommand checks for port conflicts, systemd/service
conflicts, and firewall rules on the host and can optionally fix errors
it finds.

Example:

```shell
helpernodectl preflight
```

This will only report conflicts/errors and leaves it up to the user
to fix/reconcile.

You can optionally fix systemd and firewall rules by passing the `--fix-all` flag (EXPERIMENTAL).


Example:

```shell
helpernodectl preflight --fix-all
```

Again, this is experimental and assumes you have `firewalld` running on
a RHEL 8/CentOS 8 server/node/machine/vm.

Currently the firewall rules needed for the HelperNode are:

* 6443/tcp
* 22623/tcp
* 8080/tcp
* 9000/tcp
* 9090/tcp
* 67/udp
* 546/udp
* 53/tcp
* 53/udp
* 80/tcp
* 443/tcp
* 22/tcp
* 69/udp
