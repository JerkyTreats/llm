# LLM

This project manages this Local LLM API.

## Network Architecture

This project is intended to be running in a Tailscale network.

* Contains multiple devices with different roles- NAS, GPU enabled PC, mobile phone, etc.
* Assumed to be managed by https://github.com/jerkytreats/dns
* The endpoint dns.{{internal_domain}}/swagger will respond with openapi

## API Documentation

This project includes an auto-generated OpenAPI specification for the API. The OpenAPI spec is automatically generated from the Go code using reflection and route registration.

## Initialization

This project is concerned with the initialization of local LLMs.

It will:

* Setup mount path on NAS
  * API to be assumed to be running on a separate device
  * Device