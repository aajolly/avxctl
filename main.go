/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package main

import "github.com/aajolly/avxctl/cmd"

func main() {
	cmd.Execute()
}

// PLAYGROUND
// package main

// import (
// 	"crypto/tls"
// 	"encoding/json"
// 	"fmt"
// 	"io"
// 	"net/http"
// 	"net/url"
// 	"strings"

// 	"github.com/fatih/color"
// )

// func main() {
// 	transport := &http.Transport{
// 		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
// 	}

// 	// Create a new http client with the custom transport
// 	client := &http.Client{Transport: transport}

// 	cid, _ := GetCid(client, "54.206.204.81", "10.10.0.4")
// 	fmt.Println(cid)
// }

// func GetCid(client *http.Client, ctrl_pub_ip string, ctrl_private_ip string) (string, error) {

// 	color.Blue("## Fetching Customer ID...")
// 	var respBody map[string]interface{}
// 	ctrlurl := "https://" + ctrl_pub_ip + "/v1/api"

// 	// Prepare form data
// 	data := url.Values{
// 		"action":   {"login"},
// 		"username": {"admin"},
// 		"password": {ctrl_private_ip},
// 	}
// 	// data.Set("action", "login")
// 	// data.Set("username", "admin")
// 	// data.Set("password", ctrl_private_ip)

// 	// Create a new POST request with the form data
// 	resp, err := client.Post(ctrlurl, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
// 	// req, err := http.NewRequest("POST", ctrlurl, bytes.NewBufferString(data.Encode()))
// 	// if err != nil {
// 	// 	return string(err.Error()), err
// 	// }
// 	// // reqDump, _ := httputil.DumpRequestOut(req, true)
// 	// // fmt.Printf("REQUEST:\n%s", string(reqDump)+"\n")

// 	// req.PostForm = data
// 	// req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

// 	// Make the POST request
// 	// resp, err := client.Do(req)
// 	if err != nil {
// 		return "GET CID Failed", err
// 	}
// 	defer resp.Body.Close()

// 	// Print the response body
// 	body, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		return string(err.Error()), err
// 	}
// 	json.Unmarshal(body, &respBody)
// 	cid, ok := respBody["CID"]
// 	if !ok {
// 		fmt.Printf("CID does not exist yet\n")
// 		return "CID does not exist yet", nil
// 	}
// 	cid1, _ := cid.(string)
// 	color.Blue("## Fetched Customer ID...")
// 	return cid1, nil
// }
