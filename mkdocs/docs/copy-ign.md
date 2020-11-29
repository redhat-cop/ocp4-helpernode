# Copy Ignition Configs

This command takes ignition configurations from the given directory,
and copies those files into the http contianer. For example:

Example:

```shell
helpernodectl copy-ign --dir=./ocp4/
```

This command must be run on the host that is to be the helpernode. There
is no support for copying the ignition files to an external webserver.

Usage:

```shell
helpernodectl copy-ign --dir /path/to/openshift/install/directory
```
