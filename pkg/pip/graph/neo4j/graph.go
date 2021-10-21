package neo4j

import (
	"encoding/json"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"log"
	"ngac/internal/set"
	"ngac/pkg/config"
	"ngac/pkg/operations"
	"ngac/pkg/pip/graph"
	"strings"
)

var _ graph.Graph = &Graph{}

const (
	node_not_found_msg = "node %s does not exist in the graph"
)

type Graph struct {
	config *config.Config
	driver neo4j.Driver
}

// Accepts the config file's location for Neo4j
func New(cfg string) (*Graph, error) {
	conf, err := config.LoadConfig(cfg)
	if err != nil {
		return nil, err
	}
	ret := new(Graph)
	ret.config = conf
	ret.driver = nil
	return ret, nil
}

func (g *Graph) Start() (err error) {
	if g.driver == nil {
		g.driver, err = neo4j.NewDriver(g.config.Uri, neo4j.BasicAuth(g.config.Username, g.config.Password, ""))
		if err != nil {
			return err
		}
	}

	return nil
}

func (g *Graph) Close() (err error) {
	return g.driver.Close()
}

func (g *Graph) CreatePolicyClass(name string, properties graph.PropertyMap) (*graph.Node, error) {
	if len(name) == 0 {
		return nil, fmt.Errorf("no name was provided when creating a node in the in-memory graph")
	} else if g.Exists(name) {
		return nil, fmt.Errorf("the name %s already exists in the graph", name)
	}

	if properties == nil {
		properties = graph.NewPropertyMap()
	}

	session := g.driver.NewSession(neo4j.SessionConfig{
		AccessMode:   neo4j.AccessModeWrite,
		DatabaseName: g.config.Database,
	})

	jsonStr, err := json.Marshal(properties)
	if err != nil {
		return nil, err
	}
	result, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		records, err := tx.Run(fmt.Sprintf("CREATE (n:%s { name: $name, type: $type, properties: $propertyMap }) RETURN n.name, n.type, apoc.convert.getJsonPropertyMap(n, 'properties') as properties", graph.PC.String()), map[string]interface{}{
			"name":        name,
			"type":        graph.PC.String(),
			"propertyMap": string(jsonStr),
		})
		if err != nil {
			return nil, err
		}
		record, err := records.Single()
		if err != nil {
			return nil, err
		}

		n := &graph.Node{
			Name:       record.Values[0].(string),
			Type:       graph.ToNodeType(record.Values[1].(string)),
			Properties: graph.NewPropertyMap(),
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

	return result.(*graph.Node), nil
}

func (g *Graph) CreateNode(name string, t graph.NodeType, properties graph.PropertyMap, initialParent string, additionalParents ...string) (*graph.Node, error) {
	if t == graph.PC {
		return nil, fmt.Errorf("use CreatePolicyClass to create a policy class node")
	} else if len(name) == 0 {
		return nil, fmt.Errorf("no name was provided when creating a node in the in-memory graph")
	} else if g.Exists(name) {
		return nil, fmt.Errorf("the name %s already exists in the graph", name)
	}

	if properties == nil {
		properties = graph.NewPropertyMap()
	}

	session := g.driver.NewSession(neo4j.SessionConfig{
		AccessMode:   neo4j.AccessModeWrite,
		DatabaseName: g.config.Database,
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

		n := &graph.Node{
			Name:       record.Values[0].(string),
			Type:       graph.ToNodeType(record.Values[1].(string)),
			Properties: graph.NewPropertyMap(),
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

	if err := g.Assign(name, initialParent); err != nil {
		return nil, err
	}

	for _, parent := range additionalParents {
		if err := g.Assign(name, parent); err != nil {
			return nil, err
		}
	}

	return result.(*graph.Node), nil
}

func (g *Graph) UpdateNode(name string, properties graph.PropertyMap) error {
	if name == "" {
		return fmt.Errorf("no name was provided when updating a node in the neo4j graph")
	}
	if !g.Exists(name) || name == "" {
		return fmt.Errorf("node with the name %s could not be found to update", name)
	}

	session := g.driver.NewSession(neo4j.SessionConfig{
		AccessMode:   neo4j.AccessModeWrite,
		DatabaseName: g.config.Database,
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

func (g *Graph) RemoveNode(name string) {
	session := g.driver.NewSession(neo4j.SessionConfig{
		AccessMode:   neo4j.AccessModeWrite,
		DatabaseName: g.config.Database,
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

func (mg *Graph) PolicyClasses() set.Set {
	namesPolicyClasses := set.NewSet()
	nodes := g.Nodes()
	for node := range nodes.Iter() {
		if node.(*graph.Node).Type == graph.PC {
			namesPolicyClasses.Add(node.(*graph.Node).Name)
		}
	}

	return namesPolicyClasses
}

func (g *Graph) Nodes() set.Set {
	session := g.driver.NewSession(neo4j.SessionConfig{
		AccessMode:   neo4j.AccessModeRead,
		DatabaseName: g.config.Database,
	})

	result, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		records, err := tx.Run("MATCH (n) return n.name, n.type, apoc.convert.getJsonPropertyMap(n, 'properties') as properties", nil)
		nodes := set.NewSet()
		for records.Next() {
			n := &graph.Node{
				Name:       records.Record().Values[0].(string),
				Type:       graph.ToNodeType(records.Record().Values[1].(string)),
				Properties: graph.NewPropertyMap(),
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

func (g *Graph) Exists(name string) bool {

	session := g.driver.NewSession(neo4j.SessionConfig{
		AccessMode:   neo4j.AccessModeRead,
		DatabaseName: g.config.Database,
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

func (g *Graph) Node(name string) (*graph.Node, error) {
	session := g.driver.NewSession(neo4j.SessionConfig{
		AccessMode:   neo4j.AccessModeRead,
		DatabaseName: g.config.Database,
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

		n := &graph.Node{
			Name:       record.Values[0].(string),
			Type:       graph.ToNodeType(record.Values[1].(string)),
			Properties: graph.NewPropertyMap(),
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

	return result.(*graph.Node), nil
}

func (g *Graph) NodeFromDetails(t graph.NodeType, properties graph.PropertyMap) (*graph.Node, error) {
	search := g.Search(t, properties).Iterator()
	if !search.HasNext() {
		return nil, fmt.Errorf("a node matching the criteria (%s, %v) does not exist", t.String(), properties)
	}

	return search.Next().(*graph.Node), nil
}

func (g *Graph) Search(t graph.NodeType, properties graph.PropertyMap) set.Set {
	if properties == nil {
		properties = graph.NewPropertyMap()
	}

	results := set.NewSet()
	// iterate over the nodes to find ones that match the search parameters
	for n := range g.Nodes().Iter() {
		node := n.(*graph.Node)
		if node.Type != t && t != graph.NOOP {
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

func (g *Graph) Children(name string) set.Set {
	if !g.Exists(name) {
		log.Fatalf(node_not_found_msg, name)
	}

	session := g.driver.NewSession(neo4j.SessionConfig{
		AccessMode:   neo4j.AccessModeRead,
		DatabaseName: g.config.Database,
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

func (g *Graph) Parents(name string) set.Set {
	if !g.Exists(name) {
		log.Fatalf(node_not_found_msg, name)
	}

	session := g.driver.NewSession(neo4j.SessionConfig{
		AccessMode:   neo4j.AccessModeRead,
		DatabaseName: g.config.Database,
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

func (g *Graph) Assign(child, parent string) error {
	if !g.Exists(child) {
		return fmt.Errorf(node_not_found_msg, child)
	} else if !g.Exists(parent) {
		return fmt.Errorf(node_not_found_msg, parent)
	}

	if g.IsAssigned(child, parent) {
		return fmt.Errorf("%s is already assigned to %s", parent, child)
	}

	c, _ := g.Node(child)
	p, _ := g.Node(parent)

	if err := graph.CheckAssignment(c.Type, p.Type); err != nil {
		return err
	}

	session := g.driver.NewSession(neo4j.SessionConfig{
		AccessMode:   neo4j.AccessModeWrite,
		DatabaseName: g.config.Database,
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

func (g *Graph) Deassign(child, parent string) error {
	if !g.Exists(child) {
		return fmt.Errorf(node_not_found_msg, child)
	} else if !g.Exists(parent) {
		return fmt.Errorf(node_not_found_msg, parent)
	}

	session := g.driver.NewSession(neo4j.SessionConfig{
		AccessMode:   neo4j.AccessModeWrite,
		DatabaseName: g.config.Database,
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

func (g *Graph) IsAssigned(child, parent string) bool {
	session := g.driver.NewSession(neo4j.SessionConfig{
		AccessMode:   neo4j.AccessModeRead,
		DatabaseName: g.config.Database,
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

func (g *Graph) Associate(ua, target string, ops operations.OperationSet) error {
	if !g.Exists(ua) {
		return fmt.Errorf(node_not_found_msg, ua)
	} else if !g.Exists(target) {
		return fmt.Errorf(node_not_found_msg, target)
	}

	uaNode, _ := g.Node(ua)
	targetNode, _ := g.Node(target)

	// check that the association is valid
	if err := graph.CheckAssociation(uaNode.Type, targetNode.Type); err != nil {
		return err
	}

	// if no edge exists create an association
	// if an assignment exists create a new edge for the association
	// if an association exists update it
	session := g.driver.NewSession(neo4j.SessionConfig{
		AccessMode:   neo4j.AccessModeWrite,
		DatabaseName: g.config.Database,
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

func (g *Graph) SourceAssociations(source string) (map[string]operations.OperationSet, error) {
	if !g.Exists(source) {
		return nil, fmt.Errorf(node_not_found_msg, source)
	}

	session := g.driver.NewSession(neo4j.SessionConfig{
		AccessMode:   neo4j.AccessModeRead,
		DatabaseName: g.config.Database,
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

func (g *Graph) Dissociate(ua, target string) error {
	if !g.Exists(ua) {
		return fmt.Errorf(node_not_found_msg, ua)
	} else if !g.Exists(target) {
		return fmt.Errorf(node_not_found_msg, target)
	}

	session := g.driver.NewSession(neo4j.SessionConfig{
		AccessMode:   neo4j.AccessModeWrite,
		DatabaseName: g.config.Database,
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

func (g *Graph) TargetAssociations(target string) (map[string]operations.OperationSet, error) {
	if !g.Exists(target) {
		return nil, fmt.Errorf(node_not_found_msg, target)
	}

	session := g.driver.NewSession(neo4j.SessionConfig{
		AccessMode:   neo4j.AccessModeRead,
		DatabaseName: g.config.Database,
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
func (g *Graph) reset() error {
	session := g.driver.NewSession(neo4j.SessionConfig{
		AccessMode:   neo4j.AccessModeWrite,
		DatabaseName: g.config.Database,
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
