/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package create

import (
	"context"
	"fmt"
	"net"
	"os/exec"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/spf13/cobra"
)

// controllerCmd represents the controller command
var controllerCmd = &cobra.Command{
	Use:   "controller",
	Short: "A brief description of your command",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("create controller called")
		name, _ := cmd.Flags().GetString("name")
		version, _ := cmd.Flags().GetString("ver")
		cidr, _ := cmd.Flags().GetString("cidr")
		email, _ := cmd.Flags().GetString("email")
		pwd, _ := cmd.Flags().GetString("password")
		lic, _ := cmd.Flags().GetString("customer-id")
		dt := time.Now().Format("01-02-2006 15:04:05")
		ctx := context.Background()

		flags := map[string]string{
			"name":    name,
			"version": version,
			"cidr":    cidr,
			"email":   email,
			"pwd":     pwd,
			"lic":     lic,
		}

		isFlagEmpty(flags)
		color.Green(dt + ": Validating CIDR")
		if !isIPv4CIDR(flags["cidr"]) {
			color.Magenta("%s is not a valid IPv4 CIDR, please specifiy a valid IPv4 CIDR, ex: 10.0.0.0/24", flags["cidr"])
		}
		execPath, _ := getExecPath()
		// cwd, err := os.Getwd()

		workingDir := "/workspaces/avx-single-region-aws/terraform/modules/controller"

		// Setup terraform environment and check for errors
		tf, err := tfexec.NewTerraform(workingDir, execPath)
		if err != nil {
			color.Magenta(dt+":", err)
		}
		// Create Terraform workspace and check for errors
		err = tf.WorkspaceNew(ctx, "avxctl-ctrl")
		if strings.Contains(err.Error(), "exists") {
			color.Magenta(err.Error())
			tf.WorkspaceSelect(ctx, "avxctl-ctrl")
		}
		// Initialize Terraform workspace and check for errors
		err = tf.Init(ctx, tfexec.Upgrade(true))
		if err != nil {
			color.Magenta(dt+":", err)
		}
		// Set up options for Terraform execution
		lockOption := tfexec.Lock(false)
		parallelism := tfexec.Parallelism(10)
		var1 := tfexec.Var(fmt.Sprintf("vpc_cidr=%s", cidr))
		// var2 := tfexec.Var(fmt.Sprintf("ctrl_name=%s", name))
		// var3 := tfexec.Var(fmt.Sprintf("ctrl_version=%s", version))
		// applyOptions := []tfexec.ApplyOption{
		// 	tfexec.Var("vpc_cidr=%s", cidr),
		// 	tfexec.Var("ctrl_name=%s", name),
		// 	tfexec.Var("ctrl_version=%s", version),
		// }
		// vars := []tfexec.VarOption{
		// 	{
		// 		Name:  "vpc_cidr",
		// 		Value: cidr,
		// 	},
		// 	{
		// 		Name:  "ctrl_version",
		// 		Value: version,
		// 	},
		// 	{
		// 		Name:  "ctrl_name",
		// 		Value: name,
		// 	},
		// }
		color.Green(dt + ": deploying aviatrix controller")

		// Apply terraform code and return the error
		err = tf.Apply(ctx, var1, lockOption, parallelism)
		if err != nil {
			color.Magenta(dt+":", err)
		}
		// Get terraform output
		output, err := tf.Output(ctx)
		if err != nil {
			color.Magenta(dt+":", err)
		}

		ctrlIp, _ := output["avx_ctrl_public_ip"].Value.MarshalJSON()
		cpltIp, _ := output["copilot_public_ip"].Value.MarshalJSON()
		color.Yellow(fmt.Sprintf("\nController PublicIP: %s\n", string(ctrlIp)))
		color.Yellow(fmt.Sprintf("\nController PublicIP: %s\n", string(cpltIp)))
	},
}

func init() {
	CreateCmd.AddCommand(controllerCmd)

	// Define local flags for the controller sub-command
	controllerCmd.Flags().StringP("name", "n", "avxctl-controller", "Name of the Aviatrix Controller")
	controllerCmd.Flags().StringP("ver", "v", "latest", "Aviatrix controller version to deploy")
	controllerCmd.Flags().StringP("cidr", "c", "10.0.0.0/24", "CIDR Block for deploying VPC/VNET for Aviatrix Controller")
	controllerCmd.Flags().StringP("email", "e", "", "Admininstrator email for Aviatrix Controller")
	controllerCmd.Flags().StringP("password", "p", "", "Admininistrator password for Aviatrix Controller")
	controllerCmd.Flags().StringP("customer-id", "l", "", "Aviatrix CustomerID to use")
	// controllerCmd.MarkFlagRequired("email")
	// controllerCmd.MarkFlagRequired("password")
	// controllerCmd.MarkFlagRequired("customer-id")
}

// Define a function that creates an Aviatrix Controller
func isEmpty(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}

func isIPv4CIDR(s string) bool {
	fmt.Println("1")
	_, ipv4Net, err := net.ParseCIDR(s)
	if err != nil {
		return false
	}
	return ipv4Net.IP.To4() != nil
}
func isFlagEmpty(m map[string]string) {
	for key, value := range m {
		if isEmpty(value) {
			color.Cyan("%s cannot be empty", key)
		}
	}
}
func getExecPath() (string, error) {
	// Ensure the correct version of Terraform
	execPath, err := exec.LookPath("terraform")
	return execPath, err
}
