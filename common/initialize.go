package common

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/fatih/color"
)

func getCid(client *http.Client, password string, ctrlurl string) (string, error) {

	color.Blue("## Fetching CID...")
	var respBody map[string]interface{}

	// Check whether Controller is UP
	CheckServerUp(client, ctrlurl, "Aviatrix Controller")

	// Prepare form data
	data := url.Values{
		"action":   {"login"},
		"username": {"admin"},
		"password": {password},
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
		color.Red(err.Error())
		panic(err)
	}
	json.Unmarshal(body, &respBody)

	var cid1 string

	cid, ok := respBody["CID"]
	fmt.Println(ok)
	if !ok {
		color.Red("CID does not exist yet...waiting for 60 seconds")
		time.Sleep(60 * time.Second)
		resp, err := client.Post(ctrlurl, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
		if err != nil {
			return "GET CID Failed", err
		}
		defer resp.Body.Close()
		json.Unmarshal(body, &respBody)
		fmt.Println(respBody)
		cid = respBody["CID"]
		cid1, _ = cid.(string)
		color.Green("fetched CID...")

		return cid1, nil

	} else {
		cid1, _ = cid.(string)
		color.Green("fetched CID...")

		return cid1, nil
	}

}

func setAdminEmail(client *http.Client, cid string, email string, ctrlurl string) (string, error) {
	color.Blue("## Setting Admin Email...")
	var respBody map[string]interface{}

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
	fmt.Println(respBody)
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
		color.Green("Initial Setup COMPLETE - Aviatrix Controller...")
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
		"account_name": {"admin"},
		"username":     {"admin"},
		"password":     {ctrl_private_ip},
		"what":         {"password"},
		"email":        {email},
		"old_password": {ctrl_private_ip},
		"new_password": {pwd},
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
	fmt.Println(respBody)
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
	fmt.Println(respBody)
	if respBody["return"] == true {
		color.Green("Aviatrix License Set")
	} else {
		color.Red("Failed to set Aviatrix Controller License")
		panic(errors.New("failed to set aviatrix controller license"))
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

	color.Yellow(aws_role_arn)
	color.Yellow(aws_role_ec2)
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
		"groups":             {"admin"},
		"skip_sg_config":     {"true"},
	}
	// Create a new POST request with the form data
	resp, err := client.Post(ctrlurl, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))

	if err != nil {
		color.Red("1")
		color.Red(err.Error())
	}
	defer resp.Body.Close()

	// Print the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		color.Red("2")
		color.Red(err.Error())
	}
	json.Unmarshal(body, &respBody)
	fmt.Println(respBody)
	if respBody["return"] == true {
		color.Green(respBody["results"].(string))
	} else {
		color.Red(respBody["results"].(string))
	}

}

func onboardAzure(client *http.Client,
	ctrlurl string,
	cid string,
	sid string,
	appcid string,
	appsecret string,
	appendp string,
	email string) {

	color.Blue("## Onboarding Azure Subscription...")
	var respBody map[string]interface{}

	// Define data
	data := url.Values{
		"action":                        {"setup_account_profile"},
		"CID":                           {cid},
		"account_name":                  {"azure-subscription1"},
		"cloud_type":                    {"8"},
		"account_email":                 {email},
		"arm_subscription_id":           {sid},
		"arm_application_endpoint":      {appendp},
		"arm_application_client_id":     {appcid},
		"arm_application_client_secret": {appsecret},
		"groups":                        {"admin"},
		"skip_sg_config":                {"true"},
	}
	// Create a new POST request with the form data
	resp, err := client.Post(ctrlurl, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))

	if err != nil {
		color.Red("1")
		color.Red(err.Error())
	}
	defer resp.Body.Close()

	// Print the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		color.Red("2")
		color.Red(err.Error())
	}
	json.Unmarshal(body, &respBody)
	fmt.Println(respBody)
	if respBody["return"] == true {
		color.Green(respBody["results"].(string))
	} else {
		color.Red(respBody["results"].(string))
	}

}

func enableSGMgmt(client *http.Client, ctrlurl string, cid string) {
	color.Blue("## Enabling Security Group Management...")
	var respBody map[string]interface{}

	// Define parameters
	params := url.Values{
		"action":              {"enable_controller_security_group_management"},
		"CID":                 {cid},
		"access_account_name": {"aws-account1"},
	}
	// Create a new POST request with the form data
	resp, err := client.Post(ctrlurl, "application/x-www-form-urlencoded", strings.NewReader(params.Encode()))

	if err != nil {
		color.Red("1")
		color.Red(err.Error())
	}
	defer resp.Body.Close()

	// Print the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		color.Red("2")
		color.Red(err.Error())
	}
	json.Unmarshal(body, &respBody)
	fmt.Println(respBody)
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
		"action": {"get_ctrl_aws_acct_num"},
		"CID":    {cid},
	}
	// Define URL with parameters
	queryUrl := ctrlurl + "?" + params.Encode()

	// Make the GET request
	resp, err := client.Get(queryUrl)
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
	// color.Blue(respBody["results"].(string))
	if respBody["return"] == true {
		color.Green("fetched AWS Account Number...")
		return respBody["results"].(string), nil
	} else {
		color.Red("Failed to fetch AWS Account Number...")
		return respBody["results"].(string), errors.New(respBody["results"].(string))
	}
}

func setCtrlName(client *http.Client, ctrlurl string, ctrl_name string, cid string, version string) {
	color.Blue("## Setting Aviatrix Controller Name: %s", ctrl_name)
	var respBody map[string]interface{}

	// Prepare data
	data := url.Values{
		"action":          {"set_controller_name"},
		"CID":             {cid},
		"controller_name": {ctrl_name + "-" + version},
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
	fmt.Println(respBody)
	if respBody["return"] == true {
		color.Green("Controller name set to: %s", ctrl_name)
	} else {
		color.Red("Failed to set controller name: %s", ctrl_name)
	}
}

func CtrlInitialize(client *http.Client,
	ctrl_pub_ip string,
	ctrl_private_ip string,
	// cplt_public_ip string,
	cust_id string,
	ver string,
	password string,
	email string,
	region string,
	aws_role_arn string,
	aws_role_ec2 string,
	ctrl_vpc_id string,
	ctrl_subnet_cidr string,
	ctrl_name string,
	sid string,
	appcid string,
	appsecret string,
	appendp string,
	cmd string) {

	ctrlurl := "https://" + ctrl_pub_ip + "/v1/api"
	switch cmd {
	case "create":
		color.Blue("## Creating Resources...")
		cid, _ := getCid(client, ctrl_private_ip, ctrlurl)
		setCustomerId(client, cid, cust_id, ctrlurl)
		setAdminEmail(client, cid, email, ctrlurl)
		runCtrlInitialSetup(client, cid, ver, ctrlurl)
		// Get CID again after initial setup
		cid1, _ := getCid(client, ctrl_private_ip, ctrlurl)
		setCtrlName(client, ctrlurl, ctrl_name, cid1, ver)
		aai, _ := getAWSAccountNumber(client, ctrlurl, cid1)
		onboardAWSAccount(client, ctrlurl, aai, cid1, email, aws_role_arn, aws_role_ec2)
		onboardAzure(client, ctrlurl, cid, sid, appcid, appsecret, appendp, email)
		setAdminPassword(client, ctrl_private_ip, cid1, password, email, ctrlurl)
		cid2, _ := getCid(client, password, ctrlurl)
		DeployCopilot(client, ctrlurl, cid2, "aws-account1", region, ctrl_vpc_id, ctrl_subnet_cidr, password, ver)

	case "delete":
		color.Blue("## Initiating Deletion of resources...")
		cid, _ := getCid(client, password, ctrlurl)
		DestroyCopilot(client, cid, ctrlurl)
	}

}
