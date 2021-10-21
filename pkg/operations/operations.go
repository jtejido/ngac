package operations

import (
	"ngac/internal/set"
)

const (
	WRITE                      = "write"
	READ                       = "read"
	CREATE_POLICY_CLASS        = "create policy class"
	ASSIGN_OBJECT_ATTRIBUTE    = "assign object attribute"
	ASSIGN_OBJECT_ATTRIBUTE_TO = "assign object attribute to"
	ASSIGN_OBJECT              = "assign object"
	ASSIGN_OBJECT_TO           = "assign object to"
	CREATE_NODE                = "create node"
	DELETE_NODE                = "delete node"
	UPDATE_NODE                = "update node"
	OBJECT_ACCESS              = "object access"
	ASSIGN_TO                  = "assign to"
	ASSIGN                     = "assign"
	ASSOCIATE                  = "associate"
	DISASSOCIATE               = "disassociate"
	CREATE_OBJECT              = "create object"
	CREATE_OBJECT_ATTRIBUTE    = "create object attribute"
	CREATE_USER_ATTRIBUTE      = "create user attribute"
	CREATE_USER                = "create user"
	DELETE_OBJECT              = "delete object"
	DELETE_OBJECT_ATTRIBUTE    = "delete object attribute"
	DELETE_USER_ATTRIBUTE      = "delete user attribute"
	DELETE_USER                = "delete user"
	DELETE_POLICY_CLASS        = "delete policy class"
	DEASSIGN                   = "deassign"
	DEASSIGN_FROM              = "deassign from"
	CREATE_ASSOCIATION         = "create association"
	UPDATE_ASSOCIATION         = "update association"
	DELETE_ASSOCIATION         = "delete association"
	GET_ASSOCIATIONS           = "get associations"
	RESET                      = "reset"
	GET_PERMISSIONS            = "get permissions"
	CREATE_PROHIBITION         = "create prohibition"
	UPDATE_PROHIBITION         = "update prohibition"
	VIEW_PROHIBITION           = "view prohibition"
	DELETE_PROHIBITION         = "delete prohibition"
	GET_ACCESSIBLE_CHILDREN    = "get accessible children"
	GET_PROHIBITED_OPS         = "get prohibited ops"
	GET_ACCESSIBLE_NODES       = "get accessible nodes"
	TO_JSON                    = "to json"
	FROM_JSON                  = "from json"
	ADD_OBLIGATION             = "add obligation"
	GET_OBLIGATION             = "get obligation"
	UPDATE_OBLIGATION          = "update obligation"
	DELETE_OBLIGATION          = "delete obligation"
	ENABLE_OBLIGATION          = "enable obligation"

	ALL_OPS          = "*"
	ALL_ADMIN_OPS    = "*a"
	ALL_RESOURCE_OPS = "*r"
)

var (
	admin_ops = set.NewSet(
		CREATE_POLICY_CLASS,
		ASSIGN_OBJECT_ATTRIBUTE,
		ASSIGN_OBJECT_ATTRIBUTE_TO,
		ASSIGN_OBJECT,
		ASSIGN_OBJECT_TO,
		CREATE_NODE,
		DELETE_NODE,
		UPDATE_NODE,
		ASSIGN_TO,
		ASSIGN,
		ASSOCIATE,
		DISASSOCIATE,
		CREATE_OBJECT,
		CREATE_OBJECT_ATTRIBUTE,
		CREATE_USER_ATTRIBUTE,
		CREATE_USER,
		DELETE_OBJECT,
		DELETE_OBJECT_ATTRIBUTE,
		DELETE_USER_ATTRIBUTE,
		DELETE_POLICY_CLASS,
		DELETE_USER,
		DEASSIGN,
		DEASSIGN_FROM,
		CREATE_ASSOCIATION,
		UPDATE_ASSOCIATION,
		DELETE_ASSOCIATION,
		GET_ASSOCIATIONS,
		GET_PERMISSIONS,
		CREATE_PROHIBITION,
		GET_ACCESSIBLE_CHILDREN,
		GET_PROHIBITED_OPS,
		GET_ACCESSIBLE_NODES,
		RESET,
		TO_JSON,
		FROM_JSON,
		UPDATE_PROHIBITION,
		DELETE_PROHIBITION,
		VIEW_PROHIBITION,
		ADD_OBLIGATION,
		GET_OBLIGATION,
		UPDATE_OBLIGATION,
		DELETE_OBLIGATION,
		ENABLE_OBLIGATION,
	)
)

func AdminOps() OperationSet {
	return admin_ops
}

type OperationSet set.Set

func NewOperationSet(ops ...interface{}) OperationSet {
	return set.NewSet(ops...)
}

func NewOperationSetFromSet(ops set.Set) OperationSet {
	return set.NewSetFromSlice(ops.ToSlice())
}
