package common

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/fatih/color"
)

func GetCid(ctrl_pub_ip string, ctrl_private_ip string) (string, error) {

	color.Blue("## Fetching Customer ID...")
	var respBody map[string]interface{}
	ctrlurl := "https://" + ctrl_pub_ip + "/v1/api"

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	// Create a new http client with the custom transport
	client := &http.Client{Transport: transport}

	// Prepare form data
	data := url.Values{}
	data.Set("action", "login")
	data.Set("username", "admin")
	data.Set("password", ctrl_private_ip)

	// Create a new POST request with the form data
	req, err := http.NewRequest("POST", ctrlurl, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return string(err.Error()), err
	}
	// reqDump, _ := httputil.DumpRequestOut(req, true)
	// fmt.Printf("REQUEST:\n%s", string(reqDump)+"\n")

	req.PostForm = data
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Make the POST request
	resp, err := client.Do(req)
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
		return "CID does not exist yet", nil
	}
	cid1, _ := cid.(string)
	color.Blue("## Fetched Customer ID...")
	return cid1, nil
}
