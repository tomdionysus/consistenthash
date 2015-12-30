// consistenthash is designed to provide consistent hash functionality for distributed systems.
//
// The core type is **ServerNetworkNode**, which respresents a single node in the network - i.e. its unique ID, Host address and distribution.
// **ServerNode** represents the local server node, allowing other ServerNetworkNodes to be registered and deregistered.
// The 'primary' node for any given key can be found by calling the GetNodeFor method on the ServerNode.
//
// For more information on consistent hashing please see the [Wikipedia Article](https://en.wikipedia.org/wiki/Consistent_hashing)
package consistenthash