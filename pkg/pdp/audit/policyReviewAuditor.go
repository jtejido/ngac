package audit

import (
    "fmt"
    "github.com/jtejido/ngac/internal/set"
    "github.com/jtejido/ngac/pkg/operations"
    "github.com/jtejido/ngac/pkg/pip/graph"
)

type PReviewAuditor struct {
    graph       graph.Graph
    resourceOps operations.OperationSet
}

func NewPReviewAuditor(graph graph.Graph, resourceOps operations.OperationSet) *PReviewAuditor {
    if graph == nil {
        panic("graph cannot be nil")
    }

    if resourceOps == nil {
        panic("resourceOps cannot be nil")
    }
    return &PReviewAuditor{graph, resourceOps}
}

func (pa *PReviewAuditor) Explain(userID, target string) (*Explain, error) {
    userNode, err := pa.graph.Node(userID)
    if err != nil {
        return nil, err
    }
    targetNode, err := pa.graph.Node(target)
    if err != nil {
        return nil, err
    }

    userPaths, err := pa.dfs(userNode)
    if err != nil {
        return nil, err
    }
    targetPaths, err := pa.dfs(targetNode)
    if err != nil {
        return nil, err
    }

    resolvedPaths := pa.resolvePaths(userPaths, targetPaths, target)
    perms := pa.resolvePermissions(resolvedPaths)

    return NewExplain(perms, resolvedPaths), nil
}

func (pa *PReviewAuditor) resolvePermissions(paths map[string]*PolicyClass) set.Set {
    pcPerms := make(map[string]set.Set)
    for pc, pcPaths := range paths {
        for p := range pcPaths.Paths.Iter() {
            ops := p.(*Path).Operations
            var exOps set.Set
            if v, ok := pcPerms[pc]; ok {
                exOps = v
            } else {
                exOps = set.NewSet()
            }

            exOps.AddFrom(ops)
            pcPerms[pc] = exOps
        }
    }

    perms := set.NewSet()
    var first bool = true

    for _, ops := range pcPerms {
        if first {
            perms.AddFrom(ops)
            first = false
        } else {
            if perms.Contains(operations.ALL_OPS) {
                perms.Remove(operations.ALL_OPS)
                perms.AddFrom(ops)
            } else {
                // if the ops for the pc are empty then the user has no permissions on the target
                if ops.Len() == 0 {
                    perms.Clear()
                    break
                } else if !ops.Contains(operations.ALL_OPS) {
                    perms.RetainFrom(ops)
                }
            }
        }
    }

    // remove any unknown ops
    perms.RetainFrom(pa.resourceOps)

    // if the permission set includes *, remove the * and add all resource operations
    if perms.Contains(operations.ALL_OPS) {
        perms.Remove(operations.ALL_OPS)
        perms.AddFrom(operations.AdminOps())
        perms.AddFrom(pa.resourceOps)
    } else {
        // if the permissions includes *a or *r add all the admin ops/resource ops as necessary
        if perms.Contains(operations.ALL_ADMIN_OPS) {
            perms.Remove(operations.ALL_OPS)
            perms.AddFrom(operations.AdminOps())
        } else if perms.Contains(operations.ALL_RESOURCE_OPS) {
            perms.Remove(operations.ALL_RESOURCE_OPS)
            perms.AddFrom(pa.resourceOps)
        }
    }

    return perms
}

/**
 * Given a set of paths starting at a user, and a set of paths starting at an object, return the paths from
 * the user to the target node (through an association) that belong to each policy class. A path is added to a policy
 * class' entry in the returned map if the user path ends in an association in which the target of the association
 * exists in a target path. That same target path must also end in a policy class. If the path does not end in a policy
 * class the target path is ignored.
 *
 * @param userPaths the set of paths starting with a user.
 * @param targetPaths the set of paths starting with a target node.
 * @param target the name of the target node.
 * @return the set of paths from a user to a target node (through an association) for each policy class in the system.
 * @throws PMException if there is an exception traversing the graph
 */
func (pa *PReviewAuditor) resolvePaths(userPaths, targetPaths []*edgePath, target string) map[string]*PolicyClass {
    results := make(map[string]*PolicyClass)
    for _, targetPath := range targetPaths {
        pcEdge := targetPath.edges[len(targetPath.edges)-1]

        // if the last element in the target path is a pc, the target belongs to that pc, add the pc to the results
        // skip to the next target path if it is not a policy class
        if pcEdge.target.Type != graph.PC {
            continue
        }

        var policyClass *PolicyClass
        if v, ok := results[pcEdge.target.Name]; ok {
            policyClass = v
        } else {
            policyClass = NewEmptyPolicyClass()
        }

        // compute the paths for this target path
        paths := pa.computePaths(userPaths, targetPath, target, pcEdge)

        // add all paths
        existingPaths := policyClass.Paths
        existingPaths.AddFrom(paths)

        // collect all ops
        for p := range paths.Iter() {
            policyClass.Operations.AddFrom(p.(*Path).Operations)
        }

        // update results
        results[pcEdge.target.Name] = policyClass
    }

    return results
}

func (pa *PReviewAuditor) computePaths(userPaths []*edgePath, targetPath *edgePath, target string, pcEdge *edge) set.Set {
    computedPaths := set.NewSet()

    for _, userPath := range userPaths {
        lastUserEdge := userPath.edges[len(userPath.edges)-1]

        // if the last edge does not have any ops, it is not an association, so ignore it
        if lastUserEdge.ops == nil {
            continue
        }

        for i := 0; i < len(targetPath.edges); i++ {
            curEdge := targetPath.edges[i]
            // if the target of the last edge in a user resolvedPath does not match the target of the current edge in the target
            // resolvedPath, continue to the next target edge
            lastUserEdgeTarget := lastUserEdge.target.Name
            curEdgeSource := curEdge.source.Name
            curEdgeTarget := curEdge.target.Name

            // if the target of the last edge in a user path does not match the target of the current edge in the target path
            // AND if the target of the last edge in a user path does not match the source of the current edge in the target path
            //     OR if the target of the last edge in a user path does not match the target of the explain
            // continue to the next target edge
            if lastUserEdgeTarget != curEdgeTarget && (lastUserEdgeTarget != curEdgeSource || lastUserEdgeTarget == target) {
                continue
            }

            pathToTarget := make([]*edge, 0)
            for j := 0; j <= i; j++ {
                pathToTarget = append(pathToTarget, targetPath.edges[j])
            }

            resolvedPath := pa.resolvePath(userPath, pathToTarget, pcEdge)
            if resolvedPath == nil {
                continue
            }

            nodePath := resolvedPath.toNodePath(target)
            // add resolvedPath to policy class' paths
            computedPaths.Add(nodePath)
        }
    }

    return computedPaths
}

func (pa *PReviewAuditor) resolvePath(userPath *edgePath, pathToTarget []*edge, pcEdge *edge) *resolvedPath {
    if pcEdge.target.Type != graph.PC {
        return nil
    }

    // get the operations in this path
    // the operations are the ops of the association in the user path
    // convert * to actual operations
    ops := operations.NewOperationSet()
    for _, edge := range userPath.edges {
        if edge.ops != nil {
            ops = edge.ops
            // resolve the operation set
            pa.resolveOperationSet(ops, pa.resourceOps)

            break
        }
    }

    path := newEdgePath()
    // Collections.reverse(pathToTarget);
    for i, j := 0, len(pathToTarget)-1; i < j; i, j = i+1, j-1 {
        pathToTarget[i], pathToTarget[j] = pathToTarget[j], pathToTarget[i]
    }
    for _, edge := range userPath.edges {
        path.edges = append(path.edges, edge)
    }
    for _, edge := range pathToTarget {
        path.edges = append(path.edges, edge)
    }

    return newResolvedPath(pcEdge.target, path, ops)
}

/**
 * Removes any ops in ops that are not in resourceOps and converts special ops to actual ops (*, *a, *r)
 * @param ops the set of ops to check against the resource ops
 * @param resourceOps the set of resource operations
 */
func (pa *PReviewAuditor) resolveOperationSet(ops, resourceOps operations.OperationSet) {
    // if the permission set includes *, remove the * and add all resource operations
    if ops.Contains(operations.ALL_OPS) {
        ops.Remove(operations.ALL_OPS)
        ops.AddFrom(operations.AdminOps())
        ops.AddFrom(resourceOps)
    } else {
        // if the permissions includes *a or *r add all the admin ops/resource ops as necessary
        if ops.Contains(operations.ALL_ADMIN_OPS) {
            ops.Remove(operations.ALL_ADMIN_OPS)
            ops.AddFrom(operations.AdminOps())
        }
        if ops.Contains(operations.ALL_RESOURCE_OPS) {
            ops.Remove(operations.ALL_RESOURCE_OPS)
            ops.AddFrom(resourceOps)
        }
    }

    // remove any unknown ops
    ops.Filter(func(op interface{}) bool {
        return !resourceOps.Contains(op) && !operations.AdminOps().Contains(op)
    })
}

func (pa *PReviewAuditor) dfs(start *graph.Node) ([]*edgePath, error) {
    searcher := graph.NewDFS(pa.graph)

    paths := make([]*edgePath, 0)
    propPaths := make(map[string][]*edgePath)

    visitor := func(node *graph.Node) error {
        nodePaths := make([]*edgePath, 0)

        for parent := range pa.graph.Parents(node.Name).Iter() {
            n, err := pa.graph.Node(parent.(string))
            if err != nil {
                return err
            }
            ee := newEdge(node, n, nil)
            parentPaths := propPaths[parent.(string)]
            if len(parentPaths) == 0 {
                path := newEdgePath()
                path.edges = append(path.edges, ee)
                nodePaths = append([]*edgePath{path}, nodePaths...)
            } else {
                for _, p := range parentPaths {
                    parentPath := newEdgePath()
                    for _, e := range p.edges {
                        parentPath.edges = append(parentPath.edges, newEdge(e.source, e.target, e.ops))
                    }

                    parentPath.edges = append([]*edge{ee}, parentPath.edges...)
                    nodePaths = append(nodePaths, parentPath)
                }
            }
        }

        assocs, err := pa.graph.SourceAssociations(node.Name)
        if err != nil {
            return err
        }
        for target, ops := range assocs {
            targetNode, err := pa.graph.Node(target)
            if err != nil {
                return err
            }
            path := newEdgePath()
            path.edges = append(path.edges, newEdge(node, targetNode, ops))
            nodePaths = append(nodePaths, path)
        }

        // if the node being visited is the start node, add all the found nodePaths
        // TODO there might be a more efficient way of doing this
        // we don't need the if for users, only when the target is an OA, so it might have something to do with
        // leafs vs non leafs
        if node.Name == start.Name {
            paths = nodePaths
        } else {
            propPaths[node.Name] = nodePaths
        }

        return nil
    }

    propagator := func(parentNode, childNode *graph.Node) error {
        if _, ok := propPaths[childNode.Name]; !ok {
            propPaths[childNode.Name] = make([]*edgePath, 0)
        }
        childPaths := propPaths[childNode.Name]
        parentPaths := propPaths[parentNode.Name]

        for _, p := range parentPaths {
            path := newEdgePath()
            for _, e := range p.edges {
                path.edges = append(path.edges, newEdge(e.source, e.target, e.ops))
            }

            newPath := newEdgePath()
            newPath.edges = append(newPath.edges, path.edges...)
            ee := newEdge(childNode, parentNode, nil)
            newPath.edges = append([]*edge{ee}, newPath.edges...)
            childPaths = append(childPaths, newPath)
            propPaths[childNode.Name] = childPaths
        }

        if childNode.Name == start.Name {
            paths = propPaths[childNode.Name]
        }

        return nil
    }

    err := searcher.Traverse(start, propagator, visitor, graph.PARENTS)
    if err != nil {
        return nil, err
    }
    return paths, nil
}

type resolvedPath struct {
    pc   *graph.Node
    path *edgePath
    ops  set.Set
}

func newResolvedPath(pc *graph.Node, path *edgePath, ops set.Set) *resolvedPath {
    return &resolvedPath{pc, path, ops}
}

func (rp *resolvedPath) toNodePath(target string) *Path {
    nodePath := NewEmptyPath()
    nodePath.Operations = rp.ops

    if len(rp.path.edges) == 0 {
        return nodePath
    }

    var foundAssoc bool
    for _, edge := range rp.path.edges {
        var node *graph.Node
        if !foundAssoc {
            node = edge.target
        } else {
            node = edge.source
        }

        if len(nodePath.Nodes) == 0 {
            nodePath.Nodes = append(nodePath.Nodes, edge.source)
        }

        var contains bool
        for _, v := range nodePath.Nodes {
            if v.Name == node.Name {
                contains = true
                break
            }
        }

        if !contains {
            nodePath.Nodes = append(nodePath.Nodes, node)
        }

        if edge.ops != nil {
            foundAssoc = true
            if edge.target.Name == target {
                return nodePath
            }
        }
    }

    return nodePath
}

type edgePath struct {
    edges []*edge
}

func newEdgePath() *edgePath {
    return &edgePath{make([]*edge, 0)}
}

func newEdgePathWithEdges(edges []*edge) *edgePath {
    return &edgePath{edges}
}

func (e *edgePath) String() string {
    var s string
    for i, ee := range e.edges {
        s += ee.String()
        if i < len(e.edges)-1 {
            s += ", "
        }
    }

    return s
}

type edge struct {
    source, target *graph.Node
    ops            operations.OperationSet
}

func newEdge(source, target *graph.Node, ops operations.OperationSet) *edge {
    return &edge{source, target, ops}
}

func (e *edge) String() string {
    s := fmt.Sprintf("%s (%s)-->", e.source.Name, e.source.Type.String())
    if e.ops != nil {
        var i int
        for o := range e.ops.Iter() {
            s += o.(string)
            if i < e.ops.Len()-1 {
                s += ", "
            }
        }
        s += "-->"
    } else {
        s += ""
    }

    return s + e.target.Name + "(" + e.target.Type.String() + ")"
}
