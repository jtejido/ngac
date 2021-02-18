package epp

import (
	"github.com/jtejido/ngac/pip/graph"
)

const (
	ASSIGN_TO_EVENT     = "assign to"
	ASSIGN_EVENT        = "assign"
	DEASSIGN_FROM_EVENT = "deassign from"
	DEASSIGN_EVENT      = "deassign"
	ACCESS_DENIED_EVENT = "deassign"
)

type EventContext interface {
	Event() string
	Target() *graph.Node
}

type eventContext struct {
	Event  string
	Target *graph.Node
}

func NewEventContext(event string, target *graph.Node) *eventContext {
	return &eventContext{event, target}
}

func (ctx *eventContext) Event() string {
	return ctx.Event
}

func (ctx *eventContext) Target() *graph.Node {
	return ctx.Target
}

type ObjectAccessEvent struct {
	eventContext
}

func NewObjectAccessEvent(event string, target *graph.Node) *ObjectAccessEvent {
	ans := new(ObjectAccessEvent)
	ans.Event = event
	ans.Target = target
	return ans
}

type AssignEvent struct {
	eventContext
	ParentNode *graph.Node
}

func NewAssignEvent(target, parentNode *graph.Node) *AssignEvent {
	ans := new(AssignEvent)
	ans.Event = ASSIGN_EVENT
	ans.Target = target
	ans.ParentNode = parentNode
	return ans
}

type AssignToEvent struct {
	eventContext
	ChildNode *graph.Node
}

func NewAssignToEvent(target, childNode *graph.Node) *AssignToEvent {
	ans := new(AssignToEvent)
	ans.Event = ASSIGN_TO_EVENT
	ans.Target = target
	ans.ChildNode = parentNode
	return ans
}

type DeassignEvent struct {
	eventContext
	ParentNode *graph.Node
}

func NewDeassignEvent(target, parentNode *graph.Node) *DeassignEvent {
	ans := new(DeassignEvent)
	ans.Event = DEASSIGN_EVENT
	ans.Target = target
	ans.ParentNode = parentNode
	return ans
}

type DeassignFromEvent struct {
	eventContext
	ChildNode *graph.Node
}

func NewDeassignFromEvent(target, childNode *graph.Node) *DeassignFromEvent {
	ans := new(DeassignFromEvent)
	ans.Event = DEASSIGN_FROM_EVENT
	ans.Target = target
	ans.ChildNode = parentNode
	return ans
}
