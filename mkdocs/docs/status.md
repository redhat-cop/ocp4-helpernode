# Displaying Status

The `status` subcommand shows the status of the running containers on
the host.

Example:

```shell
helpernodectl status
```

Sample Output:

```shell
Names                     Status              Image
helpernode-pxe            Up 46 minutes ago   registry.example.com:5000/alpha/helpernode/pxe:latest
helpernode-loadbalancer   Up 46 minutes ago   registry.example.com:5000/alpha/helpernode/loadbalancer:latest
helpernode-http           Up 46 minutes ago   registry.example.com:5000/alpha/helpernode/http:latest
helpernode-dhcp           Up 46 minutes ago   registry.example.com:5000/alpha/helpernode/dhcp:latest
helpernode-dns            Up 46 minutes ago   registry.example.com:5000/alpha/helpernode/dns:latest
```

This simply passes you the information provided by the contianer runtime.
