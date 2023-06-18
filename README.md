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
