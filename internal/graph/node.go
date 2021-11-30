package graph

import (
	"fmt"
	"sync"

	guuid "github.com/google/uuid"
)

func newNode(name, label string, body []byte) *Node {
	return &Node{
		id:            guuid.New().String(),
		label:         label,
		name:          name,
		Body:          body,
		relationships: []string{},
	}
}

// Node represents an item in the graph. It contains the ID of the element, the body and it's relationships to other items
type Node struct {
	sync.RWMutex
	id            string
	name          string
	label         string
	Body          []byte
	relationships []string
}

func (n *Node) GetID() string {
	return n.id
}

func (n *Node) GetName() string {
	return n.name
}

func (n *Node) GetLabel() string {
	return n.label
}

// Copy returns a duplicate of the Node. This is done to avoid unwanted changes and race conditions
func (n *Node) Copy() *Node {
	relCopy := make([]string, len(n.relationships))
	copy(relCopy, n.relationships)
	return &Node{
		id:            n.id,
		name:          n.name,
		label:         n.label,
		Body:          n.Body,
		relationships: relCopy,
	}
}

func (n *Node) addRelationship(relationship string) {
	n.Lock()
	defer n.Unlock()
	n.relationships = append(n.relationships, relationship)
}

func (n *Node) ListRelationships() []string {
	n.RLock()
	defer n.RUnlock()
	return n.relationships
}

func (n *Node) String() string {
	return fmt.Sprintf("{Asset:%s}", n.name)
}
