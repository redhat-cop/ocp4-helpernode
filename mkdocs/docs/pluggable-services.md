# Pluggable Services

You are able to launch other containers that are not part of the "core"
container group. This can be done by adding `pluggableService` to your
config. Below is an example:

```yaml
pluggableServices:
  anotherweb:
    image: quay.io/christianh814/test-webserver:latest
    ports:
      - 8899/tcp
      - 8899/udp
  myweb:
    image: quay.io/christianh814/webserver-test:latest
    ports:
      - 8888/tcp
```

You will need:

* `pluggableServices.<name>` - Name your service something. You must name it something not in the "reserved names". For example you **CANNOT** name your service `loadbalancer` as that's reserved.
* `pluggableServices.image` - This is the registry where the image is stored. You must pass the tag as well.
* `pluggableServices.ports` - This is an array of ports with their protocol.

> **NOTE**: Binding to local storage is not supported currently, but it is planned.
