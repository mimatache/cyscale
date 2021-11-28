package graph

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"
)

var (
	ErrNotFound = errors.New("could not find an element matching the request")
)

type Graph interface {
	AddNodeTypes(nodesType ...string)
	GetNodesByType(nodesType string) (Nodes, error)
	ListNodeTypes() []string
}

type Nodes interface {
}

type neighbor struct {
	ID   string
	Type string
}

type relationships map[string][]neighbor

type node struct {
	Body          json.RawMessage
	Relationships relationships
}

type nodes map[string]node

type assets struct {
	items map[string]nodes
	sync.RWMutex
}

// AddNodeTypes is used to add new types of entities to the available ones.
// This operation is done synchroniously, to avoid race conditions
func (a *assets) AddNodeTypes(nodeTypes ...string) {
	a.Lock()
	defer a.Unlock()
	for _, nodeType := range nodeTypes {
		if _, ok := a.items[nodeType]; !ok {
			a.items[nodeType] = nodes{}
		}
	}
}

// GetNodesOfType is used to retrieve Nodes of the given type from the available ones.
// This operation can be done synchroniously with its self, but blocks concurrent writes to avoid race conditions
func (a *assets) GetNodesByType(nodeType string) (Nodes, error) {
	a.RLock()
	defer a.RUnlock()
	nodes, ok := a.items[nodeType]
	if !ok {
		return nil, fmt.Errorf("%w; node type %s", ErrNotFound, nodeType)
	}
	return nodes, nil
}

// ListNodeTypes returns the types of nodes available in the assets
func (a *assets) ListNodeTypes() []string {
	a.RLock()
	defer a.RUnlock()
	nodeTypes := make([]string, len(a.items))
	i := 0
	for k := range a.items {
		nodeTypes[i] = k
		i++
	}
	return nodeTypes
}

// New creates a new empty Graph structure to be used to store new entities and establish relationships between them
func New() Graph {
	return &assets{
		items: map[string]nodes{},
	}
}
