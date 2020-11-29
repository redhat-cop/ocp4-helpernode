# Disabled Services Section

By default, the HelperNode CLI utility starts all "core" services. If
you're only using a subset of services, you can specify the ones you
don't want to start.

```yaml
disabledServices:
  - dhcp
  - pxe
```

Current list of "core" services:

* `pxe`
* `http`
* `loadbalancer`
* `dns`
* `dhcp`
