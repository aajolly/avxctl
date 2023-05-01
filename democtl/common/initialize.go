package common

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/fatih/color"
)

func checkServerUp(client *http.Client, url string, node string) {
	color.Blue("## Checking " + node + "is UP...")
	var urlresponse *http.Response = &http.Response{}

	urlresponse, _ = client.Get(url)

	for {
		if urlresponse != nil {
			if urlresponse.StatusCode == 200 {
				color.Blue(node + "is UP...")
				break
			}
		}
		urlresponse, _ = client.Get(url)
	}

}
func getCid(client *http.Client, ctrl_private_ip string, ctrlurl string) (string, error) {

	color.Blue("## Fetching CID...")
	var respBody map[string]interface{}
	// ctrlurl := "https://" + ctrl_pub_ip + "/v1/api"

	// Check whether Controller is UP
	checkServerUp(client, ctrlurl, "Aviatrix Controller")

	// Prepare form data
	data := url.Values{
		"action":   {"login"},
		"username": {"admin"},
		"password": {ctrl_private_ip},
	}

	// Create a new POST request with the form data
	resp, err := client.Post(ctrlurl, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))

	if err != nil {
		return "GET CID Failed", err
	}
	defer resp.Body.Close()

	// Print the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return string(err.Error()), err
	}
	json.Unmarshal(body, &respBody)
	cid, ok := respBody["CID"]
	if !ok {
		fmt.Printf("CID does not exist yet\n")
		return "CID does not exist yet", errors.New("CID does not exist yet")
	}
	cid1, _ := cid.(string)
	color.Green("## Fetched CID...")
	return cid1, nil
}

func setAdminEmail(client *http.Client, cid string, email string, ctrlurl string) (string, error) {
	color.Blue("## Setting Admin Email...")
	var respBody map[string]interface{}
	// ctrlurl := "https://" + ctrl_pub_ip + "/v1/api"

	//Prepare form data
	data := url.Values{
		"action":      {"add_admin_email_addr"},
		"CID":         {cid},
		"admin_email": {email},
	}

	// Create a new POST request with the form data
	resp, err := client.Post(ctrlurl, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))

	if err != nil {
		return "GET CID Failed", err
	}
	defer resp.Body.Close()

	// Print the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return string(err.Error()), err
	}
	json.Unmarshal(body, &respBody)

	if respBody["return"] == true {
		color.Green(respBody["results"].(string))
		return "OK", nil
	} else {
		color.Red(respBody["results"].(string))
		return "SetAdminEmail Failed", errors.New("SetAdminEmail Failed")
	}

}

func runCtrlInitialSetup(client *http.Client, cid string, ver string, ctrlurl string) (string, error) {
	color.Blue("## Running Initial Setup - Aviatrix Controller...")
	var respBody map[string]interface{}
	// ctrlurl := "https://" + ctrl_pub_ip + "/v1/api"

	//Prepare form data
	data := url.Values{
		"action":         {"initial_setup"},
		"CID":            {cid},
		"subaction":      {"run"},
		"target_version": {ver},
	}

	// Create a new POST request with the form data
	resp, err := client.Post(ctrlurl, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))

	if err != nil {
		return "GET CID Failed", err
	}
	defer resp.Body.Close()

	// Print the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return string(err.Error()), err
	}
	json.Unmarshal(body, &respBody)

	if respBody["return"] == true {
		color.Green("## Initial Setup COMPLETE - Aviatrix Controller...")
		return "OK", nil
	} else {
		color.Red(respBody["results"].(string))
		return "ControllerInitialSetupFailed", errors.New("ControllerInitialSetupFailed")
	}
}

func setAdminPassword(client *http.Client, ctrl_private_ip string, cid string, pwd string, email string, ctrlurl string) (string, error) {
	color.Blue("## Setting Admin Password...")

	var respBody map[string]interface{}
	// ctrlurl := "https://" + ctrl_pub_ip + "/v1/api"

	//Prepare form data
	data := url.Values{
		"action":       {"edit_account_user"},
		"CID":          {cid},
		"account_name": {"run"},
		"username":     {"admin"},
		"password":     {ctrl_private_ip},
		"what":         {"password"},
		"email":        {email},
		"old_password": {ctrl_private_ip},
		"new_password": {pwd},
	}
	// data.Set("action", "edit_account_user")
	// data.Set("CID", cid)
	// data.Set("account_name", "run")
	// data.Set("username", "admin")
	// data.Set("password", ctrl_private_ip)
	// data.Set("what", "password")
	// data.Set("email", email)
	// data.Set("old_password", ctrl_private_ip)
	// data.Set("new_password", pwd)

	// Create a new POST request with the form data
	resp, err := client.Post(ctrlurl, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))

	if err != nil {
		return "GET CID Failed", err
	}
	defer resp.Body.Close()

	// Print the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return string(err.Error()), err
	}
	json.Unmarshal(body, &respBody)

	if respBody["return"] == true {
		color.Green(respBody["results"].(string))
		return "OK", nil
	} else {
		color.Red(respBody["results"].(string))
		return "SetAdminPassword Failed", errors.New("SetAdminPassword Failed")
	}
}

func setCustomerId(client *http.Client, cid string, cust_id string, ctrlurl string) {
	color.Blue("## Setting Customer ID - Aviatrix Controller...")
	var respBody map[string]interface{}

	//Prepare form data
	data := url.Values{
		"action":      {"setup_customer_id"},
		"CID":         {cid},
		"customer_id": {cust_id},
	}

	// Create a new POST request with the form data
	resp, err := client.Post(ctrlurl, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))

	if err != nil {
		color.Red(err.Error())
	}
	defer resp.Body.Close()

	// Print the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		color.Red(err.Error())
	}
	json.Unmarshal(body, &respBody)

	if respBody["return"] == true {
		color.Green(respBody["results"].(string))
	} else {
		color.Red(respBody["results"].(string))
	}
}

func onboardAWSAccount(client *http.Client,
	ctrlurl string,
	aai string,
	cid string,
	email string,
	aws_role_arn string,
	aws_role_ec2 string) {

	color.Blue("## Onboarding AWS Account " + aai)
	var respBody map[string]interface{}

	// Prepare data
	data := url.Values{
		"action":             {"setup_account_profile"},
		"CID":                {cid},
		"account_name":       {"aws-account1"},
		"cloud_type":         {"1"},
		"account_email":      {email},
		"aws_account_number": {aai},
		"aws_iam":            {"true"},
		"aws_role_arn":       {aws_role_arn},
		"aws_role_ec2":       {aws_role_ec2},
		"groups":             {"group1"},
		"skip_sg_config":     {"true"},
	}
	// Create a new POST request with the form data
	resp, err := client.Post(ctrlurl, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))

	if err != nil {
		color.Red(err.Error())
	}
	defer resp.Body.Close()

	// Print the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		color.Red(err.Error())
	}
	json.Unmarshal(body, &respBody)

	if respBody["return"] == true {
		color.Green(respBody["results"].(string))
	} else {
		color.Red(respBody["results"].(string))
	}

}

func getAWSAccountNumber(client *http.Client, ctrlurl string, cid string) (string, error) {
	color.Blue("## Fetching AWS Account Number...")
	var respBody map[string]interface{}

	// Define query parameters
	params := url.Values{
		"action": {},
		"CID":    {cid},
	}

	// Define URL
	queryUrl := ctrlurl + "?" + params.Encode()

	// Make the GET request
	resp, err := client.Get(queryUrl)
	if err != nil {
		color.Red(err.Error())
	}
	defer resp.Body.Close()

	// Print the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		color.Red(err.Error())
	}
	json.Unmarshal(body, &respBody)

	if respBody["return"] == true {
		return respBody["results"].(string), nil
	} else {
		return respBody["results"].(string), errors.New(respBody["results"].(string))
	}
}

func CtrlInitialize(client *http.Client,
	ctrl_pub_ip string,
	ctrl_private_ip string,
	cplt_public_ip string,
	cust_id string,
	ver string,
	password string,
	email string,
	aws_role_arn string,
	aws_role_ec2 string) {

	color.Blue("## Initializing Aviatrix Controller...")
	ctrlurl := "https://" + ctrl_pub_ip + "/v1/api"
	cid, _ := getCid(client, ctrl_private_ip, ctrlurl)
	setCustomerId(client, cid, cust_id, ctrlurl)
	setAdminEmail(client, cid, email, ctrlurl)
	runCtrlInitialSetup(client, cid, ver, ctrlurl)
	aai, _ := getAWSAccountNumber(client, ctrlurl, cid)
	onboardAWSAccount(client, ctrlurl, aai, cid, email, aws_role_arn, aws_role_ec2)
	setAdminPassword(client, ctrl_private_ip, cid, password, email, ctrlurl)
}
