# ngac

## Next Generation Access Control

This is a Golang port of NIST's reference core implementation, Policy Machine.

[https://github.com/PM-Master/policy-machine-core](https://github.com/PM-Master/policy-machine-core)

This port supports **Neo4j** as our Persistent Graph. In order to run it, it will require the [APOC Core](https://neo4j.com/labs/apoc/4.1/installation/) plugin to be installed. The config file is located [here](https://github.com/jtejido/github.com/jtejido/ngac/tree/master/configs) and [this](https://github.com/jtejido/github.com/jtejido/ngac/tree/master/scripts) Cypher script can be ran to quickly serve the config's requirements.

## Find their documentation here:

[https://pm-master.github.io/pm-master/policy-machine-core/](https://pm-master.github.io/pm-master/policy-machine-core/)

## TO-DO

Be reminded that this is **!!NOT FOR PROD!!** as the APIs are still open for changes.

Obligation JSON Unmarshallers - file will be JSON (following the original's JSON schema).

Follow [https://github.com/golang-standards/project-layout](https://github.com/golang-standards/project-layout)

DTO/DAO models for various Persistent and In-Memory graph DBs.

EPP to Publish/Subscribe model
