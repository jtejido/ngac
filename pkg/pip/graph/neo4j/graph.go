package neo4j

import (
	"encoding/json"
	"fmt"
	"github.com/jtejido/ngac/internal/set"
	"github.com/jtejido/ngac/pkg/config"
	"github.com/jtejido/ngac/pkg/operations"
	g "github.com/jtejido/ngac/pkg/pip/graph"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"log"
	"strings"
)

var _ g.Graph = &graph{}

const (
	node_not_found_msg = "node %s does not exist in the graph"
)

type graph struct {
	config *config.Config
	driver neo4j.Driver
}

// Accepts the config file's location for Neo4j
func New(cfg string) (g.Graph, error) {
	conf, err := config.LoadConfig(cfg)
	if err != nil {
		return nil, err
	}
	ret := new(graph)
	ret.config = conf
	ret.driver = nil
	return ret, nil
}

func (ng *graph) Start() (err error) {
	if ng.driver == nil {
		ng.driver, err = neo4j.NewDriver(ng.config.Uri, neo4j.BasicAuth(ng.config.Username, ng.config.Password, ""))
		if err != nil {
			return err
		}
	}

	return nil
}

func (ng *graph) Close() (err error) {
	return ng.driver.Close()
}

func (ng *graph) CreatePolicyClass(name string, properties g.PropertyMap) (*g.Node, error) {
	if len(name) == 0 {
		return nil, fmt.Errorf("no name was provided when creating a node in the in-memory graph")
	} else if ng.Exists(name) {
		return nil, fmt.Errorf("the name %s already exists in the graph", name)
	}

	if properties == nil {
		properties = g.NewPropertyMap()
	}

	session := ng.driver.NewSession(neo4j.SessionConfig{
		AccessMode:   neo4j.AccessModeWrite,
		DatabaseName: ng.config.Database,
	})

	jsonStr, err := json.Marshal(properties)
	if err != nil {
		return nil, err
	}
	result, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		records, err := tx.Run(fmt.Sprintf("CREATE (n:%s { name: $name, type: $type, properties: $propertyMap }) RETURN n.name, n.type, apoc.convert.getJsonPropertyMap(n, 'properties') as properties", g.PC.String()), map[string]interface{}{
			"name":        name,
			"type":        g.PC.String(),
			"propertyMap": string(jsonStr),
		})
		if err != nil {
			return nil, err
		}
		record, err := records.Single()
		if err != nil {
			return nil, err
		}

		n := &g.Node{
			Name:       record.Values[0].(string),
			Type:       g.ToNodeType(record.Values[1].(string)),
			Properties: g.NewPropertyMap(),
		}

		if m, ok := record.Values[2].(map[string]interface{}); ok {
			for key, value := range m {
				switch value := value.(type) {
				case string:
					n.Properties[key] = value
				default:
					return nil, fmt.Errorf("Illegal type for property value found")
				}

			}
			return n, nil
		}

		return nil, fmt.Errorf("Invalid property")
	})

	session.Close()
	if err != nil {
		return nil, err
	}

	return result.(*g.Node), nil
}

func (ng *graph) CreateNode(name string, t g.NodeType, properties g.PropertyMap, initialParent string, additionalParents ...string) (*g.Node, error) {
	if t == g.PC {
		return nil, fmt.Errorf("use CreatePolicyClass to create a policy class node")
	} else if len(name) == 0 {
		return nil, fmt.Errorf("no name was provided when creating a node in the in-memory graph")
	} else if ng.Exists(name) {
		return nil, fmt.Errorf("the name %s already exists in the graph", name)
	}

	if properties == nil {
		properties = g.NewPropertyMap()
	}

	session := ng.driver.NewSession(neo4j.SessionConfig{
		AccessMode:   neo4j.AccessModeWrite,
		DatabaseName: ng.config.Database,
	})

	jsonStr, err := json.Marshal(properties)
	if err != nil {
		return nil, err
	}
	result, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		records, err := tx.Run(fmt.Sprintf("CREATE (n:%s { name: $name, type: $type, properties: $propertyMap }) RETURN n.name, n.type, apoc.convert.getJsonPropertyMap(n, 'properties') as properties", t.String()), map[string]interface{}{
			"name":        name,
			"type":        t.String(),
			"propertyMap": string(jsonStr),
		})

		if err != nil {
			return nil, err
		}
		record, err := records.Single()
		if err != nil {
			return nil, err
		}

		n := &g.Node{
			Name:       record.Values[0].(string),
			Type:       g.ToNodeType(record.Values[1].(string)),
			Properties: g.NewPropertyMap(),
		}

		if m, ok := record.Values[2].(map[string]interface{}); ok {
			for key, value := range m {
				switch value := value.(type) {
				case string:
					n.Properties[key] = value
				default:
					return nil, fmt.Errorf("Illegal type for property value found")
				}

			}
			return n, nil
		}

		return nil, fmt.Errorf("Invalid property")
	})

	session.Close()
	if err != nil {
		return nil, err
	}

	if err := ng.Assign(name, initialParent); err != nil {
		return nil, err
	}

	for _, parent := range additionalParents {
		if err := ng.Assign(name, parent); err != nil {
			return nil, err
		}
	}

	return result.(*g.Node), nil
}

func (ng *graph) UpdateNode(name string, properties g.PropertyMap) error {
	if name == "" {
		return fmt.Errorf("no name was provided when updating a node in the neo4j graph")
	}
	if !ng.Exists(name) || name == "" {
		return fmt.Errorf("node with the name %s could not be found to update", name)
	}

	session := ng.driver.NewSession(neo4j.SessionConfig{
		AccessMode:   neo4j.AccessModeWrite,
		DatabaseName: ng.config.Database,
	})

	if properties != nil {
		jsonStr, err := json.Marshal(properties)
		if err != nil {
			return err
		}
		_, err = session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
			result, err := tx.Run("MATCH (n{name: $name}) SET n.properties = $propertyMap return n", map[string]interface{}{
				"name":        name,
				"propertyMap": string(jsonStr),
			})

			for result.Next() {

			}

			if err = result.Err(); err != nil {
				return nil, err
			}

			return nil, nil
		})

		session.Close()
		return err
	}

	return nil
}

func (ng *graph) RemoveNode(name string) {
	session := ng.driver.NewSession(neo4j.SessionConfig{
		AccessMode:   neo4j.AccessModeWrite,
		DatabaseName: ng.config.Database,
	})
	session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run("MATCH (n{name: $name}) DETACH DELETE n", map[string]interface{}{
			"name": name,
		})
		if err != nil {
			log.Println(err.Error())
		}
		if _, err = result.Consume(); err != nil {
			log.Println(err.Error())
		}
		return nil, nil
	})

	session.Close()
}

func (ng *graph) PolicyClasses() set.Set {
	namesPolicyClasses := set.NewSet()
	nodes := ng.Nodes()
	for node := range nodes.Iter() {
		if node.(*g.Node).Type == g.PC {
			namesPolicyClasses.Add(node.(*g.Node).Name)
		}
	}

	return namesPolicyClasses
}

func (ng *graph) Nodes() set.Set {
	session := ng.driver.NewSession(neo4j.SessionConfig{
		AccessMode:   neo4j.AccessModeRead,
		DatabaseName: ng.config.Database,
	})

	result, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		records, err := tx.Run("MATCH (n) return n.name, n.type, apoc.convert.getJsonPropertyMap(n, 'properties') as properties", nil)
		nodes := set.NewSet()
		for records.Next() {
			n := &g.Node{
				Name:       records.Record().Values[0].(string),
				Type:       g.ToNodeType(records.Record().Values[1].(string)),
				Properties: g.NewPropertyMap(),
			}

			if m, ok := records.Record().Values[2].(map[string]interface{}); ok {
				for key, value := range m {
					switch value := value.(type) {
					case string:
						n.Properties[key] = value
					default:
						return nil, fmt.Errorf("Illegal type for property value found")
					}

				}
				nodes.Add(n)
			}
		}

		if err = records.Err(); err != nil {
			return nil, err
		}

		return nodes, nil
	})

	session.Close()
	nodes := set.NewSet()
	if err != nil {
		return nodes
	}

	return nodes.Union(result.(set.Set))
}

func (ng *graph) Exists(name string) bool {

	session := ng.driver.NewSession(neo4j.SessionConfig{
		AccessMode:   neo4j.AccessModeRead,
		DatabaseName: ng.config.Database,
	})

	result, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		records, err := tx.Run("MATCH (n{name:$name}) return count(*) as count", map[string]interface{}{
			"name": name,
		})

		if err != nil {
			return nil, err
		}
		record, err := records.Single()
		if err != nil {
			return nil, err
		}

		return record.Values[0], nil
	})

	session.Close()

	if err != nil {
		log.Println(err.Error())
		return false
	}

	count := result.(int64)
	return count > 0
}

func (ng *graph) Node(name string) (*g.Node, error) {
	session := ng.driver.NewSession(neo4j.SessionConfig{
		AccessMode:   neo4j.AccessModeRead,
		DatabaseName: ng.config.Database,
	})

	result, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		records, err := tx.Run("MATCH (n{name:$name}) return n.name, n.type, apoc.convert.getJsonPropertyMap(n, 'properties') as properties", map[string]interface{}{
			"name": name,
		})

		if err != nil {
			return nil, err
		}
		record, err := records.Single()
		if err != nil {
			return nil, err
		}

		n := &g.Node{
			Name:       record.Values[0].(string),
			Type:       g.ToNodeType(record.Values[1].(string)),
			Properties: g.NewPropertyMap(),
		}

		if m, ok := record.Values[2].(map[string]interface{}); ok {
			for key, value := range m {
				switch value := value.(type) {
				case string:
					n.Properties[key] = value
				default:
					return nil, fmt.Errorf("Illegal type for property value found")
				}

			}
			return n, nil
		}

		return nil, fmt.Errorf("Invalid property")
	})

	session.Close()

	if err != nil {
		return nil, err
	}

	return result.(*g.Node), nil
}

func (ng *graph) NodeFromDetails(t g.NodeType, properties g.PropertyMap) (*g.Node, error) {
	search := ng.Search(t, properties).Iterator()
	if !search.HasNext() {
		return nil, fmt.Errorf("a node matching the criteria (%s, %v) does not exist", t.String(), properties)
	}

	return search.Next().(*g.Node), nil
}

func (ng *graph) Search(t g.NodeType, properties g.PropertyMap) set.Set {
	if properties == nil {
		properties = g.NewPropertyMap()
	}

	results := set.NewSet()
	// iterate over the nodes to find ones that match the search parameters
	for n := range ng.Nodes().Iter() {
		node := n.(*g.Node)
		if node.Type != t && t != g.NOOP {
			continue
		}

		match := true
		for k, v := range properties {
			if node.Properties[k] != v {
				match = false
			}
		}

		if match {
			results.Add(node)
		}
	}

	return results
}

func (ng *graph) Children(name string) set.Set {
	if !ng.Exists(name) {
		log.Fatalf(node_not_found_msg, name)
	}

	session := ng.driver.NewSession(neo4j.SessionConfig{
		AccessMode:   neo4j.AccessModeRead,
		DatabaseName: ng.config.Database,
	})

	result, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		records, err := tx.Run("Match (n{name:$name})<-[:ASSIGNED_TO]-(b) return b.name", map[string]interface{}{
			"name": name,
		})
		children := set.NewSet()
		for records.Next() {
			children.Add(records.Record().Values[0].(string))
		}

		if err = records.Err(); err != nil {
			return nil, err
		}

		return children, nil
	})

	session.Close()
	nodes := set.NewSet()
	if err != nil {
		return nodes
	}

	return nodes.Union(result.(set.Set))
}

func (ng *graph) Parents(name string) set.Set {
	if !ng.Exists(name) {
		log.Fatalf(node_not_found_msg, name)
	}

	session := ng.driver.NewSession(neo4j.SessionConfig{
		AccessMode:   neo4j.AccessModeRead,
		DatabaseName: ng.config.Database,
	})

	result, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		records, err := tx.Run("Match (n{name:$name})-[:ASSIGNED_TO]->(b) return b.name", map[string]interface{}{
			"name": name,
		})
		parents := set.NewSet()
		for records.Next() {
			parents.Add(records.Record().Values[0].(string))
		}

		if err = records.Err(); err != nil {
			return nil, err
		}

		return parents, nil
	})

	session.Close()
	nodes := set.NewSet()
	if err != nil {
		return nodes
	}

	return nodes.Union(result.(set.Set))
}

func (ng *graph) Assign(child, parent string) error {
	if !ng.Exists(child) {
		return fmt.Errorf(node_not_found_msg, child)
	} else if !ng.Exists(parent) {
		return fmt.Errorf(node_not_found_msg, parent)
	}

	if ng.IsAssigned(child, parent) {
		return fmt.Errorf("%s is already assigned to %s", parent, child)
	}

	c, _ := ng.Node(child)
	p, _ := ng.Node(parent)

	if err := g.CheckAssignment(c.Type, p.Type); err != nil {
		return err
	}

	session := ng.driver.NewSession(neo4j.SessionConfig{
		AccessMode:   neo4j.AccessModeWrite,
		DatabaseName: ng.config.Database,
	})

	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		_, err := tx.Run("MATCH (a {name:$child}), (b {name:$parent}) MERGE (a)-[:ASSIGNED_TO]->(b)", map[string]interface{}{
			"child":  child,
			"parent": parent,
		})

		return nil, err
	})

	session.Close()
	return err
}

func (ng *graph) Deassign(child, parent string) error {
	if !ng.Exists(child) {
		return fmt.Errorf(node_not_found_msg, child)
	} else if !ng.Exists(parent) {
		return fmt.Errorf(node_not_found_msg, parent)
	}

	session := ng.driver.NewSession(neo4j.SessionConfig{
		AccessMode:   neo4j.AccessModeWrite,
		DatabaseName: ng.config.Database,
	})

	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		_, err := tx.Run("MATCH (a {name:$child})-[r:ASSIGNED_TO]->(b {name:$parent}) DELETE r", map[string]interface{}{
			"child":  child,
			"parent": parent,
		})

		return nil, err
	})

	session.Close()
	return err
}

func (ng *graph) IsAssigned(child, parent string) bool {
	session := ng.driver.NewSession(neo4j.SessionConfig{
		AccessMode:   neo4j.AccessModeRead,
		DatabaseName: ng.config.Database,
	})

	result, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		records, err := tx.Run("MATCH (a{name:$child})-[:ASSIGNED_TO]->(b{name:$parent}) return count(*) as count", map[string]interface{}{
			"child":  child,
			"parent": parent,
		})

		if err != nil {
			return nil, err
		}
		record, err := records.Single()
		if err != nil {
			return nil, err
		}

		return record.Values[0], nil
	})

	session.Close()

	if err != nil {
		log.Println(err.Error())
		return false
	}

	count := result.(int64)
	return count > 0
}

func (ng *graph) Associate(ua, target string, ops operations.OperationSet) error {
	if !ng.Exists(ua) {
		return fmt.Errorf(node_not_found_msg, ua)
	} else if !ng.Exists(target) {
		return fmt.Errorf(node_not_found_msg, target)
	}

	uaNode, _ := ng.Node(ua)
	targetNode, _ := ng.Node(target)

	// check that the association is valid
	if err := g.CheckAssociation(uaNode.Type, targetNode.Type); err != nil {
		return err
	}

	// if no edge exists create an association
	// if an assignment exists create a new edge for the association
	// if an association exists update it
	session := ng.driver.NewSession(neo4j.SessionConfig{
		AccessMode:   neo4j.AccessModeWrite,
		DatabaseName: ng.config.Database,
	})

	opsStr := make([]string, ops.Len())
	var i int
	for op := range ops.Iter() {
		opsStr[i] = op.(string)
		i++
	}
	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		semiformat := fmt.Sprintf("%q\n", opsStr)
		tokens := strings.Split(semiformat, " ")
		result, err := tx.Run(fmt.Sprintf(
			`
				match (ua{name:$ua})
				match (target{name:$target})
				optional match (ua)-[r:ASSOCIATION]->(target)
				WITH coalesce(r) as r1
				CALL apoc.do.when(r1 is null,
				'match (ua{name:"%s"}), (target{name:"%s"}) merge (ua)-[r:ASSOCIATION]->(target) set r.operations = %v RETURN ua, target',
				'match (ua{name:"%s"})-[r:ASSOCIATION]->(target{name:"%s"}) set r.operations = %v RETURN ua, target') YIELD value
				return value.ua, value.target
			`, ua, target, strings.Join(tokens, ", "), ua, target, strings.Join(tokens, ", ")), map[string]interface{}{
			"ua":     ua,
			"target": target,
		})
		if err != nil {
			return nil, err
		}

		for result.Next() {

		}

		return nil, result.Err()
	})

	session.Close()
	return err
}

func (ng *graph) SourceAssociations(source string) (map[string]operations.OperationSet, error) {
	if !ng.Exists(source) {
		return nil, fmt.Errorf(node_not_found_msg, source)
	}

	session := ng.driver.NewSession(neo4j.SessionConfig{
		AccessMode:   neo4j.AccessModeRead,
		DatabaseName: ng.config.Database,
	})

	result, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		records, err := tx.Run("Match (ua{name:$source})-[r:ASSOCIATION]->(target) return target.name, r.operations", map[string]interface{}{
			"source": source,
		})
		assocs := make(map[string]operations.OperationSet)
		for records.Next() {
			if v, ok := records.Record().Values[1].([]interface{}); ok {
				assocs[records.Record().Values[0].(string)] = operations.NewOperationSet(v...)
			}
		}

		if err = records.Err(); err != nil {
			return nil, err
		}

		return assocs, nil
	})

	session.Close()
	if err != nil {
		return nil, err
	}

	return result.(map[string]operations.OperationSet), nil
}

func (ng *graph) Dissociate(ua, target string) error {
	if !ng.Exists(ua) {
		return fmt.Errorf(node_not_found_msg, ua)
	} else if !ng.Exists(target) {
		return fmt.Errorf(node_not_found_msg, target)
	}

	session := ng.driver.NewSession(neo4j.SessionConfig{
		AccessMode:   neo4j.AccessModeWrite,
		DatabaseName: ng.config.Database,
	})

	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		_, err := tx.Run("MATCH (a {name:$ua})-[r:ASSOCIATION]->(b {name:$target}) DELETE r", map[string]interface{}{
			"ua":     ua,
			"target": target,
		})

		return nil, err
	})

	session.Close()
	return err
}

func (ng *graph) TargetAssociations(target string) (map[string]operations.OperationSet, error) {
	if !ng.Exists(target) {
		return nil, fmt.Errorf(node_not_found_msg, target)
	}

	session := ng.driver.NewSession(neo4j.SessionConfig{
		AccessMode:   neo4j.AccessModeRead,
		DatabaseName: ng.config.Database,
	})

	result, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		records, err := tx.Run("Match (target{name:$target})<-[r:ASSOCIATION]-(ua) return ua.name, r.operations", map[string]interface{}{
			"target": target,
		})
		assocs := make(map[string]operations.OperationSet)
		for records.Next() {
			if v, ok := records.Record().Values[1].([]interface{}); ok {
				assocs[records.Record().Values[0].(string)] = operations.NewOperationSet(v...)
			}
		}

		if err = records.Err(); err != nil {
			return nil, err
		}

		return assocs, nil
	})

	session.Close()
	if err != nil {
		return nil, err
	}

	return result.(map[string]operations.OperationSet), nil
}

// testing only
func (ng *graph) reset() error {
	session := ng.driver.NewSession(neo4j.SessionConfig{
		AccessMode:   neo4j.AccessModeWrite,
		DatabaseName: ng.config.Database,
	})
	defer session.Close()

	result, err := session.Run("match(n) detach delete n", map[string]interface{}{})
	if err != nil {
		return err
	}

	if _, err = result.Consume(); err != nil {
		return err
	}

	return nil
}
