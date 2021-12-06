package verifier

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/mimatache/cyscale/assets"
	"github.com/mimatache/cyscale/internal/graph"
)

var (
	interfaces string
	vms        string
	sgs        string
	vpcs       string
)

func Verify() *cobra.Command {
	verifyCommand := &cobra.Command{
		Use:   "verify",
		Short: "verify is used to check conditions",
	}

	verifyCommand.PersistentFlags().StringVar(&interfaces, "interfaces", "data/NetworkInterface.json", "path to file containing network interfaces to verify")
	verifyCommand.PersistentFlags().StringVar(&vms, "virtual-machines", "data/VM.json", "path to file containing VMs to verify")
	verifyCommand.PersistentFlags().StringVar(&sgs, "security-groups", "data/SecurityGroup.json", "path to file containing security groups to verify")
	verifyCommand.PersistentFlags().StringVar(&vpcs, "virtual-private-cloud", "data/VPC.json", "path to file containing VPCs to verify")

	verifyCommand.AddCommand(
		exposedVMCommand,
		vmUsingHTTPPort,
		listConnections,
	)

	return verifyCommand
}

var exposedVMCommand = &cobra.Command{
	Use:   "exposed-vms",
	Short: "exposed-vms shows which VMs are exposed to the internet (i.e.: allow connections from 0.0.0.0/0)",
	RunE: func(cmd *cobra.Command, args []string) error {
		m, err := getAssetManager(interfaces, vms, sgs, vpcs)
		if err != nil {
			return fmt.Errorf("could not load assets; %w", err)
		}
		vms := m.ListExposedVMs()
		if len(vms) == 0 {
			fmt.Println("There are no exposed VMs")
			return nil
		}
		fmt.Println("Exposed VMs:")
		for _, vm := range vms {
			fmt.Printf("\t• %s\n", vm)
		}
		return nil
	},
}

var vmUsingHTTPPort = &cobra.Command{
	Use:   "vms-using-http-port",
	Short: "vms-using-http-port shows which VMs are using the HTTP port, either directly or through an interface",
	RunE: func(cmd *cobra.Command, args []string) error {
		m, err := getAssetManager(interfaces, vms, sgs, vpcs)
		if err != nil {
			return fmt.Errorf("could not load assets; %w", err)
		}
		vms := m.ListHTTPPortVMs()
		if len(vms) == 0 {
			fmt.Println("There are no VMs using the HTTP port")
			return nil
		}
		fmt.Println("VMs using HTTP port:")
		for _, vm := range vms {
			fmt.Printf("\t• %s\n", vm)
		}
		return nil
	},
}

var listConnections = &cobra.Command{
	Use:   "list-connections",
	Short: "list-connections shows how to assets connect to each other. Example `list-connections intf1 vpc1`",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			return fmt.Errorf("list-connections requires two arguments to function correctly")
		}
		m, err := getAssetManager(interfaces, vms, sgs, vpcs)
		if err != nil {
			return fmt.Errorf("could not load assets; %w", err)
		}
		connections, err := m.ListConnections(args[0], args[1])
		if err != nil {
			return err
		}
		if len(connections) == 0 {
			fmt.Printf("There are no connections between %s and %s\n", args[0], args[1])
			return nil
		}
		fmt.Printf("Connections between %s and %s\n", args[0], args[1])
		for _, conn := range connections {
			fmt.Printf("\t• %s\n", conn)
		}
		return nil
	},
}

func getAssetManager(interfaces, vms, sgs, vpcs string) (*assets.Manager, error) {
	graph := graph.New()
	interfaceContents, err := os.ReadFile(interfaces)
	if err != nil {
		return nil, fmt.Errorf("could not read interface file %s; %w", interfaces, err)
	}
	vmContents, err := os.ReadFile(vms)
	if err != nil {
		return nil, fmt.Errorf("could not read vm file %s; %w", vms, err)
	}
	sgContents, err := os.ReadFile(sgs)
	if err != nil {
		return nil, fmt.Errorf("could not read sg file %s; %w", sgs, err)
	}
	vpcContents, err := os.ReadFile(vpcs)
	if err != nil {
		return nil, fmt.Errorf("could not read vpc file %s; %w", vpcs, err)
	}
	return assets.NewManager(graph, vpcContents, sgContents, interfaceContents, vmContents)
}
