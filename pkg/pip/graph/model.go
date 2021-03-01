package graph

import (
	"fmt"
	"github.com/jtejido/ngac/pkg/operations"
	"strings"
	"sync"
)

var (
	validAssignments = map[NodeType]map[NodeType]bool{
		PC: {},
		OA: {OA: true, PC: true},
		UA: {UA: true, PC: true},
		O:  {OA: true},
		U:  {UA: true},
	}

	validAssociations = map[NodeType]map[NodeType]bool{
		PC: {},
		OA: {},
		UA: {OA: true, UA: true},
		O:  {},
		U:  {},
	}
)

const (
	HASH_LENGTH                 = 163
	SUPER_KEYWORD               = "super"
	WILDCARD                    = "*"
	PASSWORD_PROPERTY           = "password"
	DESCRIPTION_PROPERTY        = "description"
	NAMESPACE_PROPERTY          = "namespace"
	DEFAULT_NAMESPACE           = "default"
	SOURCE_PROPERTY             = "source"
	STORAGE_PROPERTY            = "storage"
	GCS_STORAGE                 = "google"
	AWS_STORAGE                 = "amazon"
	LOCAL_STORAGE               = "local"
	CONTENT_TYPE_PROPERTY       = "content_type"
	SIZE_PROPERTY               = "size"
	PATH_PROPERTY               = "path"
	BUCKET_PROPERTY             = "bucket"
	COLUMN_INDEX_PROPERTY       = "column_index"
	ORDER_BY_PROPERTY           = "order_by"
	ROW_INDEX_PROPERTY          = "row_index"
	SESSION_USER_ID_PROPERTY    = "user_id"
	SCHEMA_COMP_PROPERTY        = "schema_comp"
	SCHEMA_COMP_SCHEMA_PROPERTY = "schema"
	SCHEMA_COMP_TABLE_PROPERTY  = "table"
	SCHEMA_COMP_ROW_PROPERTY    = "row"
	SCHEMA_COMP_COLUMN_PROPERTY = "col"
	SCHEMA_COMP_CELL_PROPERTY   = "cell"
	SCHEMA_NAME_PROPERTY        = "schema"
	COLUMN_CONTAINER_NAME       = "Columns"
	ROW_CONTAINER_NAME          = "Rows"
	COLUMN_PROPERTY             = "column"
	REP_PROPERTY                = "rep"
)

type PropertyPair [2]string

type NodeType int

/**
 * Allowed types of nodes in an NGAC Graph
 * <p>
 * OA = Object Attribute
 * UA = user attribute
 * U = User
 * O = Object
 * PC = policy class
 * OS = Operation Set
 */
const (
	NOOP NodeType = iota - 1
	OA
	UA
	U
	O
	PC
)

func (nt NodeType) String() string {
	switch nt {
	case OA:
		return "OA"
	case UA:
		return "UA"
	case U:
		return "U"
	case O:
		return "O"
	case PC:
		return "PC"
	default:
		return "nil"
	}
}

func ToNodeType(t string) NodeType {
	switch strings.ToUpper(t) {
	case "OA":
		return OA
	case "UA":
		return UA
	case "U":
		return U
	case "O":
		return O
	case "PC":
		return PC
	default:
		return NOOP
	}
}

type Node struct {
	Name       string
	Type       NodeType
	Properties PropertyMap
}

func NewNode() *Node {
	return &Node{Properties: NewPropertyMap()}
}

func NewNodeFromNode(node *Node) *Node {
	ans := &Node{Name: node.Name, Type: node.Type, Properties: node.Properties}
	if node.Properties == nil {
		ans.Properties = NewPropertyMap()
	}
	return ans
}

func NewNodeWithFields(Name string, Type NodeType, Prop PropertyMap) *Node {
	return &Node{Name, Type, Prop}
}

func NewNodeWithoutProps(Name string, Type NodeType) *Node {
	return &Node{Name, Type, nil}
}

func (n *Node) Equals(i interface{}) bool {
	if v, ok := i.(*Node); ok {
		return n.Name == v.Name
	}

	return false
}

func (n *Node) String() string {
	return fmt.Sprintf("%s:%s", n.Name, n.Type.String())
}

type PropertyMap map[string]string

func NewPropertyMap() PropertyMap {
	return make(PropertyMap)
}

func ToProperties(pairs ...PropertyPair) PropertyMap {
	props := NewPropertyMap()
	for _, p := range pairs {
		props[p[0]] = p[1]
	}

	return props
}

type Edge interface {
	From() string
	To() string
}

type Relationship struct {
	Source, Target string
}

func (r *Relationship) From() string {
	return r.Source
}

func (r *Relationship) To() string {
	return r.Target
}

func (r *Relationship) Equals(i interface{}) bool {
	if v, ok := i.(*Relationship); ok {
		return r.Source == v.Source && r.Target == v.Target
	}

	return false
}

type Association struct {
	Relationship
	Operations operations.OperationSet
}

func (a *Association) Equals(i interface{}) bool {
	if v, ok := i.(*Association); ok {
		return a.Source == v.Source && a.Target == v.Target && a.Operations.Equal(v.Operations)
	}

	return false
}

func CheckAssociation(uaType, targetType NodeType) error {
	if !validAssociations[uaType][targetType] {
		return fmt.Errorf("invalid association: %q to %q", uaType.String(), targetType.String())
	}

	return nil

}

type Assignment struct {
	Relationship
}

func (a *Assignment) Equals(i interface{}) bool {
	if v, ok := i.(*Assignment); ok {
		return a.Source == v.Source && a.Target == v.Target
	}

	return false
}

func CheckAssignment(childType, parentType NodeType) error {
	if !validAssignments[childType][parentType] {
		return fmt.Errorf("invalid assignment: %q to %q", childType.String(), parentType.String())
	}

	return nil
}
