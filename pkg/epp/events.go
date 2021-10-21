package epp

import (
	"ngac/internal/set"
	"ngac/pkg/context"
	"ngac/pkg/pip/graph"
	"ngac/pkg/pip/obligations"
)

const (
	ASSIGN_TO_EVENT     = "assign to"
	ASSIGN_EVENT        = "assign"
	DEASSIGN_FROM_EVENT = "deassign from"
	DEASSIGN_EVENT      = "deassign"
	CREATE_NODE_EVENT   = "create node"
	DELETE_NODE_EVENT   = "delete node"
	ACCESS_DENIED_EVENT = "access denied"
)

type EventContext interface {
	Event() string
	Target() *graph.Node
	UserCtx() context.Context
	MatchesPattern(*obligations.EventPattern, graph.Graph) bool
}

type eventContext struct {
	userCtx context.Context
	event   string
	target  *graph.Node
}

func NewEventContext(userCtx context.Context, event string, target *graph.Node) *eventContext {
	return &eventContext{userCtx, event, target}
}

func (ctx *eventContext) Event() string {
	return ctx.event
}

func (ctx *eventContext) Target() *graph.Node {
	return ctx.target
}

func (ctx *eventContext) UserCtx() context.Context {
	return ctx.userCtx
}

func (ctx *eventContext) MatchesPattern(pattern *obligations.EventPattern, g graph.Graph) bool {
	if pattern.Operations != nil {
		var found bool
		for _, op := range pattern.Operations {
			if op == ctx.event {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	matchSubject := pattern.Subject
	matchPolicyClass := pattern.PolicyClass
	matchTarget := pattern.Target

	return ctx.subjectMatches(g, ctx.userCtx.User(), ctx.userCtx.Process(), matchSubject) &&
		ctx.pcMatches(ctx.userCtx.User(), ctx.userCtx.Process(), matchPolicyClass) &&
		ctx.targetMatches(g, ctx.target, matchTarget)
}

func (ctx *eventContext) subjectMatches(g graph.Graph, user, process string, matchSubject *obligations.Subject) bool {
	if matchSubject == nil {
		return true
	}

	// any user
	if (matchSubject.AnyUser == nil && matchSubject.Process == nil) || (matchSubject.AnyUser != nil && len(matchSubject.AnyUser) == 0) {
		return true
	}

	// get the current user node
	userNode, err := g.Node(user)
	if err != nil {
		panic(err.Error())
	}

	if ctx.checkAnyUser(g, userNode, matchSubject.AnyUser) {
		return true
	}

	if matchSubject.User == userNode.Name {
		return true
	}

	return matchSubject.Process != nil && matchSubject.Process.Value == process
}

func (ctx *eventContext) checkAnyUser(g graph.Graph, userNode *graph.Node, anyUser []string) bool {
	if anyUser == nil || len(anyUser) == 0 {
		return true
	}

	dfs := graph.NewDFS(g)

	// check each user in the anyUser list
	// there can be users and user attributes
	for _, u := range anyUser {
		anyUserNode, err := g.Node(u)
		if err != nil {
			panic(err.Error())
		}

		// if the node in anyUser == the user than return true
		if anyUserNode.Name == userNode.Name {
			return true
		}

		// if the anyUser is not an UA, move to the next one
		if anyUserNode.Type != graph.UA {
			continue
		}

		nodes := set.NewSet()
		visitor := func(node *graph.Node) error {
			if node.Name == userNode.Name {
				nodes.Add(node.Name)
			}

			return nil
		}
		propagator := func(parent, child *graph.Node) error { return nil }

		err = dfs.Traverse(userNode, propagator, visitor, graph.PARENTS)
		if err != nil {
			panic(err)
		}

		if nodes.Contains(anyUserNode.Name) {
			return true
		}
	}

	return false
}

func (ctx *eventContext) pcMatches(user, process string, matchPolicyClass *obligations.PolicyClass) bool {
	// not yet implemented
	return true
}

func (ctx *eventContext) targetMatches(g graph.Graph, target *graph.Node, matchTarget *obligations.Target) bool {
	if matchTarget == nil {
		return true
	}

	if matchTarget.PolicyElements == nil && matchTarget.Containers == nil {
		return true
	}

	if matchTarget.Containers != nil {
		if len(matchTarget.Containers) == 0 {
			return true
		}

		// check that target is contained in any container
		containers := ctx.containersOf(g, target.Name)
		for _, evrContainer := range matchTarget.Containers {
			it := containers.Iterator()
			for it.HasNext() {
				contNode := it.Next().(*graph.Node)
				if ctx.nodesMatch(evrContainer, contNode) {
					return true
				}
			}
		}

		return false
	} else if matchTarget.PolicyElements != nil {
		if len(matchTarget.PolicyElements) == 0 {
			return true
		}

		// check that target is in the list of policy elements
		for _, evrNode := range matchTarget.PolicyElements {
			if ctx.nodesMatch(evrNode, target) {
				return true
			}
		}

		return false
	}

	return false
}

func (ctx *eventContext) containersOf(g graph.Graph, name string) set.Set {
	nodes := set.NewSet()
	parents := g.Parents(name)
	for parent := range parents.Iter() {
		nn, err := g.Node(parent.(string))
		if err != nil {
			panic(err)
		}
		nodes.Add(nn)
		nodes.AddFrom(ctx.containersOf(g, parent.(string)))
	}

	return nodes
}

func (ctx *eventContext) nodesMatch(evrNode *obligations.EvrNode, node *graph.Node) bool {
	if evrNode.Name != node.Name {
		return false
	}

	if evrNode.Type != node.Type.String() {
		return false
	}

	for k, v := range evrNode.Properties {
		if val, ok := node.Properties[k]; ok {
			if val != v {
				return false
			}
		} else {
			return false
		}
	}

	return true
}

type AssignEvent struct {
	eventContext
	ParentNode *graph.Node
}

func NewAssignEvent(userCtx context.Context, target, parentNode *graph.Node) *AssignEvent {
	ans := new(AssignEvent)
	ans.userCtx = userCtx
	ans.event = ASSIGN_EVENT
	ans.target = target
	ans.ParentNode = parentNode
	return ans
}

type AssignToEvent struct {
	eventContext
	ChildNode *graph.Node
}

func NewAssignToEvent(userCtx context.Context, target, childNode *graph.Node) *AssignToEvent {
	ans := new(AssignToEvent)
	ans.userCtx = userCtx
	ans.event = ASSIGN_TO_EVENT
	ans.target = target
	ans.ChildNode = childNode
	return ans
}

type AssociationEvent struct {
	eventContext
	Source, assocTarget *graph.Node
}

func NewAssociationEvent(userCtx context.Context, source, target *graph.Node) *AssociationEvent {
	ans := new(AssociationEvent)
	ans.userCtx = userCtx
	ans.event = "association"
	ans.target = target
	ans.assocTarget = target
	ans.Source = source
	return ans
}

func (ae *AssociationEvent) Target() *graph.Node {
	return ae.assocTarget
}

type CreateNodeEvent struct {
	eventContext
	parents set.Set
}

func NewCreateNodeEvent(userCtx context.Context, deletedNode *graph.Node, initialParent string, additionalParents ...string) *CreateNodeEvent {
	ans := new(CreateNodeEvent)
	ans.userCtx = userCtx
	ans.event = CREATE_NODE_EVENT
	ans.target = deletedNode
	ans.parents = set.NewSet()
	for _, a := range additionalParents {
		ans.parents.Add(a)
	}
	ans.parents.Add(initialParent)
	return ans
}

type DeassignEvent struct {
	eventContext
	ParentNode *graph.Node
}

func NewDeassignEvent(userCtx context.Context, target, parentNode *graph.Node) *DeassignEvent {
	ans := new(DeassignEvent)
	ans.userCtx = userCtx
	ans.event = DEASSIGN_EVENT
	ans.target = target
	ans.ParentNode = parentNode
	return ans
}

type DeassignFromEvent struct {
	eventContext
	ChildNode *graph.Node
}

func NewDeassignFromEvent(userCtx context.Context, target, childNode *graph.Node) *DeassignFromEvent {
	ans := new(DeassignFromEvent)
	ans.userCtx = userCtx
	ans.event = DEASSIGN_FROM_EVENT
	ans.target = target
	ans.ChildNode = childNode
	return ans
}

type DeleteAssociationEvent struct {
	eventContext
	Source, assocTarget *graph.Node
}

func NewDeleteAssociationEvent(userCtx context.Context, source, target *graph.Node) *DeleteAssociationEvent {
	ans := new(DeleteAssociationEvent)
	ans.userCtx = userCtx
	ans.event = "delete association"
	ans.target = target
	ans.assocTarget = target
	ans.Source = source
	return ans
}

func (ae *DeleteAssociationEvent) Target() *graph.Node {
	return ae.assocTarget
}

type DeleteNodeEvent struct {
	eventContext
	parents set.Set
}

func NewDeleteNodeEvent(userCtx context.Context, deletedNode *graph.Node, parents set.Set) *DeleteNodeEvent {
	ans := new(DeleteNodeEvent)
	ans.userCtx = userCtx
	ans.event = "delete association"
	ans.target = deletedNode
	ans.parents = parents
	return ans
}

type ObjectAccessEvent struct {
	eventContext
}

func NewObjectAccessEvent(userCtx context.Context, event string, target *graph.Node) *ObjectAccessEvent {
	ans := new(ObjectAccessEvent)
	ans.userCtx = userCtx
	ans.event = event
	ans.target = target
	return ans
}
