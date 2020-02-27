# Ethereum watcher on Rhizome Chain (tenermint-daemon)


### Register Ethereum Subscription Job : REST Call
```bash
curl -i \
  -H "Accept: application/json" \
  -H "Content-Type:application/json" \
  -X POST --data '{"handler":"{handler_name}","cas":["{contract addresses}"], "from":{block number} }' \
  http://{daemon_api_address}/v1/daemon/job/add/factory/eth_subs/jobid/{job_id_to_register}


curl -i \
  -H "Accept: application/json" \
  -H "Content-Type:application/json" \
  -X POST --data '{"handler":"erc20","cas":["0xdac17f958d2ee523a2206206994597c13d831ec7"], "from":9124875 }' \
  http://localhost:7777/v1/daemon/job/add/factory/eth_subs/jobid/erc20_tether1
```

### Job Description
```json
{
  "handler":"erc20|erc721",
  "cas":["{contract addresses}"], 
  "from":"{block number to subscribe from. If null, from latest block}" 
}
```