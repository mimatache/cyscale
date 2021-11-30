package assets_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mimatache/cyscale/assets"
	"github.com/mimatache/cyscale/internal/graph"
)

func Test_NewManager(t *testing.T) {
	intfContents, err := os.ReadFile("testdata/NetworkInterface.json")
	assert.NoError(t, err, "error reading files")

	vmContents, err := os.ReadFile("testdata/VM.json")
	assert.NoError(t, err, "error reading files")

	vpcContents, err := os.ReadFile("testdata/VPC.json")
	assert.NoError(t, err, "error reading files")

	sgContents, err := os.ReadFile("testdata/SecurityGroup.json")
	assert.NoError(t, err, "error reading files")

	grf := graph.New()
	_, err = assets.NewManager(grf, vpcContents, sgContents, intfContents, vmContents)

	assert.NoError(t, err)

	assert.Equal(t, 11, len(grf.ListNodes()))
	assert.Equal(t, 18, len(grf.ListRelationships()))
}
