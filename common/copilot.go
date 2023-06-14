package common

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/fatih/color"
)

// CoPilot Deployment functions
func DeployCopilot(
	client *http.Client,
	ctrlurl string,
	cid string,
	account_name string,
	region string,
	vpc_id string,
	subnet_cidr string,
	pwd string,
	ver string,
) {
	color.Blue("## Deploying Aviatrix CoPilot...")
	var respBody map[string]interface{}

	// Prepare Data
	data := url.Values{}
	data.Set("action", "deploy_copilot")
	data.Set("CID", cid)
	data.Set("cloud_type", "1")
	data.Set("account_name", account_name)
	data.Set("region_name", region)
	data.Set("controller_service_account_username", "admin")
	data.Set("controller_service_account_password", pwd)
	data.Set("vpc_id", vpc_id)
	data.Set("subnet", subnet_cidr)
	data.Set("instance_size", "t3.2xlarge")
	data.Set("data_volume_size", "50")
	data.Set("is_cluster", "false")

	// Create a new POST request with the form data
	resp, err := client.Post(ctrlurl, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))

	if err != nil {
		color.Red(err.Error())
		panic(err)
	}
	defer resp.Body.Close()

	// Print the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		color.Red(err.Error())
	}
	json.Unmarshal(body, &respBody)
	fmt.Println(respBody)
	if respBody["return"] == true {
		color.Green(respBody["results"].(map[string]interface{})["text"].(string))
	} else {
		color.Red(respBody["reason"].(string))
	}
}
func DestroyCopilot(client *http.Client, cid string, ctrlurl string) {
	color.Blue("## Destroy CoPilot")
	var respBody map[string]interface{}

	// Prepare Data
	data := url.Values{
		"action":     {"cleanup_copilot"},
		"CID":        {cid},
		"is_cluster": {"false"},
	}

	// Create a new POST request with the form data
	resp, err := client.Post(ctrlurl, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))

	if err != nil {
		color.Red(err.Error())
		panic(err)
	}
	defer resp.Body.Close()

	// Print the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		color.Red(err.Error())
	}
	json.Unmarshal(body, &respBody)
	fmt.Println(respBody)
	if respBody["return"] == true {
		color.Green(respBody["results"].(map[string]interface{})["result"].(string))
	} else {
		color.Red(respBody["reason"].(string))
	}

}
