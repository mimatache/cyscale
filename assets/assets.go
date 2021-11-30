package assets

type Interface struct {
	Name               string   `json:"name"`
	NetworkInterfaceID string   `json:"networkInterfaceID"`
	SecurityGroupIDs   []string `json:"securityGroupIDs"`
	VpcID              string   `json:"vpcID"`
}

type SecurityGroup struct {
	Name         string   `json:"name"`
	GroupID      string   `json:"groupID"`
	VpcID        string   `json:"vpcID"`
	ExposedPorts []int    `json:"exposedPorts"`
	Direction    string   `json:"direction"`
	IPList       []string `json:"ipList"`
}

type VirtualMachine struct {
	Name                string   `json:"name"`
	SecurityGroupIDs    []string `json:"securityGroupIDs"`
	VpcID               string   `json:"vpcID"`
	NetworkInterfaceIDs []string `json:"networkInterfaceIDs"`
}

type VirtualPrivateCloud struct {
	Name  string `json:"name"`
	VpcID string `json:"vpcID"`
}
