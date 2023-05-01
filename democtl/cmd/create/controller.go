/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package create

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/aajolly/avx-single-region-aws/democtl/common"
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

		common.IsFlagEmpty(flags)
		color.Green(dt + ": Validating CIDR")
		if !common.IsIPv4CIDR(flags["cidr"]) {
			color.Magenta("%s is not a valid IPv4 CIDR, please specifiy a valid IPv4 CIDR, ex: 10.0.0.0/24", flags["cidr"])
		}
		execPath, _ := common.GetExecPath()
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

		ctrlPubIp, _ := output["avx_ctrl_public_ip"].Value.MarshalJSON()
		ctrlPrivIp, _ := output["avx_ctrl_private_ip"].Value.MarshalJSON()
		cpltIp, _ := output["copilot_public_ip"].Value.MarshalJSON()
		aws_role_arn, _ := output["aws_role_arn"].Value.MarshalJSON()
		aws_role_ec2, _ := output["aws_role_ec2"].Value.MarshalJSON()
		color.Yellow(fmt.Sprintf("\nAviatrix Controller PublicIP: %s\n", string(ctrlPubIp)))
		color.Yellow(fmt.Sprintf("\nAviatrix CoPilot PublicIP: %s\n", string(cpltIp)))

		transport := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}

		// Create a new http client with the custom transport
		client := &http.Client{Transport: transport}

		// cid, _ := common.CtrlInitialize(client, string(ctrlPubIp), string(ctrlPrivIp))
		common.CtrlInitialize(client, string(ctrlPubIp), string(ctrlPrivIp), string(cpltIp), lic, version, pwd, email, string(aws_role_arn), string(aws_role_ec2))
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
	controllerCmd.Flags().StringP("aws-account-id", "aaid", "", "AWS Account ID for Aviatrix Controller Deployment")
	controllerCmd.Flags().StringP("azure-subs-id", "asub", "", "Azure Subscription ID for Aviatrix Controller Deployment")
	// controllerCmd.MarkFlagRequired("email")
	// controllerCmd.MarkFlagRequired("password")
	// controllerCmd.MarkFlagRequired("customer-id")
}
