package common

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os/exec"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/fatih/color"
)

func isEmpty(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}

func IsIPv4CIDR(s string) bool {
	log.Println("## Validating CIDR...")
	_, ipv4Net, err := net.ParseCIDR(s)
	if err != nil {
		return false
	}
	return ipv4Net.IP.To4() != nil
}
func IsFlagEmpty(m map[string]string) {
	for key, value := range m {
		if isEmpty(value) {
			panic(fmt.Errorf("%s cannot be empty", key))
		}
	}
}
func TrimVersion(version string) string {
	// Trim the minor release version for terraform
	parts := strings.Split(version, ".")
	if len(parts) >= 3 {
		parts = parts[:2]
	}
	return strings.Join(parts, ".")
}

func CheckServerUp(client *http.Client, url string, node string) {
	color.Blue("## Checking " + node + " is UP...")
	// var urlresponse *http.Response = &http.Response{}
	urlresponse, _ := client.Get(url)

	for {
		if urlresponse != nil {
			if urlresponse.StatusCode == 200 {
				color.Green(node + " is UP...")
				break
			}
		}
		urlresponse, _ = client.Get(url)
	}

}
func GetExecPath() (string, error) {
	// Ensure the correct version of Terraform
	execPath, err := exec.LookPath("terraform")
	return execPath, err
}
func DeleteSecurityGroupsByName(groupName string) error {
	// Load the AWS SDK configuration
	color.Yellow("## Deleting AviatrixSecurityGroups...")
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-southeast-2"))
	if err != nil {
		color.Red("1")
		return err
	}

	// Create an EC2 client
	ec2Client := ec2.NewFromConfig(cfg)

	// Describe the security groups with the specified name
	input := &ec2.DescribeSecurityGroupsInput{
		Filters: []types.Filter{
			{
				Name:   aws.String("tag:Name"),
				Values: []string{groupName},
			},
		},
	}
	result, err := ec2Client.DescribeSecurityGroups(context.TODO(), input)
	if err != nil {
		color.Red("2")
		return err
	}

	// Delete each security group with the specified name
	for _, group := range result.SecurityGroups {
		_, err := ec2Client.DeleteSecurityGroup(context.TODO(), &ec2.DeleteSecurityGroupInput{
			GroupId: group.GroupId,
		})
		if err != nil {
			color.Red("3")
			return err
		}
		color.Yellow("Deleted security group %s\n", *group.GroupName)
	}

	return nil
}
