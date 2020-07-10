package main

import (
        "errors"
        "fmt"
        "log"
        "os"
        "os/exec"
        "strings"

        "github.com/urfave/cli/v2" // imports as package "cli"
)

var (
        // build number set at compile-time
        build = "0"
)

func main() {
        app := cli.NewApp()
        app.Name = "akslabs"
        app.Usage = "Azure Kubernetes Service Labs"
        app.Version = fmt.Sprintf("0.1.%s", build)
        app.Flags = []cli.Flag{
                &cli.StringFlag{
                        Name:    "name",
                        Aliases: []string{"n"},
                        Usage:   "Name of the managed cluster",
                        Value:   "lab-cluster",
                },
                &cli.StringFlag{
                        Name:    "resource-group",
                        Aliases: []string{"g"},
                        Usage:   "Name of resource group",
                        Value:   "lab-resource-group",
                },
                &cli.StringFlag{
                        Name:    "location",
                        Aliases: []string{"l"},
                        Usage:   "Location. Values from: `az account list-locations`",
                        Value:   "westus",
                },
                &cli.StringFlag{
                        Name:    "node-vm-size",
                        Aliases: []string{"s"},
                        Usage:   " Size of Virtual Machines to create as Kubernetes nodes",
                        Value:   "Standard_B2s",
                },
                &cli.StringFlag{
                        Name:    "node-count",
                        Aliases: []string{"c"},
                        Usage:   "Number of nodes in the Kubernetes node pool",
                        Value:   "1",
                },
        }
        app.Commands = []*cli.Command{
                {
                        Name:  "deploy",
                        Usage: "deploy a lab exercise",
                        Subcommands: []*cli.Command{
                                {
                                        Name:   "lab1",
                                        Usage:  "deploy lab1",
                                        Action: deployLab1,
                                },
                        },
                },
                {
                        Name:  "validate",
                        Usage: "validate the lab results",
                        Subcommands: []*cli.Command{
                                {
                                        Name:   "lab1",
                                        Usage:  "validate lab1",
                                        Action: validateLab1,
                                },
                        }},
                {
                        Name:   "describe",
                        Usage:  "describe all labs",
                        Action: describeLabs,
                },
        }
        if err := app.Run(os.Args); err != nil {
                log.Fatal(err)
        }
}

func validate(c *cli.Context) error {
        fmt.Println("vaidate a lab: " + c.String("name"))
        return nil
}

func azLoginCheck() error {
        cmd := exec.Command("az", "account", "list")
        fmt.Println(strings.Join(cmd.Args, " "))
        stdoutStderr, err := cmd.CombinedOutput()
        if err != nil {
                return err
        }
        output := fmt.Sprintf("%s\n", stdoutStderr)
        if strings.Contains(output, "az login") {
                output := fmt.Sprintf("%s\n", stdoutStderr)
                return errors.New(output)
        }
        return nil
}

func createResourcegroup(resourceGroup string, location string) error {
        cmd := exec.Command("az", "group", "create", "--name", resourceGroup, "--location", location)
        fmt.Println(strings.Join(cmd.Args, " "))
        stdoutStderr, err := cmd.CombinedOutput()
        fmt.Printf("%s\n", stdoutStderr)
        return err
}

// check resource group and cluster
// Resource group 'aks1011' could not be found.
func resourcegroupExists(resourceGroup string) bool {
        cmd := exec.Command("az", "group", "show", "-g", resourceGroup)
        fmt.Println(strings.Join(cmd.Args, " "))
        stdoutStderr, _ := cmd.CombinedOutput()
        output := fmt.Sprintf("%s\n", stdoutStderr)
        fmt.Println(output)
        if strings.Contains(output, fmt.Sprintf("Resource group '%s' could not be found.", resourceGroup)) {
                return false
        }
        return true
}

// validate cluster exists
// The Resource 'Microsoft.ContainerService/managedClusters/demo11' under resource group 'aks101' was not found. For more details please go to https://aka.ms/ARMResourceNotFoundFix
// Resource group 'aks1011' could not be found.
func clusterExists(resourceGroup string, name string) bool {

        cmd := exec.Command("az", "aks", "show", "-g", resourceGroup, "-n", name)
        fmt.Println(strings.Join(cmd.Args, " "))
        stdoutStderr, _ := cmd.CombinedOutput()
        output := fmt.Sprintf("%s\n", stdoutStderr)
        fmt.Println(output)
        if strings.Contains(output, "For more details please go to https://aka.ms/ARMResourceNotFoundFix") {
                return false
        }
        if strings.Contains(output, fmt.Sprintf("Resource group '%s' could not be found.", resourceGroup)) {
                return false
        }
        return true
}

func describeLabs(c *cli.Context) error {
        fmt.Println("Note")
        fmt.Println("-------------")
        fmt.Println("Please run az login before the labs")
        fmt.Println("Please install kubectl")
        fmt.Println("recommanded to create a new resource group for each lab")
        fmt.Println("")
        fmt.Println("Lab 1 Networking")
        fmt.Println("----------------")
        fmt.Println("Run 'akslabs -n lab1cluster -g lab1group deploy lab1' to deploy the lab1 cluster.")
        fmt.Println("Tasks:")
        fmt.Println("  Run 'kubectl create deployment whoami --image=containous/whoami' to deploy the pod.")
        fmt.Println("  Run 'kubectl get pods' to list pod.")
        fmt.Println("  Run 'akslabs -n lab1cluster -g lab1group validate lab1' to check the results.")
        return nil
}

// Lab  1
func deployLab1(c *cli.Context) error {
        fmt.Println("Deploying cluster for lab1...")

        rgName := c.String("resource-group")

        if err := azLoginCheck(); err != nil {
                return err
        }

        if resourcegroupExists(rgName) {
                return fmt.Errorf("Please delete the resource group '%s' first", rgName)
        }

        if err := createResourcegroup(rgName, c.String("location")); err != nil {
                return err
        }

        // 1. create aks vnet / subnet
        cmd1 := exec.Command(
                "az", "network", "vnet", "create",
                "--resource-group", rgName,
                "--name", "aksvnet",
                "--location", c.String("location"),
                "--address-prefix", "10.0.0.0/8",
                "--subnet-name", "default",
                "--subnet-prefix", "10.240.0.0/16",
        )

        cmd1out, err := cmd1.CombinedOutput()
        if err != nil {
                fmt.Printf("%s\n", cmd1out)
                return err
        }

        // 2. get subnet ID
        cmd2 := exec.Command(
                "az", "network", "vnet", "subnet", "list",
                "--resource-group", rgName,
                "--vnet-name", "aksvnet",
                "--query", "[0].id",
                "--output", "tsv",
        )

        subnetIDByte, err := cmd2.CombinedOutput()
        if err != nil {
                return err
        }

        // 3. create AKS cluster and add authroized ip range
        cmd3 := exec.Command(
                "az", "aks", "create",
                "--resource-group", rgName,
                "--name", c.String("name"),
                "--location", c.String("location"),
                "--node-vm-size", c.String("node-vm-size"),
                "--node-count", c.String("node-count"),
                "--network-plugin", "azure",
                "--service-cidr", "10.0.0.0/16",
                "--vnet-subnet-id", strings.TrimSpace(fmt.Sprintf("%s", subnetIDByte)),
                "--api-server-authorized-ip-ranges", "193.168.1.0/24",
        )
        // fmt.Println(strings.Join(cmd3.Args, " "))
        stdoutStderr, err := cmd3.CombinedOutput()
        fmt.Printf("%s\n", stdoutStderr)
        if err != nil {
                return err
        }

        // 4. Create a new subnet to conflict with the service CIDR
        cmd4 := exec.Command(
                "az", "network", "vnet", "subnet", "create",
                "--resource-group", rgName,
                "--vnet-name", "aksvnet",
                "--name", "gateway",
                "--address-prefix", "10.0.0.0/24",
        )

        cmd4out, err := cmd4.CombinedOutput()
        if err != nil {
                fmt.Printf("%s\n", cmd4out)
                return err
        }

        fmt.Println("lab1 ready to go")
        return nil
}

func validateLab1(c *cli.Context) error {

        cmd := exec.Command("kubectl", "get", "po", "-l", "app=whoami", "--no-headers")
        cmdout, err := cmd.CombinedOutput()
        if err != nil {
                fmt.Printf("Lab1 validation failed!\n%s", cmdout)
                return err
        }
        // expected output
        // whoami-6c79b8c8d-ztgzc   1/1   Running   0     9d

        output := fmt.Sprintf("%s", cmdout)

        if strings.HasPrefix(output, "whoami-") && strings.Contains(output, "1/1   Running   0") {
                fmt.Println("Congratulations! You have passed the lab1 test. Good job!")
                return nil
        }
        return fmt.Errorf("Lab1 validation failed!\n%s", output)
}
