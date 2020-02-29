#!/bin/bash

curl -i \
  -H "Accept: application/json" \
  -H "Content-Type:application/json" \
  -X POST --data '{"handler":"erc20","cas":["0xdac17f958d2ee523a2206206994597c13d831ec7"], "from":9124875 }' \
  http://localhost:7777/v1/daemon/job/add/factory/eth_subs/jobid/erc20_tether1

curl -i \
  -H "Accept: application/json" \
  -H "Content-Type:application/json" \
  -X POST --data '{"handler":"erc20","cas":["0xdac17f958d2ee523a2206206994597c13d831ec7"]}' \
  http://localhost:7777/v1/daemon/job/add/factory/eth_subs/jobid/erc20_tether2

curl -i \
  -H "Accept: application/json" \
  -H "Content-Type:application/json" \
  -X POST --data '{"handler":"erc20","cas":["0xdac17f958d2ee523a2206206994597c13d831ec7"], "from":9124875 }' \
  http://localhost:7777/v1/daemon/job/add/factory/eth_subs/jobid/erc20_tether3

curl -i \
  -H "Accept: application/json" \
  -H "Content-Type:application/json" \
  -X POST --data '{"handler":"erc20","cas":["0xdac17f958d2ee523a2206206994597c13d831ec7"], "from":9124875 }' \
  http://localhost:7777/v1/daemon/job/add/factory/eth_subs/jobid/erc20_tether4

curl -i \
  -H "Accept: application/json" \
  -H "Content-Type:application/json" \
  -X POST --data '{"handler":"erc20","cas":["0xdac17f958d2ee523a2206206994597c13d831ec7"], "from":9124875 }' \
  http://localhost:7777/v1/daemon/job/add/factory/eth_subs/jobid/erc20_tether5

curl -i \
  -H "Accept: application/json" \
  -H "Content-Type:application/json" \
  -X POST --data '{"handler":"erc20","cas":["0xdac17f958d2ee523a2206206994597c13d831ec7"], "from":9124875 }' \
  http://localhost:7777/v1/daemon/job/add/factory/eth_subs/jobid/erc20_tether6

curl -i \
  -H "Accept: application/json" \
  -H "Content-Type:application/json" \
  -X POST --data '{"handler":"erc20","cas":["0xdac17f958d2ee523a2206206994597c13d831ec7"], "from":9124875 }' \
  http://localhost:7777/v1/daemon/job/add/factory/eth_subs/jobid/erc20_tether7

curl -i \
  -H "Accept: application/json" \
  -H "Content-Type:application/json" \
  -X POST --data '{"handler":"erc20","cas":["0xdac17f958d2ee523a2206206994597c13d831ec7"], "from":9124875 }' \
  http://localhost:7777/v1/daemon/job/add/factory/eth_subs/jobid/erc20_tether8

curl -i \
  -H "Accept: application/json" \
  -H "Content-Type:application/json" \
  -X POST --data '{"handler":"erc20","cas":["0xdac17f958d2ee523a2206206994597c13d831ec7"], "from":9124875 }' \
  http://localhost:7777/v1/daemon/job/add/factory/eth_subs/jobid/erc20_tether9

curl -i \
  -H "Accept: application/json" \
  -H "Content-Type:application/json" \
  -X POST --data '{"handler":"erc20","cas":["0xdac17f958d2ee523a2206206994597c13d831ec7"], "from":9124875 }' \
  http://localhost:7777/v1/daemon/job/add/factory/eth_subs/jobid/erc20_tether10

curl -i \
  -H "Accept: application/json" \
  -H "Content-Type:application/json" \
  -X POST --data '{"handler":"erc20","cas":["0xdac17f958d2ee523a2206206994597c13d831ec7"], "from":9124875 }' \
  http://localhost:7777/v1/daemon/job/add/factory/eth_subs/jobid/erc20_tether11

curl -i \
  -H "Accept: application/json" \
  -H "Content-Type:application/json" \
  -X POST --data '{"handler":"erc20","cas":["0xdac17f958d2ee523a2206206994597c13d831ec7"], "from":9124875 }' \
  http://localhost:7777/v1/daemon/job/add/factory/eth_subs/jobid/erc20_tether12

curl -i \
  -H "Accept: application/json" \
  -H "Content-Type:application/json" \
  -X POST --data '{"handler":"erc20","cas":["0xdac17f958d2ee523a2206206994597c13d831ec7"], "from":9124875 }' \
  http://localhost:7777/v1/daemon/job/add/factory/eth_subs/jobid/erc20_tether13


curl -i \
  -H "Accept: application/json" \
  -H "Content-Type:application/json" \
  -X POST --data '{"source":"erc20_tether13", "dataType":"erc20", "logType":"none" }' \
  http://localhost:7777/v1/daemon/job/add/factory/eth_log/jobid/log_tether13


curl -i \
  -H "Accept: application/json" \
  -H "Content-Type:application/json" \
  -X POST --data '{"handler":"erc20","cas":["0xdac17f958d2ee523a2206206994597c13d831ec7"], "from":9124875 }' \
  http://localhost:7777/v1/daemon/job/add/factory/eth_subs/jobid/erc20_tether14


curl -i \
  -H "Accept: application/json" \
  -H "Content-Type:application/json" \
  -X POST --data '{"source":"erc20_tether14", "dataType":"erc20", "logType":"none" }' \
  http://localhost:7777/v1/daemon/job/add/factory/eth_log/jobid/log_tether14


#curl -i \
#-H "Accept: application/json" \
#-H "Content-Type:application/json" \
#-X POST --data '{"handler":"erc20","cas":["0xB8c77482e45F1F44dE1745F52C74426C631bDD52","0x514910771af9ca656af840dff83e8264ecf986ca"] }' \
#http://localhost:7777/v1/daemon/job/add/factory/eth_subs/jobid/erc20_others

#curl -i \
#-H "Accept: application/json" \
#-H "Content-Type:application/json" \
#-X POST --data '{"handler":"erc721","cas":["0x06012c8cf97bead5deae237070f9587f8e7a266d","0x0e3a2a1f2146d86a604adc220b4967a898d7fe07"] }' \
#http://localhost:7777/v1/daemon/job/add/factory/eth_subs/jobid/erc721_tokens
