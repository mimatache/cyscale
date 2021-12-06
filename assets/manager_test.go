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

func Test_ExposedVMs(t *testing.T) {
	intfContents, err := os.ReadFile("testdata/NetworkInterface.json")
	assert.NoError(t, err, "error reading files")

	vmContents, err := os.ReadFile("testdata/VM.json")
	assert.NoError(t, err, "error reading files")

	vpcContents, err := os.ReadFile("testdata/VPC.json")
	assert.NoError(t, err, "error reading files")

	sgContents, err := os.ReadFile("testdata/SecurityGroup.json")
	assert.NoError(t, err, "error reading files")

	grf := graph.New()
	m, err := assets.NewManager(grf, vpcContents, sgContents, intfContents, vmContents)
	assert.NoError(t, err)

	vms := m.ListExposedVMs()
	assert.Equal(t, 2, len(vms))
	assert.Contains(t, vms, "VM_1")
	assert.Contains(t, vms, "VM_2")
}

func Test_ListHTTPPortVMs(t *testing.T) {
	intfContents, err := os.ReadFile("testdata/NetworkInterface.json")
	assert.NoError(t, err, "error reading files")

	vmContents, err := os.ReadFile("testdata/VM.json")
	assert.NoError(t, err, "error reading files")

	vpcContents, err := os.ReadFile("testdata/VPC.json")
	assert.NoError(t, err, "error reading files")

	sgContents, err := os.ReadFile("testdata/SecurityGroup.json")
	assert.NoError(t, err, "error reading files")

	grf := graph.New()
	m, err := assets.NewManager(grf, vpcContents, sgContents, intfContents, vmContents)
	assert.NoError(t, err)

	vms := m.ListHTTPPortVMs()
	assert.Equal(t, 1, len(vms))
	assert.Contains(t, vms, "VM_1")
}

func Test_ListConnections(t *testing.T) {
	intfContents, err := os.ReadFile("testdata/NetworkInterface.json")
	assert.NoError(t, err, "error reading files")

	vmContents, err := os.ReadFile("testdata/VM.json")
	assert.NoError(t, err, "error reading files")

	vpcContents, err := os.ReadFile("testdata/VPC.json")
	assert.NoError(t, err, "error reading files")

	sgContents, err := os.ReadFile("testdata/SecurityGroup.json")
	assert.NoError(t, err, "error reading files")

	grf := graph.New()
	m, err := assets.NewManager(grf, vpcContents, sgContents, intfContents, vmContents)
	assert.NoError(t, err)

	cons, err := m.ListConnections("VM_1", "vpc-06bcacc5531641a68")
	assert.NoError(t, err)
	assert.Equal(t, 6, len(cons))
	assert.Contains(t, cons, "{Asset:VM_1}->{rel:VM_1-part_of-vpc-06bcacc5531641a68}->{Asset:vpc-06bcacc5531641a68}")
	assert.Contains(t, cons, "{Asset:VM_1}->{rel:VM_1-part_of-sg-095531efae90566d5}->{Asset:sg-095531efae90566d5}->{rel:sg-095531efae90566d5-part_of-vpc-06bcacc5531641a68}->{Asset:vpc-06bcacc5531641a68}")
	assert.Contains(t, cons, "{Asset:VM_1}->{rel:VM_1-using-eni-0c02d0e2602622897}->{Asset:eni-0c02d0e2602622897}->{rel:eni-0c02d0e2602622897-part_of-vpc-06bcacc5531641a68}->{Asset:vpc-06bcacc5531641a68}")
	assert.Contains(t, cons, "{Asset:VM_1}->{rel:VM_1-using-eni-0c02d0e2602622897}->{Asset:eni-0c02d0e2602622897}->{rel:eni-0c02d0e2602622897-part_of-sg-095531efae90566d5}->{Asset:sg-095531efae90566d5}->{rel:sg-095531efae90566d5-part_of-vpc-06bcacc5531641a68}->{Asset:vpc-06bcacc5531641a68}")
	assert.Contains(t, cons, "{Asset:VM_1}->{rel:VM_1-using-eni-0c1000541fb09e879}->{Asset:eni-0c1000541fb09e879}->{rel:eni-0c1000541fb09e879-part_of-vpc-06bcacc5531641a68}->{Asset:vpc-06bcacc5531641a68}")
	assert.Contains(t, cons, "{Asset:VM_1}->{rel:VM_1-using-eni-0c1000541fb09e879}->{Asset:eni-0c1000541fb09e879}->{rel:eni-0c1000541fb09e879-part_of-sg-095531efae90566d5}->{Asset:sg-095531efae90566d5}->{rel:sg-095531efae90566d5-part_of-vpc-06bcacc5531641a68}->{Asset:vpc-06bcacc5531641a68}")
}
