# ngac

**Next Generation Access Control**

This is a Golang port of NIST's reference core implementation from Policy Machine (written in Java).

[https://github.com/PM-Master/policy-machine-core](https://github.com/PM-Master/policy-machine-core)

**Find their documentation here:**

[https://pm-master.github.io/pm-master/policy-machine-core/](https://pm-master.github.io/pm-master/policy-machine-core/)

**TO-DO**

At the moment, it's translated verbatim. Once all functions are translated, re-factoring will begin (which includes locking mechanisms on structs/fields, memoizing here and there, and multi-thread stuff when necessary).

EPP shall be a PubSub hub model.

Obligation file will be json (following the original's json schema) as golang don't have yaml official support (and third party libs shall be avoided if possible) in their toolchain.

Neo4J (Persisted) and MemGraph (In-memory) support - since both uses Cypher. (DAO for different GraphDB had to be implemented)