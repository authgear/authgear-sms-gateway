# Authgear SMS Gateway

Authgear SMS Gateway is a HTTP server that recieves a send SMS request and invoke send action in corresponding provider.
The request will be redirected to corresponding provider by a set of rules.

## Pre-requisite

Install tools specified in `.tool-versions`

### Install dependencies

```sh
$ make vendor
```

### Set up environment

```sh
$ cp .env.example .env
$ cp var/sms_service_provider_config.example.yaml var/sms_service_provider_config.yaml
```

## Run

```
$ make start
```

## Example

### Send

```sh
$ curl --request POST \
  --url http://localhost:8091/send \
  --header 'Content-Type: application/json' \
  --data '{
    "to": "+85298765432",
    "body": "Your OTP is 123456"
}'
```
