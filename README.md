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
$ cp var/authgear-sms-gateway.example.yaml var/authgear-sms-gateway.yaml
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
    "app_id": "accounts",
    "to": "+85298765432",
    "body": "Your OTP is 123456"
    "language_tag": "zh-HK",
    "template_name": "verficiation_sms.txt",
    "template_variables": {}
}'
```
