#!/bin/bash

set -e

rm -rf swagger.json
json-refs resolve --filter relative ./spec/root.yml >> swagger.json
rm -rf restapi/
rm -rf models/
swagger generate server -A sendmailserviceproxy -f swagger.json --model-package=models

./scripts/copy_post_code_gen.sh

rm -rf client/
swagger generate client -f swagger.json -A sendmailserviceproxy
