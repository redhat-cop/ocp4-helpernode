# Get Clients

This will get the  needed clients from the holding container. Which is,
by default, the http container. It saves these in your current working
directory.

Example:

```shell
helpernodectl get-clients
```

Currently, the clients that are provided are:

* helm
* kubectl
* oc
* openshift-install

These are provided via a compressed tarball (i.e. `tar.gz`). It is left
to the user to extract them in the proper `$PATH` location.
