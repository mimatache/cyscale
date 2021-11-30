package assets

import (
	"encoding/json"
	"fmt"

	"github.com/mimatache/cyscale/internal/graph"
)

const (
	InterfaceType     = "interface"
	VpcType           = "vpc"
	SecurityGroupType = "securityGroup"
	VirtualMacineType = "vm"
)

func NewManager(graph *graph.Graph, vpcData, sgData, interfaceData, vmData []byte) (*Manager, error) {
	m := &Manager{
		graph: graph,
		items: map[string]string{},
	}
	if err := m.loadVPCs(vpcData); err != nil {
		return nil, err
	}
	if err := m.loadSGs(sgData); err != nil {
		return nil, err
	}
	if err := m.loadInterfaces(interfaceData); err != nil {
		return nil, err
	}
	if err := m.loadVMs(vmData); err != nil {
		return nil, err
	}
	return m, nil
}

type Manager struct {
	graph *graph.Graph
	items map[string]string
}

func (m *Manager) loadInterfaces(data []byte) error {
	interfaces := []Interface{}
	if err := json.Unmarshal(data, &interfaces); err != nil {
		return fmt.Errorf("could not unmarshal interfaces; %w", err)
	}
	for _, v := range interfaces {
		interfaceBody, err := json.Marshal(v)
		if err != nil {
			return fmt.Errorf("could not marshal interface; %w", err)
		}
		node := m.graph.InsertNode(v.NetworkInterfaceID, InterfaceType, interfaceBody)

		vpcs := m.graph.ListNodes(graph.FilterNodesByName(v.VpcID))
		if len(vpcs) == 0 {
			n := m.graph.InsertNode(v.VpcID, VpcType, []byte{})
			vpcs = append(vpcs, n)
		}
		for _, vpc := range vpcs {
			if _, err := m.graph.AddRelationship(node.GetID(), vpc.GetID(), "part_of"); err != nil {
				return err
			}
		}

		for _, sg := range v.SecurityGroupIDs {
			sgs := m.graph.ListNodes(graph.FilterNodesByName(sg))
			if len(sgs) == 0 {
				n := m.graph.InsertNode(sg, SecurityGroupType, []byte{})
				sgs = append(sgs, n)
			}
			for _, sgNode := range sgs {
				if _, err := m.graph.AddRelationship(node.GetID(), sgNode.GetID(), "part_of"); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (m *Manager) loadVPCs(data []byte) error {
	vpcs := []VirtualPrivateCloud{}
	if err := json.Unmarshal(data, &vpcs); err != nil {
		return fmt.Errorf("could not unmarshal interfaces; %w", err)
	}
	for _, v := range vpcs {
		vpcBody, err := json.Marshal(v)
		if err != nil {
			return fmt.Errorf("could not marshal interface; %w", err)
		}
		m.graph.InsertNode(v.VpcID, VpcType, vpcBody)
	}
	return nil
}

func (m *Manager) loadVMs(data []byte) error {
	vms := []VirtualMachine{}
	if err := json.Unmarshal(data, &vms); err != nil {
		return fmt.Errorf("could not unmarshal vms; %w", err)
	}
	for _, v := range vms {
		vmBody, err := json.Marshal(v)
		if err != nil {
			return fmt.Errorf("could not marshal vms; %w", err)
		}
		node := m.graph.InsertNode(v.Name, VirtualMacineType, vmBody)

		vpcs := m.graph.ListNodes(graph.FilterNodesByName(v.VpcID))
		if len(vpcs) == 0 {
			n := m.graph.InsertNode(v.VpcID, VpcType, []byte{})
			vpcs = append(vpcs, n)
		}
		for _, vpc := range vpcs {
			if _, err := m.graph.AddRelationship(node.GetID(), vpc.GetID(), "part_of"); err != nil {
				return err
			}
		}

		for _, sg := range v.SecurityGroupIDs {
			sgs := m.graph.ListNodes(graph.FilterNodesByName(sg))
			if len(sgs) == 0 {
				n := m.graph.InsertNode(sg, SecurityGroupType, []byte{})
				sgs = append(sgs, n)
			}
			for _, sgNode := range sgs {
				if _, err := m.graph.AddRelationship(node.GetID(), sgNode.GetID(), "part_of"); err != nil {
					return err
				}
			}
		}

		for _, intfID := range v.NetworkInterfaceIDs {
			intfs := m.graph.ListNodes(graph.FilterNodesByName(intfID))
			if len(intfs) == 0 {
				n := m.graph.InsertNode(intfID, SecurityGroupType, []byte{})
				intfs = append(intfs, n)
			}
			for _, intfNode := range intfs {
				if _, err := m.graph.AddRelationship(node.GetID(), intfNode.GetID(), "using"); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (m *Manager) loadSGs(data []byte) error {
	sgs := []SecurityGroup{}
	if err := json.Unmarshal(data, &sgs); err != nil {
		return fmt.Errorf("could not unmarshal sgs; %w", err)
	}
	for _, v := range sgs {
		sgBody, err := json.Marshal(v)
		if err != nil {
			return fmt.Errorf("could not marshal sgs; %w", err)
		}
		node := m.graph.InsertNode(v.GroupID, SecurityGroupType, sgBody)

		vpcs := m.graph.ListNodes(graph.FilterNodesByName(v.VpcID))
		if len(vpcs) == 0 {
			n := m.graph.InsertNode(v.VpcID, VpcType, []byte{})
			vpcs = append(vpcs, n)
		}
		for _, vpc := range vpcs {
			if _, err := m.graph.AddRelationship(node.GetID(), vpc.GetID(), "part_of"); err != nil {
				return err
			}
		}
	}
	return nil
}
