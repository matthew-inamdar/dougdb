# DougDB
DougDB is a distributed key-value store making use of the following technologies/patterns:
* Raft for replication.
* SSTables for optimising for high write throughput.
* WAL for fault-tolerance.
* Facade for entrypoint and partition routing on keys to multiple clusters.
* Synchronous consistency and eventual consistency option (eventual consistency to allow reads from followers).

This is not intended for production use, but as a fun exercise to learn and experiment.

## Todo
* Hydrating log from startup should apply committed entries to state machine.
* Entries need to have defined structure and checksum to validate good writes.

## Future work
* Cluster membership changes.

# Notes
* Need 5 nodes in the cluster
* The `Doug` gRPC service is what clients will communicate with
* The `Doug` gRPC service will issue RPCs to other node's `Raft` gRPC service
* The `Doug` gRPC service is what will redirect clients to the current leader
* A candidate node can time out and start a new election an increment the term
* Leader's need to send periodic heartbeats with empty entries if no activity
* Followers need to trigger leader election if no entry received from leader for some time - called `election timeout`
* Beginning an election: 
  * Transitions the follower to candidate state and increments current term
  * Candidate votes for itself and issued RequestVote RPc in parallel to all nodes
  * Candidate will remain in this state until either:
    * It wins the election
    * Another server establishes itself as leader
    * The election times out
* A candidate wins if it receives a majority for the same term
* A node can only vote for a single candidate per term
* Incrementing a term must nullify the current vote
* Election timeouts are randomised to mitigate split brain elections
* Candidates and followers both have election timeouts (also used for missed heartbeats)
* We'll use 150-300ms as the election timeout random bounds
* When a candidate becomes leader, it initialises all `nextIndex` values to the
index after the latest entry in the new leader's log
* 
