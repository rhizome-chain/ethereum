# Example to use Ethereum watcher

## Run daemons on the single node
### Init 
```
TMHOME=chainroot1 go run ./examples/. init --chain-id=daemon-chain --node-name=node1 \
       --force-rewrite=true \
       --p2p.allow_duplicate_ip=true \
       --mempool.size=300000 \
       --mempool.max_txs_bytes=107374182400 \
       --mempool.max_tx_bytes=104857600 \
       --consensus.timeout_commit=1s \
       --daemon_api_addr=0.0.0.0:7777 \
       --daemon_alive_threshold=4 \
       --eth_network_url=wss://mainnet.infura.io/ws/v3/21f6dd123d1845b489095a9afa39a81d

export MASTER_ID=$(TMHOME=chainroot1 go run ./examples/. show_node_id)
TMHOME=chainroot2 go run ./examples/. init --chain-id=daemon-chain --node-name=node2 \
       --force-rewrite=true \
       --rpc.laddr=tcp://127.0.0.1:16657 \
       --p2p.allow_duplicate_ip=true \
       --p2p.laddr="tcp://0.0.0.0:16656" \
       --p2p.persistent_peers=${MASTER_ID}@127.0.0.1:26656 \
       --mempool.size=300000 \
       --mempool.max_txs_bytes=107374182400 \
       --mempool.max_tx_bytes=104857600 \
       --instrumentation.prometheus_listen_addr=:16660 \
       --daemon_api_addr=0.0.0.0:7778 \
       --daemon_alive_threshold=4 \
       --eth_network_url=wss://mainnet.infura.io/ws/v3/21f6dd123d1845b489095a9afa39a81d

export MASTER_ID=$(TMHOME=chainroot1 go run ./examples/. show_node_id)
TMHOME=chainroot3 go run ./examples/. init --chain-id=daemon-chain --node-name=node3 \
       --force-rewrite=true \
       --rpc.laddr=tcp://127.0.0.1:17757 \
       --p2p.allow_duplicate_ip=true \
       --p2p.laddr="tcp://0.0.0.0:17756" \
       --p2p.persistent_peers=${MASTER_ID}@127.0.0.1:26656 \
       --mempool.size=300000 \
       --mempool.max_txs_bytes=107374182400 \
       --mempool.max_tx_bytes=104857600 \
       --instrumentation.prometheus_listen_addr=:17760 \
       --daemon_api_addr=0.0.0.0:7779 \
       --daemon_alive_threshold=4 \
       --eth_network_url=wss://mainnet.infura.io/ws/v3/21f6dd123d1845b489095a9afa39a81d

export MASTER_ID=$(TMHOME=chainroot1 go run ./examples/. show_node_id)
TMHOME=chainroot4 go run ./examples/. init --chain-id=daemon-chain --node-name=node4 \
       --force-rewrite=true \
       --rpc.laddr=tcp://127.0.0.1:18857 \
       --p2p.allow_duplicate_ip=true \
       --p2p.laddr="tcp://0.0.0.0:18856" \
       --p2p.persistent_peers=${MASTER_ID}@127.0.0.1:26656 \
       --mempool.size=300000 \
       --mempool.max_txs_bytes=107374182400 \
       --mempool.max_tx_bytes=104857600 \
       --instrumentation.prometheus_listen_addr=:18860 \
       --daemon_api_addr=0.0.0.0:7780 \
       --daemon_alive_threshold=4 \
       --eth_network_url=wss://mainnet.infura.io/ws/v3/21f6dd123d1845b489095a9afa39a81d

```
* Copy genesis.json and paste to chainroot2/config, chainroot3/config

### Run Nodes
``` 
TMHOME=chainroot1 go run ./examples/. start 
TMHOME=chainroot2 go run ./examples/. start 
TMHOME=chainroot3 go run ./examples/. start 
TMHOME=chainroot4 go run ./examples/. start

```

### Register Sample Jobs
``` 
./cmd/reg_jobs.sh 
```
### Deregister Sample Jobs
``` 
./cmd/del_jobs.sh 
```

