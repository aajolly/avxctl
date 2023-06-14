/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package delete

import (
	"context"
	"crypto/tls"
	"log"
	"net/http"
	"strings"
	"time"

	"path/filepath"

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
	Short: "Delete an Aviatrix Controller",
	Long:  `Delete an Aviatrix Controller`,
	Run: func(cmd *cobra.Command, args []string) {
		// cwd, err := os.Getwd()

		filePath := viper.GetString("file")

		// Extrat the directory path
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
		// Parse attributes from config file
		name := viper.GetString("demo.controller.name")
		viper.SetDefault("demo.controller.name", "avxctl-ctlr") // Set default value for Controller name
		region := viper.GetString("demo.controller.region")
		version := viper.GetString("demo.controller.version")
		viper.SetDefault("demo.controller.version", "latest") // Set default value for Controller version
		email := viper.GetString("demo.controller.email")
		pwd := viper.GetString("demo.controller.password")
		lic := viper.GetString("demo.controller.customerId")
		color.Blue("## Creating Aviatrix Controller version: %s", version)
		dt := time.Now().Format("01-02-2006 15:04:05")

		ctx := context.Background()
		execPath, _ := common.GetExecPath()

		workingDir := "/workspaces/avx-single-region-aws/terraform/modules/controller"

		// Setup terraform environment and check for errors
		tf, err := tfexec.NewTerraform(workingDir, execPath)
		if err != nil {
			color.Magenta(err.Error())
		}
		// Create Terraform workspace and check for errors
		err = tf.WorkspaceNew(ctx, "avxctl-ctrl")
		if strings.Contains(err.Error(), "exists") {
			color.Magenta(err.Error())
			color.Blue("## Selecting existing workspace...")
			tf.WorkspaceSelect(ctx, "avxctl-ctrl")
		}

		// Initialize Terraform workspace and check for errors
		err = tf.Init(ctx, tfexec.Upgrade(true))
		if err != nil {
			color.Magenta(err.Error())
		}
		// Set up options for Terraform execution
		lockOption := tfexec.Lock(false)
		parallelism := tfexec.Parallelism(10)

		color.Blue("## Deleting Aviatrix Controller...")

		// Get terraform output
		output, err := tf.Output(ctx)
		if err != nil {
			color.Magenta(err.Error())
			panic(err)
		}

		ctrlPubIp, _ := output["avx_ctrl_public_ip"].Value.MarshalJSON()

		// Create a new http client with the custom transport
		transport := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}

		client := &http.Client{Transport: transport}

		common.CtrlInitialize(
			client,
			strings.Trim(string(ctrlPubIp), "\""),
			"",
			lic,
			version,
			pwd,
			email,
			region,
			"",
			"",
			"",
			"",
			name,
			"delete",
		)

		// Execute terraform destroy and return the error
		err = tf.Destroy(ctx, lockOption, parallelism)
		if strings.Contains(err.Error(), "DependencyViolation") {
			color.Magenta(err.Error())
			err := common.DeleteSecurityGroupsByName("AviatrixSecurityGroup")
			if err != nil {
				color.Magenta(dt+":", err.Error())
			}
		}
		color.Blue("## Attempting deletion again...")

		err = tf.Destroy(ctx, lockOption, parallelism)
		if err != nil {
			color.Magenta(dt+":", err.Error())
		}
	},
}

func init() {
	DeleteCmd.AddCommand(controllerCmd)
	controllerCmd.Flags().StringVarP(&filePath, "file", "f", "", "Path to the config file")
	viper.BindPFlag("file", controllerCmd.Flags().Lookup("file"))
	controllerCmd.MarkFlagRequired("file")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// controllerCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// controllerCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
