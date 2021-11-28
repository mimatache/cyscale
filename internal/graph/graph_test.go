package graph_test

import (
	"errors"
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mimatache/cyscale/internal/graph"
)

const (
	dragonType  = "dragon"
	unicornType = "unicorn"
	puppyType   = "puppy"
)

var tests = []struct {
	name          string
	toAdd         []string
	shouldFind    []string
	shouldNotFind []string
}{
	{
		name:          "add one",
		toAdd:         []string{puppyType},
		shouldFind:    []string{puppyType},
		shouldNotFind: []string{unicornType, dragonType},
	},
	{
		name:          "add multiple",
		toAdd:         []string{puppyType},
		shouldFind:    []string{puppyType, unicornType},
		shouldNotFind: []string{dragonType},
	},
}

func Test_Graph_AddTypes(t *testing.T) {
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			grf := graph.New()
			grf.AddNodeTypes(test.shouldFind...)
			nodeTypes := grf.ListNodeTypes()
			assert.Equal(t, len(test.shouldFind), len(nodeTypes), "missmatch between expected type count and found type count")
		})
	}
}
func Test_Graph_AddTypesConcurrently(t *testing.T) {
	concurrencyCount := 100
	types := make([]string, concurrencyCount)
	for i := 0; i < concurrencyCount; i++ {
		types[i] = fmt.Sprintf("type-%d", i)
	}
	var wg sync.WaitGroup
	wg.Add(concurrencyCount)
	grf := graph.New()
	for i := 0; i < concurrencyCount; i++ {
		go func(index int) {
			defer wg.Done()
			grf.AddNodeTypes(fmt.Sprintf("type-%d", index))
		}(i)
	}
	wg.Wait()
	assert.Equal(t, concurrencyCount, len(grf.ListNodeTypes()), "not all node types were added")
}
func Test_Graph_GetByType(t *testing.T) {
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			grf := graph.New()
			grf.AddNodeTypes(test.shouldFind...)
			for _, nodeType := range test.shouldFind {
				_, err := grf.GetNodesByType(nodeType)
				assert.NoError(t, err, "failed to find type")
			}
		})
	}
}
func Test_Graph_GetByType_Missing(t *testing.T) {
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			grf := graph.New()
			grf.AddNodeTypes(test.shouldFind...)
			for _, nodeType := range test.shouldNotFind {
				_, err := grf.GetNodesByType(nodeType)
				assert.True(t, errors.Is(err, graph.ErrNotFound), "returned error did not match expected")
			}
		})
	}
}
