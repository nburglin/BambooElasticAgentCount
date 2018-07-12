package main

import (
	"path/filepath"
	"os"
	"log"
	"github.com/spf13/viper"
	"net/http"
	"crypto/tls"
	"time"
	"encoding/json"
	"io/ioutil"
	"fmt"
	"strconv"
)

type Config struct {
	Url string `json:baseurl`
	Username string `json:username`
	Password string `json:password`
	SkipSslVerify bool `json:skipsslverify`
}

//JSON object that is returned for each agent
type AgentInfoJSON struct {
	Id int `json: id`
	Name string `json: name`
	Type string `json: type`
	Active bool `json: active`
	Enabled bool `json: enabled`
	Busy bool `json: busy`
}

//The actual response from Bamboo is an array of agent objects
type Agents []AgentInfoJSON

func readConfigs() Config {
	//Read in Config File
	var appConf Config
	currentDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatalf("Error finding working directory info. \nError: %v", err)
	}

	viper.SetConfigName("settings")
	viper.AddConfigPath(currentDir)
	viper.SetDefault("InsecureSkipVerify", false)
	err = viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading settings file located in %s. \nError: %v", currentDir, err)
	}

	err = viper.Unmarshal(&appConf)
	if err != nil {
		log.Fatalf("Error parsing data from settings file. \nError: %v", err)
	}

	return appConf
}

func createHttpClient() http.Client {
	//Create the http client and set up the request
	client := http.Client {
		Timeout: time.Second * 30,
	}

	return client
}

func formatHttpRequest(appConf Config) *http.Request{
	req, err := http.NewRequest(http.MethodGet, appConf.Url, nil)
	if err != nil {
		log.Fatalf("Error while setting up the GET request \nError: %v", err)
	}

	req.SetBasicAuth(appConf.Username, appConf.Password)

	return req
}

func parseResponse(resp *http.Response) Agents {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body from API. \nError: %v", err)
	}

	//Parse the json body
	agents := Agents{}
	err = json.Unmarshal(body, &agents)
	if err != nil {
		log.Fatalf("Error parsing the json of body. \nError: %v", err)
	}

	return agents
}

func countElasticAgents(agents Agents) string {
	elasticAgentCount := 0
	for _, agent := range agents {
		if agent.Type == "ELASTIC" {
			elasticAgentCount ++
		}
	}

	elasticAgentCountStr := strconv.Itoa(elasticAgentCount)

	return elasticAgentCountStr

}

func main() {

	appConf := readConfigs()

	//Disable security check of the http certificate to work with self signed cert if set in config file
	if appConf.SkipSslVerify {
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	//Create http client and format the request
	client := createHttpClient()
	req := formatHttpRequest(appConf)

	//Perform the GET
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error performing GET call to the url. \nError: %v", err)
	}

	//Pull out the body from the response and parse into json
    agents := parseResponse(resp)

	//Find how many of these agents are elastic
    elasticAgentCountStr := countElasticAgents(agents)
	fmt.Println("<prtg> <result> <channel>Elastic Agent Count</channel> <value>" + elasticAgentCountStr + "</value> </result> </prtg>")

}