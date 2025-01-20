# Helper Node email notifications installation

This quickstart will get postfix up and running on the bastion host so that you can test OpenShift's monitoring stack Alertmanager notifications.

~~~
ansible-playbook -e @vars.yaml tasks/main.yml -v -e postfix_install=true
~~~

How to update the alertmanager configuration?

1) extract the alertmanager configuration => 

~~~
oc get secret alertmanager-main -o go-template='{{ index .data "alertmanager.yaml"|base64decode}}' -n openshift-monitoring > alertmanager.yaml
~~~

2) Edit the `alertmanager.yaml` file with your required changes.


3) update the alertmanager configuration =>

~~~
oc create secret generic alertmanager-main --from-file=alertmanager.yaml --dry-run=client -o=yaml -n openshift-monitoring | oc replace secret --filename=-
~~~

---
# Helper Node Alertmanager notification configuration

1) How to send the alerts to the multiple receivers in RHOCP 4? - Red Hat Alertmanager is not sending the critical alerts Alertmanager was configured to send the critical alerts but it's not working Alertmanager is not sending the alerts to all the routes It was configured.
* [link](https://access.redhat.com/solutions/6612991) <= sending alerts to multiple receivers

2) Send dummy alerts to alertmanager in OpenShift 4 - A receiver has to be tested, for instance, with a Critical alert. A real critical alert cannot be forced in production.
* [link](https://access.redhat.com/solutions/6828481) <= sending dummy alerts to alertmanager

3) To test the email receiver by generating a critical alert, run the following. **Please note**, it is imperative to read the links in the previous steps 1,2.

~~~
oc exec alertmanager-main-0 -n openshift-monitoring -- amtool alert add --alertmanager.url http://localhost:9093 alertname=myalarm  --start="2022-03-18T00:00:00-00:00" severity=critical
~~~

4) To verify that you received an alert, on the bastion host run the following.

~~~
cat /var/spool/mail/incoming
~~~

---

Working example of the `alertmanager.yaml` configuration. Replace `bastion.ocp4.example.com` with the actual hostname that you used in the `vars.yaml` => `helper.name`

* Entries that should be left as-is:
  * smtp_smarthost: 'smtp.ocp4.example.com:25'
  * smarthost: 'smtp.ocp4.example.com:25'

~~~
global:
  resolve_timeout: 5m
  smtp_from: openshift@bastion.ocp4.example.com
  smtp_smarthost: 'smtp.ocp4.example.com:25'
  smtp_hello: openshift@bastion.ocp4.example.com
  smtp_require_tls: false
  smtp_auth_username: 'user'
  smtp_auth_password: 'password'
inhibit_rules:
  - equal:
      - namespace
      - alertname
    source_matchers:
      - severity = critical
    target_matchers:
      - severity =~ warning|info
  - equal:
      - namespace
      - alertname
    source_matchers:
      - severity = warning
    target_matchers:
      - severity = info
  - equal:
      - namespace
    source_matchers:
      - alertname = InfoInhibitor
    target_matchers:
      - severity = info
receivers:
  - name: Critical
    email_configs:
      - to: incoming@bastion.ocp4.example.com
        from: openshift@bastion.ocp4.example.com
        smarthost: 'smtp.ocp4.example.com:25'
        hello: ocp4.example.com
        require_tls: false
  - name: Default
  - name: 'null'
  - name: Watchdog
  - name: Warning
    email_configs:
      - to: incoming@bastion.ocp4.example.com
route:
  group_by:
    - namespace
  group_interval: 5m
  group_wait: 30s
  receiver: Default
  repeat_interval: 12h
  routes:
    - matchers:
        - alertname = Watchdog
      receiver: Watchdog
    - matchers:
        - alertname = InfoInhibitor
      receiver: 'null'
    - receiver: Critical
      continue: true
      matchers:
        - severity = critical
    - receiver: Warning
      matchers:
        - severity = warning
~~~