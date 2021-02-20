package graph

import (
	"fmt"
	"github.com/jtejido/ngac/internal/omap"
	"github.com/jtejido/ngac/operations"
	"strings"
)

var (
	validAssignments  = omap.NewOrderedMap()
	validAssociations = omap.NewOrderedMap()
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

func init() {
	validAssociations.Add(PC, []NodeType{})
	validAssociations.Add(OA, []NodeType{})
	validAssociations.Add(O, []NodeType{})
	validAssociations.Add(UA, []NodeType{UA, OA})
	validAssociations.Add(U, []NodeType{})

	validAssignments.Add(PC, []NodeType{})
	validAssignments.Add(OA, []NodeType{PC, OA})
	validAssignments.Add(O, []NodeType{OA})
	validAssignments.Add(UA, []NodeType{UA, PC})
	validAssignments.Add(U, []NodeType{UA})
}

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
	ALL NodeType = iota - 1
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
		return ALL
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

func (n *Node) Equals(i interface{}) bool {
	if v, ok := i.(*Node); ok {
		return n.Name == v.Name
	}

	return false
}

func (n *Node) String() string {
	return fmt.Sprintf("%s:%s", n.Name, n.Type.String())
}

type PropertyMap omap.OrderedMap

func NewPropertyMap() PropertyMap {
	return omap.NewOrderedMap()
}

func ToProperties(pairs ...PropertyPair) PropertyMap {
	props := omap.NewOrderedMap()
	for _, p := range pairs {
		props.Add(p[0], p[1])
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
	c, _ := validAssociations.Get(uaType)
	check := c.([]NodeType)
	for _, nt := range check {
		if nt == targetType {
			return nil
		}
	}

	return fmt.Errorf("cannot associate a node of type %s to a node of type %s", uaType.String(), targetType.String())
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
	c, _ := validAssignments.Get(childType)
	check := c.([]NodeType)
	for _, nt := range check {
		if nt == parentType {
			return nil
		}
	}

	return fmt.Errorf("cannot assign a node of type %s to a node of type %s", childType.String(), parentType.String())
}
