# Pushpal-API  [![Build Status](https://travis-ci.com/dernise/pushpal-api.svg?token=AbEANjysKDJ24sgJwcmH&branch=master)](https://travis-ci.com/dernise/pushpal-api)

This is the main repo of the pushpal API

## How to install

Run openssl genrsa -out pushpal.rsa 1024 to generate the private key. Store this key at the root level.

Run openssl rsa -in pushpal.rsa -pubout > pushpal.rsa.pub to generate the public key. Store this key at the root level.
