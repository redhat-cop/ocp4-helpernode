# Starting the HelperNode Services

The `start` subcommand  will start the containers needed for the HelperNode to run. It will run, and configure, the services depending on what YAML configuration is passed.

Example:

```shell
helpernodectl start --config helpernode.yaml
```

The `--config` needs to be passed, unless you performed a `save` of the
config file (which saves it under `~/.helper.yaml`).

# Preflight Checks

A preflight check is performed when you issue the `start` command. This
ensures that there is no conflicts before starting the services.

This will sometimes produce a "false positive" when making changes and
wanting to restart the services. To skip the preflight steps during
startup by passing `--skip-preflight` or  the shorthand `-s`.

Example:

```shell
helpernodectl start --config helpernode.yaml --skip-preflight
```

This will start the services without checking for errors first. Use
with caution.

# Starting Individual Services

By default, `start` will start all "core" services. You can start
individua services if, for example, you stopped the `http` container to
make a config change in your YAML file; and need to start it up again.

Example:

```shell
helpernodectl start --config helpernode.yaml http
```

You can start multiple services by using a comma (`,`) as a delimiter.

Example:

```shell
helpernodectl start --config helpernode.yaml http,pxe
```

The same caveat of preflight checks applies. You can pass `--skip-preflight` (or the `-s` shorthand) to skip these checks.

Example:

```shell
helpernodectl start --config helpernode.yaml dhcp,pxe --skip-preflight
```
