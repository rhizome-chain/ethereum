#!/bin/bash

curl -i -H "Accept: application/json" -H "Content-Type:application/json" \
  -X DELETE http://localhost:7777/v1/daemon/job/erc20_tether1

curl -i -H "Accept: application/json" -H "Content-Type:application/json" \
  -X DELETE http://localhost:7777/v1/daemon/job/log_tether1

curl -i -H "Accept: application/json" -H "Content-Type:application/json" \
  -X DELETE http://localhost:7777/v1/daemon/job/erc20_tether2

curl -i -H "Accept: application/json" -H "Content-Type:application/json" \
  -X DELETE http://localhost:7777/v1/daemon/job/log_tether2

curl -i -H "Accept: application/json" -H "Content-Type:application/json" \
  -X DELETE http://localhost:7777/v1/daemon/job/erc20_tether3

curl -i -H "Accept: application/json" -H "Content-Type:application/json" \
  -X DELETE http://localhost:7777/v1/daemon/job/log_tether3

curl -i -H "Accept: application/json" -H "Content-Type:application/json" \
  -X DELETE http://localhost:7777/v1/daemon/job/erc20_tether4

curl -i -H "Accept: application/json" -H "Content-Type:application/json" \
  -X DELETE http://localhost:7777/v1/daemon/job/log_tether4

curl -i -H "Accept: application/json" -H "Content-Type:application/json" \
  -X DELETE http://localhost:7777/v1/daemon/job/erc20_tether5

curl -i -H "Accept: application/json" -H "Content-Type:application/json" \
  -X DELETE http://localhost:7777/v1/daemon/job/log_tether5
