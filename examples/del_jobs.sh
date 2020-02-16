#!/bin/bash

curl -i -H "Accept: application/json" -H "Content-Type:application/json" \
     -X DELETE http://localhost:7777/v1/daemon/job/erc20_tether
curl -i -H "Accept: application/json" -H "Content-Type:application/json" \
     -X DELETE http://localhost:7777/v1/daemon/job/erc20_others
curl -i -H "Accept: application/json" -H "Content-Type:application/json" \
     -X DELETE http://localhost:7777/v1/daemon/job/erc721_tokens