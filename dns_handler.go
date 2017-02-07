package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"net/http"
	"net/url"
	"strings"
	"github.com/bitly/go-simplejson"
)

func getCurrentIP(url string) (string, error) {
	response, err := http.Get(url)

	if err != nil {
		log.Println(url)
		return "", err
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	// is ip address
	return string(body), nil
}

func generateBody(content url.Values) url.Values {
	body := url.Values{}
	body.Add("login_token", config.LoginToken)
	body.Add("format", "json")
	body.Add("lang", "en")
	body.Add("error_on_empty", "no")

	if content != nil {
		for k := range content {
			body.Add(k, content.Get(k))
		}
	}

	return body
}

func getSubDomain(domain string, name string, dnsType string) (string, string) {
	log.Println("debug:", domain, name)
	var ret, ip string
	value := url.Values{}
	value.Add("domain", domain)
	value.Add("sub_domain", name)

	response, err := postData("/Record.List", value)

	if err != nil {
		log.Println("Failed to get domain list")
		return "", ""
	}

	sjson, parseErr := simplejson.NewJson([]byte(response))

	if parseErr != nil {
		log.Println(parseErr)
		return "", ""
	}

	if sjson.Get("status").Get("code").MustString() == "1" {
		records, _ := sjson.Get("records").Array()

		if len(records) == 0 {
			log.Println("records slice is empty.")
		}

		for _, d := range records {
			m := d.(map[string]interface{})
			if m["name"] == name && m["type"] == dnsType {
				ret = m["id"].(string)
				ip = m["value"].(string)
				break
			}
		}

	} else {
		log.Println("get_subdomain:status code:", sjson.Get("status").Get("code").MustString())
	}

	return ret, ip
}

func updateIP(domain string, subDomainID string, subDomainName string, ip string, ttl int) {
	value := url.Values{}
	value.Add("domain", domain)
	value.Add("record_id", subDomainID)

	value.Add("sub_domain", subDomainName)
	value.Add("record_type", "A")
	value.Add("record_line", "默认")
	value.Add("value", ip)
	value.Add("ttl", strconv.Itoa(ttl))


	response, err := postData("/Record.Modify", value)

	if err != nil {
		log.Println("Failed to update record to new IP!")
		log.Println(err)
		return
	}

	sjson, parseErr := simplejson.NewJson([]byte(response))

	if parseErr != nil {
		log.Println(parseErr)
		return
	}

	if sjson.Get("status").Get("code").MustString() == "1" {
		log.Println("New IP updated!")
	}

}

func postData(url string, content url.Values) (string, error) {
	client := &http.Client{}
	values := generateBody(content)
	req, _ := http.NewRequest("POST", "https://dnsapi.cn"+url, strings.NewReader(values.Encode()))

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", fmt.Sprintf("hmldd DDNS/0.0.1 (%s)", "hmldd@qq.com"))

	response, err := client.Do(req)

	if err != nil {
		log.Println("Post failed...")
		log.Println(err)
		return "", err
	}

	defer response.Body.Close()
	resp, _ := ioutil.ReadAll(response.Body)

	return string(resp), nil
}
