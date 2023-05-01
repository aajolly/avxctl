package common

import (
	"log"
	"net"
	"os/exec"
	"strings"

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
			color.Cyan("%s cannot be empty", key)
		}
	}
}
func GetExecPath() (string, error) {
	// Ensure the correct version of Terraform
	execPath, err := exec.LookPath("terraform")
	return execPath, err
}
