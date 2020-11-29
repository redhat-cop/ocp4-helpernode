# Stopping the HelperNode Services

The `stop` subcommand  will stop the HelperNode containers currently
running.

Example:

```shell
helpernodectl stop
```

# Stopping Individual Services

By default, `stop` will stop all "core" services. You can stop
individual services if, for example, you want the  `http` container to
be shutdown since you don't need it post install.

Example:

```shell
helpernodectl stop http
```

You can stop multiple services by using a comma (`,`) as a delimiter.

Example:

```shell
helpernodectl stop dhcp,pxe
```
