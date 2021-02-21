package tx

import (
    "fmt"
    "github.com/jtejido/ngac/internal/set"
    "github.com/jtejido/ngac/operations"
    "github.com/jtejido/ngac/pip/graph"
    "sync"
)

type txGraphCommitter interface {
    Committer() Committer
}

type txGraphCommitterImpl struct {
    c Committer
}

func (c *txGraphCommitterImpl) Committer() Committer {
    return c.c
}

type txGraphDeassignCommitter struct {
    c             Committer
    parent, child string
}

func (c *txGraphDeassignCommitter) Committer() Committer {
    return c.c
}

type txGraphDissociateCommitter struct {
    c          Committer
    ua, target string
    operations operations.OperationSet
}

func (c *txGraphDissociateCommitter) Committer() Committer {
    return c.c
}

type TxGraph struct {
    sync.RWMutex
    targetGraph  graph.Graph
    nodes        map[string]*graph.Node
    pcs          set.Set
    assignments  map[string]set.Set
    associations map[string]map[string]operations.OperationSet
    cmds         map[string][]txGraphCommitter
    idx          []string
}

func NewTxGraph(g graph.Graph) *TxGraph {
    ans := new(TxGraph)
    ans.targetGraph = g
    ans.nodes = make(map[string]*graph.Node)
    ans.pcs = set.NewSet()
    ans.assignments = make(map[string]set.Set)
    ans.associations = make(map[string]map[string]operations.OperationSet)
    ans.cmds = make(map[string][]txGraphCommitter)
    ans.idx = make([]string, 0)
    return ans
}

func (tx *TxGraph) CreatePolicyClass(name string, properties graph.PropertyMap) (*graph.Node, error) {
    pc := graph.NewNodeWithFields(name, graph.PC, properties)
    tx.pcs.Add(pc)
    tx.nodes[name] = pc

    _, found := tx.cmds["create_policy_class"]
    if !found {
        tx.cmds["create_policy_class"] = make([]txGraphCommitter, 0)
        tx.idx = append(tx.idx, "create_policy_class")
    }
    tx.cmds["create_policy_class"] = append(tx.cmds["create_policy_class"], &txGraphCommitterImpl{
        c: func() error {
            _, err := tx.targetGraph.CreatePolicyClass(name, properties)
            return err
        }},
    )

    return pc, nil
}

func (tx *TxGraph) CreateNode(name string, t graph.NodeType, properties graph.PropertyMap, initialParent string, additionalParents ...string) (*graph.Node, error) {
    node := graph.NewNodeWithFields(name, t, properties)
    tx.nodes[name] = node

    parents := set.NewSet()
    parents.Add(initialParent)
    for _, p := range additionalParents {
        parents.Add(p)
    }

    // check that the parents exist in the tx or target graph
    for parent := range parents.Iter() {
        _, found := tx.nodes[parent.(string)]
        if !(found || tx.targetGraph.Exists(parent.(string))) {
            return nil, fmt.Errorf("parent %s does not exist", parent.(string))
        }
    }
    _, found := tx.cmds["create_node"]
    if !found {
        tx.cmds["create_node"] = make([]txGraphCommitter, 0)
        tx.idx = append(tx.idx, "create_node")
    }

    tx.cmds["create_node"] = append(tx.cmds["create_node"], &txGraphCommitterImpl{
        c: func() error {
            it := parents.Iterator()
            var ip string
            assert(it.HasNext())
            if it.HasNext() {
                ip = it.Next().(string)
                parents.Remove(ip)
            }
            pts := make([]string, parents.Len())
            var i int
            for p := range parents.Iter() {
                pts[i] = p.(string)
                i++
            }
            _, err := tx.targetGraph.CreateNode(name, t, properties, ip, pts...)
            return err
        }},
    )

    return node, nil
}

func (tx *TxGraph) UpdateNode(name string, properties graph.PropertyMap) (err error) {
    var node *graph.Node

    if v, found := tx.nodes[name]; found {
        node = v
    } else if tx.targetGraph.Exists(name) {
        node, err = tx.targetGraph.Node(name)
        if err != nil {
            return
        }
    } else {
        return fmt.Errorf("node %s does not exist", name)
    }

    node.Properties = properties
    if node.Type == graph.PC {
        tx.pcs.Add(node)
    }

    tx.nodes[name] = node
    _, found := tx.cmds["update_node"]
    if !found {
        tx.cmds["update_node"] = make([]txGraphCommitter, 0)
        tx.idx = append(tx.idx, "update_node")
    }
    tx.cmds["update_node"] = append(tx.cmds["update_node"], &txGraphCommitterImpl{
        c: func() error {
            // node, err := tx.targetGraph.Node(name)
            // if err != nil {
            //     return err
            // }
            // originalProperties := node.Properties
            return tx.targetGraph.UpdateNode(name, properties)
        }},
    )
    return nil
}

func (tx *TxGraph) RemoveNode(name string) {
    if _, found := tx.nodes[name]; found {
        delete(tx.nodes, name)
        tx.pcs.Remove(graph.NewNodeWithoutProps(name, graph.PC))
    }
    _, found := tx.cmds["remove_node"]
    if !found {
        tx.cmds["remove_node"] = make([]txGraphCommitter, 0)
        tx.idx = append(tx.idx, "remove_node")
    }
    tx.cmds["remove_node"] = append(tx.cmds["remove_node"], &txGraphCommitterImpl{
        c: func() error {
            tx.targetGraph.RemoveNode(name)
            return nil
        }},
    )
}

func (tx *TxGraph) Exists(name string) bool {
    _, found := tx.nodes[name]
    return found || tx.targetGraph.Exists(name)
}

func (tx *TxGraph) PolicyClasses() set.Set {
    pcs := set.NewSet()
    for pc := range tx.pcs.Iter() {
        pcs.Add(pc.(*graph.Node).Name)
    }

    pcs.AddFrom(tx.targetGraph.PolicyClasses())

    return pcs
}

func (tx *TxGraph) Nodes() set.Set {
    node := make([]*graph.Node, len(tx.nodes))
    var i int
    for _, v := range tx.nodes {
        node[i] = v
        i++
    }
    nodeSet := set.NewSet()
    for _, nt := range tx.nodes {
        nodeSet.Add(graph.NewNodeFromNode(nt))
    }
    tx.nodes = make(map[string]*graph.Node)
    for n := range tx.targetGraph.Nodes().Iter() {
        nt := n.(*graph.Node)
        tx.nodes[nt.Name] = nt
    }
    for _, nt := range tx.nodes {
        nodeSet.Add(graph.NewNodeFromNode(nt))
    }

    return nodeSet
}

func (tx *TxGraph) Node(name string) (*graph.Node, error) {
    if n, found := tx.nodes[name]; found {
        return graph.NewNodeFromNode(n), nil
    }
    n, err := tx.targetGraph.Node(name)
    if err != nil {
        return nil, err
    }

    return graph.NewNodeFromNode(n), nil
}

func (tx *TxGraph) NodeFromDetails(t graph.NodeType, properties graph.PropertyMap) (*graph.Node, error) {
    search := tx.Search(t, properties)
    if search.Len() == 0 {
        return nil, fmt.Errorf("node with type (%s) with properties %q does not exist", t.String(), properties)
    }
    it := search.Iterator()
    assert(it.HasNext())
    next := it.Next()
    return graph.NewNodeFromNode(next.(*graph.Node)), nil
}

func (tx *TxGraph) Search(t graph.NodeType, properties graph.PropertyMap) set.Set {
    // check tx first
    txNodes := tx.txSearch(t, properties)
    search := tx.targetGraph.Search(t, properties)
    for node := range search.Iter() {
        if _, found := txNodes[node.(*graph.Node).Name]; found {
            continue
        }

        txNodes[node.(*graph.Node).Name] = graph.NewNodeFromNode(node.(*graph.Node))
    }
    s := set.NewSet()
    for _, v := range txNodes {
        s.Add(v)
    }

    return s
}

func (tx *TxGraph) txSearch(t graph.NodeType, properties graph.PropertyMap) map[string]*graph.Node {
    if properties == nil {
        properties = graph.NewPropertyMap()
    }

    results := make(map[string]*graph.Node)
    // iterate over the nodes to find ones that match the search parameters
    for _, node := range tx.nodes {
        // if the type parameter is not null and the current node type does not equal the type parameter, do not add
        if t != graph.ALL && node.Type != t {
            continue
        }

        add := true
        for _, key := range properties.Keys() {
            checkValue, _ := properties.Get(key)
            foundValue, _ := node.Properties.Get(key)
            // if the property provided in the search parameters is null or *, continue to the next property
            if checkValue != foundValue {
                add = false
                break
            }
        }

        if add {
            results[node.Name] = graph.NewNodeFromNode(node)
        }
    }

    return results
}

func (tx *TxGraph) Children(name string) set.Set {
    // get children from the target graph
    children := set.NewSet()
    if tx.targetGraph.Exists(name) {
        children.AddFrom(tx.targetGraph.Children(name))
    }

    // add the children from the tx

    for child, val := range tx.assignments {
        parents := val
        if !parents.Contains(name) {
            continue
        }

        children.Add(child)
    }

    // remove and deassigns
    for _, txCmd := range tx.cmds["deassign"] {
        if txCmd.(*txGraphDeassignCommitter).parent != name {
            continue
        }

        children.Remove(txCmd.(*txGraphDeassignCommitter).child)
    }

    return children
}

func (tx *TxGraph) Parents(name string) set.Set {
    // get children from the target graph
    parents := set.NewSet()
    if tx.targetGraph.Exists(name) {
        parents.AddFrom(tx.targetGraph.Parents(name))
    }

    // add the children from the tx
    if v, found := tx.assignments[name]; found {
        parents.AddFrom(v)
    }
    // remove and deassigns
    for _, txCmd := range tx.cmds["deassign"] {
        if txCmd.(*txGraphDeassignCommitter).child != name {
            continue
        }

        parents.Remove(txCmd.(*txGraphDeassignCommitter).child)
    }

    return parents
}

func (tx *TxGraph) Assign(child, parent string) error {
    parents := set.NewSet()
    if v, found := tx.assignments[child]; found {
        parents.AddFrom(v)
    }
    parents.Add(parent)
    tx.assignments[child] = parents
    _, found := tx.cmds["assign"]
    if !found {
        tx.cmds["assign"] = make([]txGraphCommitter, 0)
        tx.idx = append(tx.idx, "assign")
    }

    tx.cmds["assign"] = append(tx.cmds["assign"], &txGraphCommitterImpl{
        c: func() error {
            return tx.targetGraph.Assign(child, parent)
        }},
    )

    return nil
}

func (tx *TxGraph) Deassign(child, parent string) error {
    parents := set.NewSet()
    if v, found := tx.assignments[child]; found {
        parents.AddFrom(v)
    }
    parents.Remove(parent)
    tx.assignments[child] = parents
    _, found := tx.cmds["deassign"]
    if !found {
        tx.cmds["deassign"] = make([]txGraphCommitter, 0)
        tx.idx = append(tx.idx, "deassign")
    }

    tx.cmds["deassign"] = append(tx.cmds["deassign"], &txGraphDeassignCommitter{
        c: func() error {
            return tx.targetGraph.Deassign(child, parent)
        },
        child:  child,
        parent: parent,
    })

    return nil
}

func (tx *TxGraph) IsAssigned(child, parent string) bool {
    res := set.NewSet()
    if v, found := tx.assignments[child]; found {
        res.AddFrom(v)
    }
    return res.Contains(parent) || tx.targetGraph.IsAssigned(child, parent)
}

func (tx *TxGraph) Associate(ua, target string, ops operations.OperationSet) error {
    assocs := make(map[string]operations.OperationSet)
    if v, found := tx.associations[ua]; found {
        assocs = v
    }

    assocs[target] = ops
    tx.associations[ua] = assocs

    _, found := tx.cmds["associate"]
    if !found {
        tx.cmds["associate"] = make([]txGraphCommitter, 0)
        tx.idx = append(tx.idx, "associate")
    }

    tx.cmds["associate"] = append(tx.cmds["associate"], &txGraphCommitterImpl{
        c: func() error {
            return tx.targetGraph.Associate(ua, target, ops)
        },
    })

    return nil
}

func (tx *TxGraph) Dissociate(ua, target string) error {
    assocs := make(map[string]operations.OperationSet)
    if v, found := tx.associations[ua]; found {
        assocs = v
    }
    delete(assocs, target)
    tx.associations[ua] = assocs
    _, found := tx.cmds["dissociate"]
    if !found {
        tx.cmds["dissociate"] = make([]txGraphCommitter, 0)
        tx.idx = append(tx.idx, "dissociate")
    }

    tx.cmds["dissociate"] = append(tx.cmds["dissociate"], &txGraphDissociateCommitter{
        c: func() error {
            // assocs, err := tx.targetGraph.SourceAssociations(ua)
            // if err != nil {
            //     return err
            // }
            // operations = assocs.Get(target);
            return tx.targetGraph.Dissociate(ua, target)
        },
        ua:     ua,
        target: target,
    })
    return nil
}

func (tx *TxGraph) SourceAssociations(source string) (map[string]operations.OperationSet, error) {
    // get target graph associations
    sourceAssociations, err := tx.targetGraph.SourceAssociations(source)
    if err != nil {
        return nil, err
    }
    // get tx associations
    if v, found := tx.associations[source]; found {
        for key, val := range v {
            sourceAssociations[key] = val
        }
    }

    // remove any dissociates
    for _, txCmd := range tx.cmds["dissociate"] {
        if txCmd.(*txGraphDissociateCommitter).ua != source {
            continue
        }

        delete(sourceAssociations, txCmd.(*txGraphDissociateCommitter).target)
    }

    return sourceAssociations, nil

}

func (tx *TxGraph) TargetAssociations(target string) (map[string]operations.OperationSet, error) {
    // get target graph associations
    targetAssociations, err := tx.targetGraph.TargetAssociations(target)
    if err != nil {
        return nil, err
    }
    // get tx associations
    for _, val := range tx.associations {
        assocs := val
        if _, ok := assocs[target]; !ok {
            continue
        }

        os := operations.NewOperationSet()
        v, found := assocs[target]
        if found {
            os.AddFrom(v)
        }
        targetAssociations[target] = os
    }

    // remove any dissociates
    for _, txCmd := range tx.cmds["dissociate"] {
        if txCmd.(*txGraphDissociateCommitter).ua != target {
            continue
        }

        delete(targetAssociations, txCmd.(*txGraphDissociateCommitter).target)
    }

    return targetAssociations, nil
}

func (tx *TxGraph) Commit() (err error) {
    tx.Lock()
    for _, v := range tx.idx {
        for _, txCmd := range tx.cmds[v] {
            f := txCmd.Committer()
            if err = f(); err != nil {
                return
            }
        }
    }
    tx.Unlock()
    return nil
}
