# Anomaly generation for Open5GS

This directory contains tooling and example configurations to generate
abnormal call flows with `gnbsim`. The goal is to exercise error handling
paths on an Open5GS 5G core by injecting malformed or unexpected signalling
messages.

## Available scenarios

* `malformed-nas` – sends a NAS message with random bytes which should be
  rejected by the core.
* `service-before-registration` – issues a Service Request before completing
  the Registration procedure.
* `duplicate-registration` – transmits two Registration Requests back to back.
* `auth-response-before-request` – emits an Authentication Response without a
  preceding challenge.
* `security-complete-before-command` – sends Security Mode Complete before the
  network issues Security Mode Command.
* `deregistration-before-registration` – attempts a Deregistration prior to any
  Registration.
* `empty-nas` – forwards a NAS message with a zero-length payload.

## Running an anomaly

1. Adapt the gNB/AMF addresses in the example configuration to match your
   deployment.
2. Execute `gnbsim` with the chosen configuration, e.g.:

   ```sh
   gnbsim --cfg anomaly/malformed-nas.yaml
   ```

   or

   ```sh
   gnbsim --cfg anomaly/duplicate-registration.yaml
   ```

Logs produced by Open5GS during these runs can be collected as anomalous
samples for analysis.
