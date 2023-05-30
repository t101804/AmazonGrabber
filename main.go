// Copyright (c) 2023 By @CallMeRep
// Buy apikey on telegram directly with me @CallMeRep
//
//	add update auto customization
package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gojek/heimdall/v7/httpclient"
	"github.com/manifoldco/promptui"
	"github.com/valyala/fastjson"
)

func clearScreen() {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	} else {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

type myHTTPClient struct {
	client http.Client
}

func (c *myHTTPClient) Do(request *http.Request) (*http.Response, error) {
	return c.client.Do(request)
}

func parser(data string) error {
	parser := fastjson.Parser{}
	lines := strings.Split(data, "\n")
	ipFile, err := os.OpenFile("ip.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("Error creating ip.txt:", err)
		return err
	}
	defer ipFile.Close()
	ptrFile, err := os.OpenFile("sites.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("Error creating ptr.txt:", err)
		return err
	}
	defer ptrFile.Close()
	for _, line := range lines {
		if len(line) == 0 {
			// Skip empty lines
			continue
		}
		// Parse the line as a JSON object
		line = strings.TrimSuffix(line, ",")
		value, err3 := parser.Parse(line)
		if err3 != nil {
			fmt.Println("Error parsing JSON:", err, line)
			continue
		}
		// Extract the "ip" and "ptr" fields from the JSON object
		ip := string(value.GetStringBytes("ip"))
		ptrArray := value.GetArray("ptr")
		var ptrs []string
		for i := 0; i < len(ptrArray); i++ {
			ptrs = append(ptrs, string(ptrArray[i].GetStringBytes()))
		}
		// Write the extracted data to output files
		_, err = fmt.Fprintln(ipFile, ip)
		if err != nil {
			fmt.Println("Error writing to ip.txt:", err)
			continue
		}
		for _, ptr := range ptrs {
			_, err = fmt.Fprintln(ptrFile, ptr)
			if err != nil {
				fmt.Println("Error writing to ptr.txt:", err)
				continue
			}
		}
	}
	return err

}

func apikey() string {
	file, err := os.Open("config.conf")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "your_apikey_token=") {
			return strings.Split(line, "=")[1]
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return ""
}

func api_server() string {
	file, err := os.Open("config.conf")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "server_address=") {
			return strings.Split(line, "=")[1]
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return ""
}
func main() {
	// Make 10 requests to Google.com
	clearScreen()
	numRequests := 10
	// Define the prompt questions
	prompt := promptui.Prompt{
		Label:   "How much results that you want per loopings ?",
		Default: "10000",
	}
	results, err := prompt.Run()
	if err != nil {
		fmt.Println(err)
		return
	}
	clearScreen()
	prompt_loops := promptui.Prompt{
		Label: "How much loopings ?",
	}
	loop, err := prompt_loops.Run()
	if err != nil {
		fmt.Println(err)
		return
	}
	clearScreen()
	numRequests, err = strconv.Atoi(loop)
	if err != nil {
		fmt.Println(err)
	}
	validateApiKey := strings.ReplaceAll(apikey(), "\n", "")
	// Create a WaitGroup to synchronize the goroutines
	var wg sync.WaitGroup
	wg.Add(numRequests)
	start := time.Now()
	// Start the requests in goroutines
	for i := 0; i < numRequests; i++ {
		// Make the HTTP request to the api
		transport := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
		client := httpclient.NewClient(httpclient.WithHTTPClient(&myHTTPClient{client: http.Client{Transport: transport}}))
		site := fmt.Sprintf("%s%s", api_server(), results)
		req, err := http.NewRequest("GET", site, nil)
		if err != nil {
			log.Fatal(err)
		}
		req.Header.Add("X-Api-Key", validateApiKey)
		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("Error fetching api : %s\n", err.Error())
		} else {
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				fmt.Println(err)
			}
			go func() {
				if !strings.Contains(string(body), "Authentication") {
					fmt.Printf("Fetched %s%s\n", api_server(), results)
					err := parser(string(body))
					if err != nil {
						log.Fatal(err)
					}
					resp.Body.Close()
				} else {
					fmt.Println("Buy VIP To Use Mass Scrapping of AWS ")
				}
				wg.Done()
			}()

		}

		// Notify the WaitGroup that the request is complete

	}

	// Wait for all requests to complete
	wg.Wait()
	elapsed := time.Since(start)
	fmt.Printf("Elapsed time: %s\n", elapsed)
}
