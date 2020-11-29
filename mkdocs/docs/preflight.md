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
a RHEL/CentOS 8 server/node/machine/vm.
