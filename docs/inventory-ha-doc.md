# How to configure inventory for high availability environment

In case while running helpernode playbook from one of the helpernode servers, ensure the rest of the helpernodes are added to the inventory file.

```
[vmhost]
localhost ansible_connection=local
192.168.67.3 ansible_connection=ssh ansible_user=root
```

In case while running helpernode playbook from a remote server, ensure all helpernodes are added to the inventory file.

```
[vmhost]
192.168.67.2 ansible_connection=ssh ansible_user=root
192.168.67.3 ansible_connection=ssh ansible_user=root
```

**NOTE**: Ensure SSH connectivity between all the helpernodes is working fine.
