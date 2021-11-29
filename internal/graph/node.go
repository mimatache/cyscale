package graph

import (
	"sync"

	guuid "github.com/google/uuid"
)

func newNode(label string, body []byte) *Node {
	return &Node{
		id:            guuid.New().String(),
		label:         label,
		Body:          body,
		relationships: []Relationship{},
	}
}

// Node represents an item in the graph. It contains the ID of the element, the body and it's relationships to other items
type Node struct {
	sync.RWMutex
	id            string
	label         string
	Body          []byte
	relationships []Relationship
}

func (n *Node) GetID() string {
	return n.id
}

func (n *Node) GetLabel() string {
	return n.label
}

// Copy returns a duplicate of the Node. This is done to avoid unwanted changes and race conditions
func (n *Node) Copy() *Node {
	relCopy := make([]Relationship, len(n.relationships))
	copy(relCopy, n.relationships)
	return &Node{
		id:            n.id,
		label:         n.label,
		Body:          n.Body,
		relationships: relCopy,
	}
}

func (n *Node) addRelationship(relationship Relationship) {
	n.Lock()
	defer n.Unlock()
	n.relationships = append(n.relationships, relationship)
}

func (n *Node) ListRelationships(filters ...FilterRelationship) []Relationship {
	n.RLock()
	defer n.RUnlock()
	if len(filters) == 0 {
		relCopy := make([]Relationship, len(n.relationships))
		copy(relCopy, n.relationships)
		return relCopy
	}
	relCopy := make([]Relationship, 0, len(n.relationships))
	for _, rel := range n.relationships {
		for _, filter := range filters {
			if filter(rel) {
				relCopy = append(relCopy, rel)
			}
		}
	}
	return relCopy
}
