## LDV

This is a simplified version of our DagChain system. Specifically,
- All the wallet addresses are maintained by a genesis node (#1 node)
- All the transactions are issued by the genesis node (#1 node), and then are broadcast to the other nodes
- Each time, an address will be selected randomly as the sender, while another address will be selected randomly as 
the receiver
- To simulate the different views of tips due to the network latency, the tips cited by a transaction will not be 
deleted from the tip collection immediately. Instead, these tips will be marked by their `cited` flags and be deleted 
from the tip collection later, either periodically or by a outer cli command
