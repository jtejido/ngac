package decider

import (
	"github.com/jtejido/ngac/internal/set"
	"github.com/jtejido/ngac/operations"
	"github.com/jtejido/ngac/pip/graph"
	"github.com/jtejido/ngac/pip/prohibitions"
)

var (
	prd *PReviewDecider
	_   Decider = prd
)

// An implementation of the Decider interface that uses an in memory NGAC graph
type PReviewDecider struct {
	graph        graph.Graph
	prohibitions prohibitions.Prohibitions
	ResourceOps  operations.OperationSet
}

func NewPReviewDecider(graph graph.Graph, resourceOps operations.OperationSet) *PReviewDecider {
	return NewPReviewDeciderWithProhibitions(graph, prohibitions.NewMemProhibitions(), resourceOps)
}

func NewPReviewDeciderWithProhibitions(graph graph.Graph, prohibs prohibitions.Prohibitions, resourceOps operations.OperationSet) *PReviewDecider {
	if graph == nil {
		panic("graph cannot be nil")
	}
	if resourceOps == nil {
		panic("resourceOps cannot be nil")
	}
	if prohibs == nil {
		prohibs = prohibitions.NewMemProhibitions()
	}

	d := new(PReviewDecider)
	d.graph = graph
	d.prohibitions = prohibs
	d.ResourceOps = resourceOps
	return d
}

func (pr *PReviewDecider) Check(subject, process, target string, perms ...interface{}) bool {
	allowed := pr.List(subject, process, target)

	if len(perms) == 0 {
		return allowed.Len() > 0
	}
	return allowed.Contains(perms...)
}

func (pr *PReviewDecider) List(subject, process, target string) set.Set {
	perms := set.NewSet()

	// traverse the user side of the graph to get the associations
	userCtx, err := pr.processUserDAG(subject, process)
	if err != nil {
		return perms
	}

	if len(userCtx.borderTargets) == 0 {
		return perms
	}

	// traverse the target side of the graph to get permissions per policy class
	targetCtx, err := pr.processTargetDAG(target, userCtx)
	if err != nil {
		return perms
	}
	// resolve the permissions
	return pr.resolvePermissions(userCtx, targetCtx, target)
}

func (pr *PReviewDecider) Filter(subject, process string, nodes set.Set, perms ...interface{}) set.Set {
	n := set.NewSet()
	for nn := range nodes.Iter() {
		node := nn.(string)
		if pr.Check(subject, process, node, perms...) {
			n.Add(node)
		}
	}

	return n
}

func (pr *PReviewDecider) Children(subject, process, target string, perms ...interface{}) set.Set {
	children := pr.graph.Children(target)
	return pr.Filter(subject, process, children, perms...)
}

func (pr *PReviewDecider) CapabilityList(subject, process string) map[string]set.Set {
	results := make(map[string]set.Set)

	//get border nodes.  Can be OA or UA.  Return empty set if no OAs are reachable
	userCtx, err_u := pr.processUserDAG(subject, process)
	if err_u != nil {
		return results
	}
	if len(userCtx.borderTargets) == 0 {
		return results
	}

	for borderTarget, _ := range userCtx.borderTargets {
		n, err := pr.graph.Node(borderTarget)
		if err != nil {
			return results
		}
		objects := pr.ascendants(n.Name, make(map[string]set.Set))
		for object := range objects.Iter() {
			objn := object.(string)
			// run dfs on the object
			targetCtx, err_t := pr.processTargetDAG(objn, userCtx)
			if err_t != nil {
				return results
			}
			permissions := pr.resolvePermissions(userCtx, targetCtx, objn)
			results[objn] = permissions
		}
	}

	return results
}

func (pr *PReviewDecider) GenerateACL(target, process string) map[string]set.Set {
	acl := make(map[string]set.Set)

	search := pr.graph.Search(graph.U, nil)
	for user := range search.Iter() {
		usern := user.(*graph.Node)
		list := pr.List(usern.Name, process, target)
		acl[usern.Name] = list
	}

	return acl
}

func (pr *PReviewDecider) resolvePermissions(uctx *userContext, tctx *targetContext, target string) set.Set {
	allowed := pr.resolveAllowedPermissions(tctx)
	pr.resolveSpecialPermissions(allowed)
	allowed.Filter(func(op interface{}) bool {
		return !pr.ResourceOps.Contains(op) && !operations.AdminOps().Contains(op)
	})

	denied := pr.resolveProhibitions(uctx, tctx, target)
	allowed.RemoveFrom(denied)

	return allowed
}

func (pr *PReviewDecider) resolveAllowedPermissions(tctx *targetContext) set.Set {
	pcMap := tctx.pcSet

	allowed := set.NewSet()
	first := true
	for _, ops := range pcMap {
		if first {
			allowed.AddFrom(ops)
			first = false
		} else {
			if allowed.Contains(operations.ALL_OPS) {
				// clear all of the existing permissions because the intersection already had *
				// all permissions can be added
				allowed.Clear()
				allowed.AddFrom(ops)
			} else {
				// if the ops for the pc are empty then the user has no permissions on the target
				if ops.Len() == 0 {
					allowed.Clear()
					break
				} else if !ops.Contains(operations.ALL_OPS) {
					allowed.RetainFrom(ops)
				}
			}
		}
	}

	return allowed
}

func (pr *PReviewDecider) resolveSpecialPermissions(permissions set.Set) {
	// if the permission set includes *, remove the * and add all resource operations

	if permissions.Contains(operations.ALL_OPS) {
		permissions.Remove(operations.ALL_OPS)
		permissions.AddFrom(operations.AdminOps())
		permissions.AddFrom(pr.ResourceOps)
	} else {
		// if the permissions includes *a or *r add all the admin ops/resource ops as necessary
		if permissions.Contains(operations.ALL_ADMIN_OPS) {
			permissions.Remove(operations.ALL_ADMIN_OPS)
			permissions.AddFrom(operations.AdminOps())
		}
		if permissions.Contains(operations.ALL_RESOURCE_OPS) {
			permissions.Remove(operations.ALL_RESOURCE_OPS)
			permissions.AddFrom(pr.ResourceOps)
		}
	}
}

func (pr *PReviewDecider) resolveProhibitions(uctx *userContext, tctx *targetContext, target string) set.Set {
	denied := set.NewSet()
	prohibs := uctx.prohibitions
	reachedTargets := tctx.reachedTargets

	for p := range prohibs.Iter() {
		proh := p.(*prohibitions.Prohibition)

		inter := proh.Intersection
		containers := proh.Containers()

		var addOps bool
		for contName, isComplement := range containers {
			if target == contName {
				addOps = false
				if inter {
					// if the target is a container and the prohibition evaluates the intersection
					// the whole prohibition is not satisfied
					break
				} else {
					// continue checking the remaining conditions
					continue
				}
			}
			if !isComplement && reachedTargets.Contains(contName) || isComplement && !reachedTargets.Contains(contName) {
				addOps = true

				// if the prohibition is not intersection, one satisfied container condition means
				// the prohibition is satisfied
				if !inter {
					break
				}

				continue
			}

			// since the intersection requires the target to satisfy each node condition in the prohibition
			// if one is not satisfied then the whole is not satisfied
			addOps = false

			// if the prohibition is the intersection, one unsatisfied container condition means the whole
			// prohibition is not satisfied
			if inter {
				break
			}

		}

		if addOps {
			denied.AddFrom(proh.Operations)
		}
	}

	return denied
}

/**
 * Perform a depth first search on the object side of the graph.  Start at the target node and recursively visit nodes
 * until a policy class is reached.  On each node visited, collect any operation the user has on the target. At the
 * end of each dfs iteration the visitedNodes map will contain the operations the user is permitted on the target under
 * each policy class.
 *
 * @param target      the name of the current target node.
 */
func (pr *PReviewDecider) processTargetDAG(target string, userCtx *userContext) (*targetContext, error) {
	borderTargets := userCtx.borderTargets

	visitedNodes := make(map[string]map[string]set.Set)
	reachedTargets := set.NewSet()

	visitor := func(node *graph.Node) error {
		// add this node to reached prohibited targets if it has any prohibitions
		reachedTargets.Add(node.Name)
		var nodeCtx map[string]set.Set
		if v, ok := visitedNodes[node.Name]; ok {
			nodeCtx = v
		} else {
			nodeCtx = make(map[string]set.Set)
		}

		if len(nodeCtx) == 0 {
			visitedNodes[node.Name] = nodeCtx
		}

		if node.Type == graph.PC {
			nodeCtx[node.Name] = set.NewSet()
		} else {
			if uaOps, ok := borderTargets[node.Name]; ok {
				for pc, v := range nodeCtx {
					var pcOps set.Set
					if v != nil {
						pcOps = v
					} else {
						pcOps = set.NewSet()
					}

					pcOps.AddFrom(uaOps)
					nodeCtx[pc] = pcOps
				}
			}
		}

		return nil
	}

	propagator := func(parent, child *graph.Node) error {
		parentCtx := visitedNodes[parent.Name]
		var nodeCtx map[string]set.Set
		if v, ok := visitedNodes[child.Name]; ok {
			nodeCtx = v
		} else {
			nodeCtx = make(map[string]set.Set)
		}

		for name, v := range parentCtx {
			var ops set.Set
			if v2, ok := nodeCtx[name]; ok {
				ops = v2
			} else {
				ops = set.NewSet()
			}

			ops.AddFrom(v)
			nodeCtx[name] = ops
		}

		visitedNodes[child.Name] = nodeCtx

		return nil
	}

	ss := graph.NewDFS(pr.graph)
	n, err := pr.graph.Node(target)
	if err != nil {
		return nil, err
	}
	err = ss.Traverse(n, propagator, visitor, graph.PARENTS)
	if err != nil {
		return nil, err
	}

	return &targetContext{visitedNodes[target], reachedTargets}, nil
}

/**
 * Find the target nodes that are reachable by the subject via an association. This is done by a breadth first search
 * starting at the subject node and walking up the user side of the graph until all user attributes the subject is assigned
 * to have been visited.  For each user attribute visited, get the associations it is the source of and store the
 * target of that association as well as the operations in a map. If a target node is reached multiple times, add any
 * new operations to the already existing ones.
 *
 * @return a Map of target nodes that the subject can reach via associations and the operations the user has on each.
 */
func (pr *PReviewDecider) processUserDAG(subject, process string) (*userContext, error) {
	searcher := graph.NewBFS(pr.graph)
	start, err := pr.graph.Node(subject)
	if err != nil {
		return nil, err
	}

	borderTargets := make(map[string]set.Set)
	reachedProhibitions := set.NewSet(pr.prohibitions.ProhibitionsFor(process).ToSlice()...)

	// if the start node is an UA, get it's associations
	if start.Type == graph.UA {
		assocs, err := pr.graph.SourceAssociations(start.Name)
		if err != nil {
			return nil, err
		}

		pr.collectAssociations(assocs, borderTargets)
	}

	visitor := func(node *graph.Node) error {
		reachedProhibitions.AddFrom(pr.prohibitions.ProhibitionsFor(node.Name))

		//get the parents of the subject to start bfs on user side
		parents := pr.graph.Parents(node.Name)
		it := parents.Iterator()
		for it.HasNext() {
			parentNode := it.Next().(string)
			//get the associations the current parent node is the source of
			assocs, err := pr.graph.SourceAssociations(parentNode)

			if err != nil {
				return err
			}

			//collect the target and operation information for each association
			pr.collectAssociations(assocs, borderTargets)
			//add all of the current parent node's parents to the queue
			parents.AddFrom(pr.graph.Parents(parentNode))
		}

		return nil
	}

	// nothing is being propagated
	propagator := func(parent, child *graph.Node) error { return nil }

	// start the bfs
	if err := searcher.Traverse(start, propagator, visitor, graph.PARENTS); err != nil {
		return nil, err
	}

	return &userContext{borderTargets, reachedProhibitions}, nil
}

func (pr *PReviewDecider) collectAssociations(assocs map[string]operations.OperationSet, borderTargets map[string]set.Set) {
	for target, ops := range assocs {
		exOps := borderTargets[target]
		//if the target is not in the map already, put it
		//else add the found operations to the existing ones.
		if exOps == nil {
			borderTargets[target] = ops
		} else {
			ops.AddFrom(exOps)
			borderTargets[target] = ops
		}
	}
}

func (pr *PReviewDecider) ascendants(vNode string, cache map[string]set.Set) set.Set {
	ascendants := set.NewSet()
	ascendants.Add(vNode)

	children := pr.graph.Children(vNode)
	if children.Len() == 0 {
		return ascendants
	}

	ascendants.AddFrom(children)
	for s := range children.Iter() {
		child := s.(string)
		if _, ok := cache[child]; !ok {
			cache[child] = pr.ascendants(child, cache)
		}
		ascendants.AddFrom(cache[child])
	}

	return ascendants
}

type userContext struct {
	borderTargets map[string]set.Set
	prohibitions  set.Set
}

type targetContext struct {
	pcSet          map[string]set.Set
	reachedTargets set.Set
}
