# Disabled Services Section

By default, the HelperNode services assumes you will use the Load Balancer
functionality. In the case of BareMetal IPI, vSphere IPI, OpenStack IPI,
and RHV IPI; the load balancing is done on the platform itself. You can
flag this in the YAML file.

```yaml
ipiip: 
  api: "192.168.7.200"
  apps: "192.168.7.201"
```

> NOTE: These IPs must **NOT** be in use, as they'll
> be assigned by the OpenShift installer itself. Please see the [official documentation](https://docs.openshift.com/container-platform/4.6/installing/installing_bare_metal_ipi/ipi-install-prerequisites.html#network-requirements_ipi-install-prerequisites) for more details.


Doing this, you'll no longer need the `loadbalancer` service. You can disable this in the YAML as well.

```yaml
disabledServices:
  - loadbalancer
```

> You can also disable `pxe` as most IPI installers don't need it. This
> depends on the IPI implementation. Again, please consult the official
> documentation.
