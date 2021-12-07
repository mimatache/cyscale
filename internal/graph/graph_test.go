package graph_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mimatache/cyscale/internal/graph"
)

const (
	dragonType = "dragon"
	puppyType  = "puppy"
)

type puppy struct {
	Power int    `json:"power"`
	Name  string `json:"name"`
}

var (
	bobita     = "Bobita"
	bobitaBody = []byte("{\"name\":\"Bobita\", \"power\":500}")
	azor       = "Azor"
	azorBody   = []byte("{\"name\":\"Azor\", \"power\":457}")
	smaug      = "Smaug"
	smaugBody  = []byte("{\"name\":\"Azor\", \"power\":457, \"canFly\":true}")
)

func Test_Graph_Insert(t *testing.T) {
	grf := graph.New()
	createdNode := grf.InsertNode(bobita, puppyType, bobitaBody)
	node, err := grf.GetNodeByID(createdNode.GetID())
	assert.NoError(t, err)
	assert.Equal(t, bobitaBody, node.Body)
	assert.Equal(t, puppyType, node.GetLabel())
}

func Test_Graph_Insert_Imutable(t *testing.T) {
	grf := graph.New()
	createdNode := grf.InsertNode(bobita, puppyType, bobitaBody)
	node, err := grf.GetNodeByID(createdNode.GetID())
	node.Body = []byte{}
	assert.NoError(t, err)
	assert.NotEqual(t, bobitaBody, node.Body)
}

func Test_Graph_AddConcurrently(t *testing.T) {
	concurrencyCount := 100
	types := make([]string, concurrencyCount)
	for i := 0; i < concurrencyCount; i++ {
		types[i] = fmt.Sprintf("type-%d", i)
	}
	var wg sync.WaitGroup
	wg.Add(concurrencyCount)
	grf := graph.New()
	for i := 0; i < concurrencyCount; i++ {
		go func(i int) {
			defer wg.Done()
			grf.InsertNode(fmt.Sprintf("item-%d", i), puppyType, []byte{})
		}(i)
	}
	wg.Wait()
	nodes := grf.ListNodes()
	assert.Equal(t, concurrencyCount, len(nodes), "not all node types were added")
}

func Test_Graph_AddRelationshipConcurrently(t *testing.T) {
	concurrencyCount := 100
	types := make([]string, concurrencyCount)
	grf := graph.New()
	createdNodeOne := grf.InsertNode(bobita, puppyType, bobitaBody)
	createdNodeTwo := grf.InsertNode(azor, puppyType, azorBody)
	for i := 0; i < concurrencyCount; i++ {
		types[i] = fmt.Sprintf("type-%d", i)
	}
	var wg sync.WaitGroup
	wg.Add(concurrencyCount)
	for i := 0; i < concurrencyCount; i++ {
		go func(i int) {
			defer wg.Done()
			_, err := grf.AddRelationship(createdNodeOne.GetID(), createdNodeTwo.Copy().GetID(), fmt.Sprintf("item-%d", i))
			assert.NoError(t, err)
		}(i)
	}
	wg.Wait()
	rels := grf.ListRelationships()
	assert.Equal(t, concurrencyCount, len(rels), "not all node types were added")
}
func Test_Graph_GetNodes_Missing(t *testing.T) {
	grf := graph.New()
	_, err := grf.GetNodeByID("bobitaNodeID")
	assert.True(t, errors.Is(err, graph.ErrNotFound))
}

func Test_Graph_List(t *testing.T) {
	grf := graph.New()
	grf.InsertNode(bobita, puppyType, bobitaBody)
	grf.InsertNode(azor, puppyType, azorBody)
	grf.InsertNode(smaug, dragonType, smaugBody)
	foundNodes := grf.ListNodes()
	assert.Equal(t, 3, len(foundNodes))
}

func Test_Graph_ListNodes_FilterByLabel(t *testing.T) {
	grf := graph.New()
	bNode := grf.InsertNode(bobita, puppyType, bobitaBody)
	grf.InsertNode(smaug, dragonType, smaugBody)
	foundNodes := grf.ListNodes(graph.FilterNodesByLabel(puppyType))
	assert.Equal(t, 1, len(foundNodes))
	assert.Equal(t, bNode.GetID(), foundNodes[0].GetID())
}

func Test_Graph_ListNodes_FilterByName(t *testing.T) {
	grf := graph.New()
	bNode := grf.InsertNode(bobita, puppyType, bobitaBody)
	grf.InsertNode(smaug, dragonType, smaugBody)
	foundNodes := grf.ListNodes(graph.FilterNodesByName(bobita))
	assert.Equal(t, 1, len(foundNodes))
	assert.Equal(t, bNode.GetID(), foundNodes[0].GetID())
}

func Test_Graph_ListNodes_Filter(t *testing.T) {
	grf := graph.New()
	whereCond := func(body *graph.Node) bool {
		pup := puppy{}
		if err := json.Unmarshal(body.Body, &pup); err != nil {
			return false
		}
		return pup.Power > 499
	}
	bNode := grf.InsertNode(bobita, puppyType, bobitaBody)
	grf.InsertNode(azor, puppyType, azorBody)
	foundNodes := grf.ListNodes(whereCond)
	assert.Equal(t, 1, len(foundNodes))
	assert.Equal(t, bNode.GetID(), foundNodes[0].GetID())
}

func Test_Graph_AddRelationship(t *testing.T) {
	grf := graph.New()
	bNode := grf.InsertNode(bobita, puppyType, bobitaBody)
	aNode := grf.InsertNode(azor, puppyType, azorBody)
	rel1, err := grf.AddRelationship(bNode.GetID(), aNode.GetID(), "friends")
	assert.NoError(t, err)
	rel2, err := grf.AddRelationship(bNode.GetID(), aNode.GetID(), "competitors")
	assert.NoError(t, err)
	bNode, err = grf.GetNodeByID(bNode.GetID())
	assert.NoError(t, err)
	assert.Equal(t, 2, len(bNode.ListRelationships()))
	assert.Contains(t, bNode.ListRelationships(), rel1.ID)
	assert.Contains(t, bNode.ListRelationships(), rel2.ID)
}

func Test_Graph_AddRelationship_NoFrom(t *testing.T) {
	grf := graph.New()
	aNode := grf.InsertNode(azor, puppyType, azorBody)
	_, err := grf.AddRelationship("bNode.GetID()", aNode.GetID(), "friends")
	assert.Error(t, err)
}

func Test_Graph_AddRelationship_NoTo(t *testing.T) {
	grf := graph.New()
	aNode := grf.InsertNode(azor, puppyType, azorBody)
	_, err := grf.AddRelationship(aNode.GetID(), "bNode.GetID()", "friends")
	assert.Error(t, err)
}

func Test_Graph_GetRelationship(t *testing.T) {
	grf := graph.New()
	bNode := grf.InsertNode(bobita, puppyType, bobitaBody)
	aNode := grf.InsertNode(azor, puppyType, azorBody)
	initialRel, err := grf.AddRelationship(bNode.GetID(), aNode.GetID(), "friends")
	assert.NoError(t, err)
	foundRel, err := grf.GetRelationshipByID(initialRel.ID)
	assert.NoError(t, err)
	assert.Equal(t, initialRel, foundRel)
}

func Test_Graph_GetRelationship_NotFound(t *testing.T) {
	grf := graph.New()
	_, err := grf.GetRelationshipByID("fake")
	assert.Error(t, err)
	assert.ErrorIs(t, err, graph.ErrNotFound)
}

func Test_Graph_ListRelationships(t *testing.T) {
	grf := graph.New()
	bNode := grf.InsertNode(bobita, puppyType, bobitaBody)
	aNode := grf.InsertNode(azor, puppyType, azorBody)
	dNode := grf.InsertNode(smaug, dragonType, smaugBody)
	rel1, err := grf.AddRelationship(bNode.GetID(), aNode.GetID(), "friends")
	assert.NoError(t, err)
	rel2, err := grf.AddRelationship(bNode.GetID(), aNode.GetID(), "competitors")
	assert.NoError(t, err)
	rel3, err := grf.AddRelationship(bNode.GetID(), dNode.GetID(), "enemies")
	assert.NoError(t, err)
	rels := grf.ListRelationships()
	assert.Equal(t, 3, len(rels))
	assert.Contains(t, rels, rel1)
	assert.Contains(t, rels, rel2)
	assert.Contains(t, rels, rel3)
}

func Test_Graph_ListRelationships_Filter(t *testing.T) {
	grf := graph.New()
	bNode := grf.InsertNode(bobita, puppyType, bobitaBody)
	aNode := grf.InsertNode(azor, puppyType, azorBody)
	dNode := grf.InsertNode(smaug, dragonType, smaugBody)
	rel1, err := grf.AddRelationship(bNode.GetID(), aNode.GetID(), "friends")
	assert.NoError(t, err)
	_, err = grf.AddRelationship(bNode.GetID(), aNode.GetID(), "competitors")
	assert.NoError(t, err)
	_, err = grf.AddRelationship(bNode.GetID(), dNode.GetID(), "enemies")
	assert.NoError(t, err)
	rels := grf.ListRelationships(graph.FilterRelByLabel("friends"))
	assert.Equal(t, 1, len(rels))
	assert.Contains(t, rels, rel1)
}

func Test_Graph_ListRelationships_FilterByFrom(t *testing.T) {
	grf := graph.New()
	bNode := grf.InsertNode(bobita, puppyType, bobitaBody)
	aNode := grf.InsertNode(azor, puppyType, azorBody)
	dNode := grf.InsertNode(smaug, dragonType, smaugBody)
	rel1, err := grf.AddRelationship(bNode.GetID(), aNode.GetID(), "friends")
	assert.NoError(t, err)
	_, err = grf.AddRelationship(dNode.GetID(), bNode.GetID(), "enemies")
	assert.NoError(t, err)
	rels := grf.ListRelationships(graph.FilterRelByFrom(bNode.GetID()))
	assert.Equal(t, 1, len(rels))
	assert.Contains(t, rels, rel1)
}

func Test_Graph_ListRelationships_FilterByTo(t *testing.T) {
	grf := graph.New()
	bNode := grf.InsertNode(bobita, puppyType, bobitaBody)
	aNode := grf.InsertNode(azor, puppyType, azorBody)
	dNode := grf.InsertNode(smaug, dragonType, smaugBody)
	rel1, err := grf.AddRelationship(bNode.GetID(), aNode.GetID(), "friends")
	assert.NoError(t, err)
	_, err = grf.AddRelationship(dNode.GetID(), bNode.GetID(), "enemies")
	assert.NoError(t, err)
	rels := grf.ListRelationships(graph.FilterRelByTo(aNode.GetID()))
	assert.Equal(t, 1, len(rels))
	assert.Contains(t, rels, rel1)
}

func Test_Graph_ListConnections(t *testing.T) {
	grf := graph.New()
	bNode := grf.InsertNode(bobita, puppyType, bobitaBody)
	aNode := grf.InsertNode(azor, puppyType, azorBody)
	dNode := grf.InsertNode(smaug, dragonType, smaugBody)
	_, err := grf.AddRelationship(bNode.GetID(), aNode.GetID(), "friends")
	assert.NoError(t, err)
	_, err = grf.AddRelationship(bNode.GetID(), dNode.GetID(), "enemies")
	assert.NoError(t, err)
	_, err = grf.AddRelationship(aNode.GetID(), dNode.GetID(), "enemies")
	assert.NoError(t, err)
	_, err = grf.AddRelationship(aNode.GetID(), bNode.GetID(), "friends")
	assert.NoError(t, err)
	_, err = grf.AddRelationship(aNode.GetID(), bNode.GetID(), "friends")
	assert.NoError(t, err)
	bNode, err = grf.GetNodeByID(bNode.GetID())
	assert.NoError(t, err)
	dNode, err = grf.GetNodeByID(dNode.GetID())
	assert.NoError(t, err)
	cons := grf.ListConnections(bNode, dNode)
	assert.Len(t, cons, 2)
}
