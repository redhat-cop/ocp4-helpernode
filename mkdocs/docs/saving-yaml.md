# Saving YAML Config

The `save` subcommand is used to save a configuration file needed to
start the HelperNode services. Once saved, you no longer have to pass
`--config` to the `start` subcommand.

Example:

```shell
helpernodectl save -f helpernode.yaml
```

This will save the provided config file to ` ~/.helper.yaml`
