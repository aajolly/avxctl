/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package create

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aajolly/avx-single-region-aws/democtl/common"
	"github.com/fatih/color"
	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Variable to store the file path
var filePath string

// controllerCmd represents the controller command
var controllerCmd = &cobra.Command{
	Use:   "controller",
	Short: "Creates an Aviatrix Controller",
	Long:  `Creates an Aviatrix Controller with specific version, 6.7+`,
	Run: func(cmd *cobra.Command, args []string) {
		// filePath := viper.GetString("file")

		// Extract the directory path
		dirPath := filepath.Dir(filePath)
		color.Blue("Config File Directory: %s", dirPath)

		// Extract the file name
		fileName := filepath.Base(filePath)
		color.Blue("Config File Name: %s", fileName)
		viper.AddConfigPath(dirPath)
		viper.SetConfigName(fileName)
		viper.SetConfigType("yaml")
		err := viper.ReadInConfig()
		if err != nil {
			log.Fatalf("fatal error config file: %v", err)
		}
		config := common.Config{}
		err = viper.Unmarshal(&config)
		if err != nil {
			log.Fatalf("unable to decode into struct: %v", err)
		}
		color.Green(viper.GetString("demo.controller.region"))

		// Parse attributes from config file
		name := viper.GetString("demo.controller.name")
		viper.SetDefault("demo.controller.name", "avxctl-ctlr") // Set default value for Controller name
		region := viper.GetString("demo.controller.region")
		version := viper.GetString("demo.controller.version")
		viper.SetDefault("demo.controller.version", "latest") // Set default value for Controller version
		cidr := viper.GetString("demo.controller.vpcCidr")
		viper.SetDefault("demo.controller.vpcCidr", "10.0.0.0/24") // Set default value for Controller VPC
		email := viper.GetString("demo.controller.email")
		pwd := viper.GetString("demo.controller.password")
		lic := viper.GetString("demo.controller.customerId")
		key := viper.GetString("demo.controller.keypair")
		color.Blue("## Creating Aviatrix Controller version: %s", version)
		dt := time.Now().Format("01-02-2006 15:04:05")
		ctx := context.Background()

		// flags := map[string]string{
		// 	"name":    name,
		// 	"region":  region,
		// 	"version": version,
		// 	"cidr":    cidr,
		// 	"email":   email,
		// 	"pwd":     pwd,
		// 	"lic":     lic,
		// 	"key":     key,
		// }

		// common.IsFlagEmpty(flags)
		// if !common.IsIPv4CIDR(flags["cidr"]) {
		// 	color.Magenta("%s is not a valid IPv4 CIDR, please specifiy a valid IPv4 CIDR, ex: 10.0.0.0/24", flags["cidr"])
		// }
		execPath, _ := common.GetExecPath()
		cwd, _ := os.Getwd()

		workingDir := fmt.Sprintf("%s/terraform/modules/controller", cwd)
		trimmedVersion := common.TrimVersion(version)

		// Setup terraform environment and check for errors
		tf, err := tfexec.NewTerraform(workingDir, execPath)
		if err != nil {
			color.Magenta(dt+":", err.Error())
			panic(err)
		}
		// Create Terraform workspace and check for errors
		err = tf.WorkspaceNew(ctx, "avxctl-ctrl")
		if strings.Contains(err.Error(), "exists") {
			color.Magenta(err.Error())
			color.Blue("## Selecting existing terraform workspace...")
			tf.WorkspaceSelect(ctx, "avxctl-ctrl")
		}
		// Initialize Terraform workspace and check for errors
		err = tf.Init(ctx, tfexec.Upgrade(true))
		if err != nil {
			color.Magenta(dt+":", err.Error())
			panic(err)
		}
		// Set up options for Terraform execution
		lockOption := tfexec.Lock(false)
		parallelism := tfexec.Parallelism(10)
		var1 := tfexec.Var(fmt.Sprintf("vpc_cidr=%s", cidr))
		var2 := tfexec.Var(fmt.Sprintf("region=%s", region))
		var3 := tfexec.Var(fmt.Sprintf("ctrl_version=%s", trimmedVersion))
		var4 := tfexec.Var(fmt.Sprintf("keypair=%s", key))

		color.Green(dt + ": deploying aviatrix controller")

		// Apply terraform code and return the error
		err = tf.Apply(ctx, var1, var2, var3, var4, lockOption, parallelism)
		if err != nil {
			color.Magenta(err.Error())
			panic(err)
		}
		// Get terraform output
		output, err := tf.Output(ctx)
		if err != nil {
			color.Magenta(err.Error())
			panic(err)
		}

		ctrlPubIp, _ := output["avx_ctrl_public_ip"].Value.MarshalJSON()
		ctrlPrivIp, _ := output["avx_ctrl_private_ip"].Value.MarshalJSON()
		// cpltIp, _ := output["copilot_public_ip"].Value.MarshalJSON()
		aws_role_arn, _ := output["aws_role_arn"].Value.MarshalJSON()
		aws_role_ec2, _ := output["aws_role_ec2"].Value.MarshalJSON()
		ctrl_vpc_id, _ := output["avxctl_ctrl_vpc_id"].Value.MarshalJSON()
		ctrl_subnet_cidr, _ := output["avxctl_ctrl_subnet_cidr"].Value.MarshalJSON()
		color.Yellow(fmt.Sprintf("\nAviatrix Controller PublicIP: %s\n", strings.Trim(string(ctrlPubIp), "\"")))
		// color.Yellow(fmt.Sprintf("\nAviatrix CoPilot PublicIP: %s\n", strings.Trim(string(cpltIp), "\"")))

		transport := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}

		// Create a new http client with the custom transport
		client := &http.Client{Transport: transport}

		// Initialize Aviatrix Controller
		common.CtrlInitialize(
			client,
			strings.Trim(string(ctrlPubIp), "\""),
			strings.Trim(string(ctrlPrivIp), "\""),
			// strings.Trim(string(cpltIp), "\""),
			lic,
			version,
			pwd,
			email,
			region,
			strings.Trim(string(aws_role_arn), "\""),
			strings.Trim(string(aws_role_ec2), "\""),
			strings.Trim(string(ctrl_vpc_id), "\""),
			strings.Trim(string(ctrl_subnet_cidr), "\""),
			name,
			"create",
		)
	},
}

func init() {
	CreateCmd.AddCommand(controllerCmd)

	// Define local flags for the controller sub-command
	// controllerCmd.Flags().StringP("name", "n", "avxctl-controller", "Name of the Aviatrix Controller")
	// controllerCmd.Flags().StringP("ver", "v", "latest", "Aviatrix controller version to deploy")
	// controllerCmd.Flags().StringP("cidr", "c", "10.0.0.0/24", "CIDR Block for deploying Aviatrix Controller")
	// controllerCmd.Flags().StringP("email", "e", "", "Admininstrator email for Aviatrix Controller")
	// controllerCmd.Flags().StringP("password", "p", "", "Admininistrator password for Aviatrix Controller")
	// controllerCmd.Flags().StringP("customer-id", "l", "", "Aviatrix CustomerID to use")
	// controllerCmd.Flags().StringP("region", "r", "", "CSP Region for deploying Aviatrix Controller")
	// controllerCmd.Flags().StringP("keypair", "k", "", "KeyPair to use for deploying Aviatrix Controller")
	// controllerCmd.Flags().String("aws-account-id", "", "AWS Account ID for Aviatrix Controller Deployment")
	// controllerCmd.Flags().String("azure-subs-id", "", "Azure Subscription ID for Aviatrix Controller Deployment")
	controllerCmd.Flags().StringVarP(&filePath, "file", "f", "", "Path to the config file")
	viper.BindPFlag("file", controllerCmd.Flags().Lookup("file"))
	// controllerCmd.MarkFlagRequired("email")
	// controllerCmd.MarkFlagRequired("password")
	// controllerCmd.MarkFlagRequired("customer-id")
	// controllerCmd.MarkFlagRequired("region")
	// controllerCmd.MarkFlagRequired("keypair")
	controllerCmd.MarkFlagRequired("file")
}
