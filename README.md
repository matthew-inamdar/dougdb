# DougDB
DougDB is a distributed key-value store making use of the following technologies/patterns:
* Raft for replication.
* SSTables for optimising for high write throughput.
* WAL for fault-tolerance.
* Facade for entrypoint and partition routing on keys to multiple clusters.

This is not intended for production use, but as a fun exercise to learn and experiment.
