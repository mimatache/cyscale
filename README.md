# cyscale
The purpose of this project is to find basic security misconfigurations, given the provided input.

## Usage
` cyscale-cli `:
```
This is a simple application that parses json files containing information about cloud topologies and performs simple security violation scans.
This should not be treated as an exhaustive security scan of your cloud environment

Usage:
  cyscale-cli [command]

Available Commands:
  completion  generate the autocompletion script for the specified shell
  help        Help about any command
  license     Show license information
  verify      verify is used to check conditions
  version     Show version information

Flags:
  -h, --help   help for cyscale-cli

Use "cyscale-cli [command] --help" for more information about a command.
```

## Checks
`cyscale-cli verify`:
```
verify is used to check conditions

Usage:
  cyscale-cli verify [command]

Available Commands:
  exposed-vms         exposed-vms shows which VMs are exposed to the internet (i.e.: allow connections from 0.0.0.0/0)
  list-connections    list-connections shows how two assets connect to each other. Example `list-connections intf1 vpc1`
  vms-using-http-port vms-using-http-port shows which VMs are using the HTTP port, either directly or through an interface

Flags:
  -h, --help                           help for verify
      --interfaces string              path to file containing network interfaces to verify (default "data/NetworkInterface.json")
      --security-groups string         path to file containing security groups to verify (default "data/SecurityGroup.json")
      --virtual-machines string        path to file containing VMs to verify (default "data/VM.json")
      --virtual-private-cloud string   path to file containing VPCs to verify (default "data/VPC.json")

Use "cyscale-cli verify [command] --help" for more information about a command.
```

## Building the binary
`make` -> this will build a binary in `{PROJECT}/.build/cyscale-cli/_bin

## Building a docker image
`make docker-build`