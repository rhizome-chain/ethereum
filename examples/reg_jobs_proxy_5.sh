#!/bin/bash

curl -i \
  -H "Accept: application/json" \
  -H "Content-Type:application/json" \
  -X POST --data '{"handler":"erc20","cas":["0xdac17f958d2ee523a2206206994597c13d831ec7"], "from":9119747 }' \
  http://localhost:7777/v1/daemon/job/add/factory/eth_subs/jobid/erc20_tether1

curl -i \
  -H "Accept: application/json" \
  -H "Content-Type:application/json" \
  -X POST --data '{"source":"erc20_tether1", "dataType":"erc20", "logType":"json" }' \
  http://localhost:7777/v1/daemon/job/add/factory/eth_log/jobid/log_tether1

curl -i \
  -H "Accept: application/json" \
  -H "Content-Type:application/json" \
  -X POST --data '{"handler":"erc20","cas":["0xdac17f958d2ee523a2206206994597c13d831ec7"], "from":9119747}' \
  http://localhost:7777/v1/daemon/job/add/factory/eth_subs/jobid/erc20_tether2

curl -i \
  -H "Accept: application/json" \
  -H "Content-Type:application/json" \
  -X POST --data '{"source":"erc20_tether2", "dataType":"erc20", "logType":"simple" }' \
  http://localhost:7777/v1/daemon/job/add/factory/eth_log/jobid/log_tether2


curl -i \
  -H "Accept: application/json" \
  -H "Content-Type:application/json" \
  -X POST --data '{"handler":"erc20","cas":["0xdac17f958d2ee523a2206206994597c13d831ec7"], "from":9119747 }' \
  http://localhost:7777/v1/daemon/job/add/factory/eth_subs/jobid/erc20_tether3

curl -i \
  -H "Accept: application/json" \
  -H "Content-Type:application/json" \
  -X POST --data '{"source":"erc20_tether3", "dataType":"erc20", "logType":"none" }' \
  http://localhost:7777/v1/daemon/job/add/factory/eth_log/jobid/log_tether3

curl -i \
  -H "Accept: application/json" \
  -H "Content-Type:application/json" \
  -X POST --data '{"handler":"erc20","cas":["0xdac17f958d2ee523a2206206994597c13d831ec7"], "from":9119747 }' \
  http://localhost:7777/v1/daemon/job/add/factory/eth_subs/jobid/erc20_tether4

curl -i \
  -H "Accept: application/json" \
  -H "Content-Type:application/json" \
  -X POST --data '{"source":"erc20_tether4", "dataType":"erc20", "logType":"none" }' \
  http://localhost:7777/v1/daemon/job/add/factory/eth_log/jobid/log_tether4

curl -i \
  -H "Accept: application/json" \
  -H "Content-Type:application/json" \
  -X POST --data '{"handler":"erc20","cas":["0xdac17f958d2ee523a2206206994597c13d831ec7"], "from":9119747 }' \
  http://localhost:7777/v1/daemon/job/add/factory/eth_subs/jobid/erc20_tether5

curl -i \
  -H "Accept: application/json" \
  -H "Content-Type:application/json" \
  -X POST --data '{"source":"erc20_tether5", "dataType":"erc20", "logType":"none" }' \
  http://localhost:7777/v1/daemon/job/add/factory/eth_log/jobid/log_tether5
