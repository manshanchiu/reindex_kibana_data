package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	indexes := "indexname"
	index_arr := strings.Split(indexes, ",")
	for _, v := range index_arr {
		setNewMapping(v + "temp")
		reIndex(v, v+"temp")
		deleteIndex(v)
		setNewMapping(v)
		reIndex(v+"temp", v)
		deleteIndex(v + "temp")
	}
}

func setNewMapping(index_name string) {
	jsonStr := []byte(`{
		"mappings": {
			"properties": {
			  "sla_breached": {
				"type": "text",
				"fields": {
				  "keyword": {
					"type": "keyword",
					"ignore_above": 256
				  }
				}
			  }
			}
		}
	  }`)
	req, _ := http.NewRequest("PUT", os.Getenv("KIBANA_ENDPOINT")+index_name, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(os.Getenv("KIBANA_USERNAME"), "KIBANA_PASSWORD")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(bodyBytes))
}

func reIndex(old_name string, new_name string) {
	s := `{
		"source": {
		  "index": "{{old_name}}"
		},
		"dest": {
		  "index": "{{new_name}}"
		}
	  }`
	s = strings.Replace(s, "{{old_name}}", old_name, 1)
	s = strings.Replace(s, "{{new_name}}", new_name, 1)
	jsonStr := []byte(s)
	req, _ := http.NewRequest("POST", os.Getenv("KIBANA_ENDPOINT")+"_reindex?refresh=true", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(os.Getenv("KIBANA_USERNAME"), "KIBANA_PASSWORD")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(bodyBytes))
}

func deleteIndex(index_name string) {
	req, _ := http.NewRequest("DELETE", os.Getenv("KIBANA_ENDPOINT")+index_name, nil)
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(os.Getenv("KIBANA_USERNAME"), "KIBANA_PASSWORD")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(bodyBytes))
}
