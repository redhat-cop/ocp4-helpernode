# Restarting the HelperNode Services

The `restart` subcommand  will stop then start running HelperNode
sevices. The services that are restrated depend on what is in your YAML
file that you've saved. (see [saving yaml](saving-yaml.md) for more information)

Example:

```shell
helpernodectl restart
```

# Restarting Individual Services

By default, `restart` will stop then start all services listed in your YAML. You can restart
individual services if, for example, you need to restart the http container.

Example:

```shell
helpernodectl restart http
```

You can restart multiple services by using a comma (`,`) as a delimiter.

Example:

```shell
helpernodectl restart http,pxe
```
