# Example to use Ethereum watcher


### Init
```
TMHOME=chainroot1 go run ./examples/. init --chain-id=daemon-chain --node-name=node1
TMHOME=chainroot2 go run ./examples/. init --chain-id=daemon-chain --node-name=node2
TMHOME=chainroot3 go run ./examples/. init --chain-id=daemon-chain --node-name=node3
```
* In chainroot1/config/config.toml file, you should set allow_duplicate_ip true to enable to run multi nodes on the same machine
* Copy config.toml and genesis.json and paste to chainroot2/config, chainroot3/config
* In chainroot2 and chainroot3 /config/config.toml, change proxy_app port, laddr port, prometheus_listen_addr port to avoid port conflict

### check node id
```
TMHOME=chainroot1 go run ./examples/. show_node_id
```

### Run Nodes
``` 
TMHOME=chainroot1 go run ./examples/. start  

export MASTER_ID=$(TMHOME=chainroot1 go run ./examples/. show_node_id)
TMHOME=chainroot2 go run ./examples/. start --p2p.persistent_peers=${MASTER_ID}@127.0.0.1:26656 --daemon.api_addr=0.0.0.0:7778

export MASTER_ID=$(TMHOME=chainroot1 go run ./examples/. show_node_id)
TMHOME=chainroot3 go run ./examples/. start --p2p.persistent_peers=${MASTER_ID}@127.0.0.1:26656 --daemon.api_addr=0.0.0.0:7779

```

### Register Sample Jobs
``` 
./cmd/reg_jobs.sh 
```
### Deregister Sample Jobs
``` 
./cmd/del_jobs.sh 
```

