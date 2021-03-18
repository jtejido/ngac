# ngac

## Next Generation Access Control

This is a Golang port of NIST's reference core implementation from Policy Machine (written in Java).

[https://github.com/PM-Master/policy-machine-core](https://github.com/PM-Master/policy-machine-core)

This port supports **Neo4j** as our Persistent Graph. In order to run it, it will require the [APOC Core](https://neo4j.com/labs/apoc/4.1/installation/) plugin to be installed. The config file is located [here](https://github.com/jtejido/ngac/configs) and [this](https://github.com/jtejido/ngac/scripts) Cypher script can be ran to quickly serve the config's requirements.

## Find their documentation here:

[https://pm-master.github.io/pm-master/policy-machine-core/](https://pm-master.github.io/pm-master/policy-machine-core/)

## TO-DO

Obligation JSON Unmarshallers - file will be JSON (following the original's JSON schema).

Follow [https://github.com/golang-standards/project-layout](https://github.com/golang-standards/project-layout)
