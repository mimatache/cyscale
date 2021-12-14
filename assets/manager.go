package assets

import (
	"encoding/json"
	"fmt"
	"log"

	graph "github.com/curious-kitten/assets"
)

const (
	InterfaceType     = "interface"
	VpcType           = "vpc"
	SecurityGroupType = "securityGroup"
	VirtualMacineType = "vm"
)

// NewManager creates a new instance of an asset manager, allong with loading the data that will be used by it
func NewManager(graph *graph.Graph, vpcData, sgData, interfaceData, vmData []byte) (*Manager, error) {
	m := &Manager{
		graph: graph,
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
			if _, err := m.graph.AddRelationship(node, vpc, "part_of"); err != nil {
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
				if _, err := m.graph.AddRelationship(node, sgNode, "part_of"); err != nil {
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
			if _, err := m.graph.AddRelationship(node, vpc, "part_of"); err != nil {
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
				if _, err := m.graph.AddRelationship(node, sgNode, "part_of"); err != nil {
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
				if _, err := m.graph.AddRelationship(node, intfNode, "using"); err != nil {
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
			if _, err := m.graph.AddRelationship(node, vpc, "part_of"); err != nil {
				return err
			}
		}
	}
	return nil
}

// ListExposedVMs returns a list of VMs that accept connections from 0.0.0.0/0
func (m *Manager) ListExposedVMs() []string {
	rule := func(node graph.Node) bool {
		sg := SecurityGroup{}
		if err := json.Unmarshal(node.Body, &sg); err != nil {
			// only printing the error, since an error here means there is no useful information to extract, but we still need to continue checking
			// TODO: consider adding a check for invalid bodies?
			log.Printf("error: unable to unmarshal security group %s; %s; this might indicate corrupt data \n", node.GetName(), err.Error())
		}
		for _, network := range sg.IPList {
			if network == "0.0.0.0/0" && sg.Direction == "inbound" {
				return true
			}
		}
		return false
	}
	return m.findVMsBySecurityIssue(rule)
}

// ListHTTPPortVMs returns a list of VMs that have port 80 opened, either directly on the VM, or on a connected interface
func (m *Manager) ListHTTPPortVMs() []string {
	rule := func(node graph.Node) bool {
		sg := SecurityGroup{}
		if err := json.Unmarshal(node.Body, &sg); err != nil {
			// only printing the error, since an error here means there is no useful information to extract, but we still need to continue checking
			// TODO: consider adding a check for invalid bodies?
			log.Printf("error: unable to unmarshal security group %s; %s; this might indicate corrupt data \n", node.GetName(), err.Error())
		}
		for _, port := range sg.ExposedPorts {
			if port == 80 && sg.Direction == "inbound" {
				return true
			}
		}
		return false
	}
	return m.findVMsBySecurityIssue(rule)
}

// ListConnections list all possible relationship chains between the 2 points
func (m *Manager) ListConnections(from, to string) ([]string, error) {
	fromNodes := m.graph.ListNodes(graph.FilterNodesByName(from))
	if len(fromNodes) != 1 {
		return []string{}, fmt.Errorf("could not uniquely identify node %s; found %d elements", from, len(fromNodes))
	}
	toNodes := m.graph.ListNodes(graph.FilterNodesByName(to))
	if len(fromNodes) != 1 {
		return []string{}, fmt.Errorf("could not uniquely identify node %s; found %d elements", to, len(toNodes))
	}
	chains := m.graph.ListConnections(fromNodes[0], toNodes[0])
	connections := make([]string, len(chains))
	for i, v := range chains {
		connections[i] = v.String()
	}
	return connections, nil
}

// findVMsBySecurityIssue searches for VMs that have connections to a SecurityGroup that is in violation of the given rule
func (m *Manager) findVMsBySecurityIssue(rule graph.FilterNodes) []string {
	exposedVMs := []string{}
	// get security groups that habe port 80 opened
	openedSecurityGroups := m.graph.ListNodes(
		graph.FilterNodesByLabel(SecurityGroupType),
		rule)
	vms := m.graph.ListNodes(graph.FilterNodesByLabel(VirtualMacineType))
	for _, sg := range openedSecurityGroups {
		for _, vm := range vms {
			if cons := m.graph.ListConnections(vm, sg); len(cons) > 0 {
				exposedVMs = append(exposedVMs, vm.GetName())
			}
		}
	}
	return exposedVMs
}
