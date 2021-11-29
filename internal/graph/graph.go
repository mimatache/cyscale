package graph

import (
	"errors"
	"fmt"
	"sync"
)

var (
	ErrNotFound = errors.New("could not find an element matching the request")
)

// FilterNodes is used as an interface for filtering functionality. This allows each user to provide their own way of filtering different items
type FilterNodes func(node *Node) bool

func FilterNodesByLabel(label string) FilterNodes {
	return func(node *Node) bool {
		return node.label == label
	}
}

type FilterRelationship func(rel Relationship) bool

func FilterRelByLabel(label string) FilterRelationship {
	return func(rel Relationship) bool {
		return rel.Label == label
	}
}

func FilterRelByTo(toID string) FilterRelationship {
	return func(rel Relationship) bool {
		return rel.To == toID
	}
}

func FilterRelByFrom(fromID string) FilterRelationship {
	return func(rel Relationship) bool {
		return rel.From == fromID
	}
}

func New() *Graph {
	return &Graph{
		nodes:         map[string]*Node{},
		relationships: map[string]Relationship{},
	}
}

// Graph represents a collection of different nodes of the same type
type Graph struct {
	sync.RWMutex
	nodes         map[string]*Node
	relationships map[string]Relationship
}

// InsertNode adds a new node to the graph
func (g *Graph) InsertNode(label string, body []byte) *Node {
	g.Lock()
	defer g.Unlock()
	node := newNode(label, body)
	g.nodes[node.id] = node
	return node.Copy()
}

// GetNodeByID returns the node that has the given ID
func (g *Graph) GetNodeByID(id string) (*Node, error) {
	g.RLock()
	defer g.RUnlock()
	item, ok := g.nodes[id]
	if !ok {
		return nil, fmt.Errorf("%w; node with id '%s'", ErrNotFound, id)
	}
	return item, nil
}

// ListNodes returns a map of all the nodes that match all the where clauses provided.
func (g *Graph) ListNodes(where ...FilterNodes) []*Node {
	g.RLock()
	defer g.RUnlock()
	matchingNodes := make([]*Node, 0, len(g.nodes))
	for _, item := range g.nodes {
		matches := true
		for _, clause := range where {
			if ok := clause(item); !ok {
				matches = false
				break
			}
		}
		if matches {
			matchingNodes = append(matchingNodes, item)
		}
	}

	return matchingNodes
}

// Link is used to establish a relationship between the current item and another one
func (g *Graph) AddRelationship(fromID, toID, label string) (Relationship, error) {
	fromNode, err := g.GetNodeByID(fromID)
	if err != nil {
		return Relationship{}, fmt.Errorf("getNodeByID %s; %w", fromID, err)
	}

	toNode, err := g.GetNodeByID(toID)
	if err != nil {
		return Relationship{}, fmt.Errorf("getNodeByID %s; %w", fromID, err)
	}

	rel := newRelationship(fromNode, toNode, label)
	fromNode.addRelationship(rel)
	g.relationships[rel.ID] = rel

	return rel, nil
}

func (g *Graph) GetRelationshipByID(id string) (Relationship, error) {
	g.RLock()
	defer g.RUnlock()
	item, ok := g.relationships[id]
	if !ok {
		return Relationship{}, fmt.Errorf("%w; relationship with id '%s'", ErrNotFound, id)
	}
	return item, nil
}

func (g *Graph) ListRelationships(filters ...FilterRelationship) []Relationship {
	g.RLock()
	defer g.RUnlock()
	matchingRelationships := make([]Relationship, 0, len(g.relationships))
	for _, item := range g.relationships {
		matches := true
		for _, clause := range filters {
			if ok := clause(item); !ok {
				matches = false
				break
			}
		}
		if matches {
			matchingRelationships = append(matchingRelationships, item)
		}
	}

	return matchingRelationships
}
