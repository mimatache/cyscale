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

	verifyCommand.PersistentFlags().StringVar(&interfaces, "interfaces", "testdata/NetworkInterface.json", "path to file containing network interfaces to verify")
	verifyCommand.PersistentFlags().StringVar(&vms, "virtual-machines", "testdata/VM.json", "path to file containing VMs to verify")
	verifyCommand.PersistentFlags().StringVar(&sgs, "security-groups", "testdata/SecurityGroup.json", "path to file containing security groups to verify")
	verifyCommand.PersistentFlags().StringVar(&vpcs, "virtual-private-cloud", "testdata/VPC.json", "path to file containing VPCs to verify")

	verifyCommand.AddCommand(
		exposedVMCommand,
	)

	return verifyCommand
}

var exposedVMCommand = &cobra.Command{
	Use:   "exposed-vms",
	Short: "exposed-vms shows which VMs are exposed to the internet",
	RunE: func(cmd *cobra.Command, args []string) error {
		_, err := getAssetManager(interfaces, vms, sgs, vpcs)
		if err != nil {
			return fmt.Errorf("could not load data; %w", err)
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
