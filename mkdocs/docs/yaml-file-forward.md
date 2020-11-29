# The helpernode.yaml File

The configuration file for the HelperNode is a YAML file that can
be named anything, but will be referred to in this documentation as
`helpernode.yaml`. This configuration file is used as the "source of
truth" to configure and start the services needed to run the HelperNode.

This config file needs to be either saved (by using the `save` subcommand)
or passed to the `start` command.

```shell
helpernodectl start --config helpernode.yaml
```

When each service starts, the `helpernode.yaml` file is passed to the container and each individual container uses this YAML file to configure and start the services.

> NOTE: Currently, no validation is done by the container. It's "garbage in/garbage out" in it's current form.

# Versioning and Architecture

The current version of the YAML manifest is `v2`. Only `x86_64` is currently supported.

```yaml
version: v2
arch: "x86_64"
```

Plans are in place to support `PPCLE` and `ARM`.

> ARM support will come at after OpenShift supports ARM.
