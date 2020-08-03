package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"time"

	"github.com/cloudflare/cloudflare-go"
	"merlin-ddns-cloudflare/config"
)

func init() {
	_ = config.LoadCfg("config.toml")
}

func main() {
	// Construct a new API object
	api, err := cloudflare.New(config.O.Api.ApiToken, config.O.Api.Email)
	if err != nil {
		NotifyDdnsState(false)
		log.Fatal(err)
		return
	}

	zoneID, err := api.ZoneIDByName(config.O.Api.Domain)
	if err != nil {
		NotifyDdnsState(false)
		log.Fatal(err)
		return
	}
	// Fetch all records for a zone
	recs, err := api.DNSRecords(zoneID, cloudflare.DNSRecord{Name: config.O.Api.SubDomain, Type: "A"})
	if err != nil {
		NotifyDdnsState(false)
		log.Fatal(err)
		return
	}
	f := func() (interface{}, error) {
		return getIp(config.O.App.GetIpFromUrl)
	}
	result, err := retry(f, 4)
	if err != nil {
		NotifyDdnsState(false)
		return
	}
	ip, ok := result.(string)
	if ok {
		for _, r := range recs {
			err := api.UpdateDNSRecord(zoneID, r.ID, cloudflare.DNSRecord{Content: ip, Type: "A"})
			if err != nil {
				NotifyDdnsState(false)
				fmt.Println(err.Error())
				return
			}
		}
	}

	NotifyDdnsState(true)
}

func NotifyDdnsState(success bool) {
	if success {
		command := exec.Command("/sbin/ddns_custom_updated", "1")
		err := command.Run()
		if err != nil {
			fmt.Println("update record success.but set flag failed")
			return
		}
		fmt.Println("update record success.")
	} else {
		command := exec.Command("/sbin/ddns_custom_updated", "0")
		err := command.Run()
		if err != nil {
			fmt.Println("update record failed.and set flag failed")
			return
		}
		fmt.Println("update record failed.")
	}
}

func getIp(echoUrl string) (interface{}, error) {
	var client http.Client
	client.Timeout = 100 * time.Second
	get, errGetIp := client.Get(echoUrl)
	if errGetIp != nil {
		return "", errGetIp
	}
	all, errParseIp := ioutil.ReadAll(get.Body)
	defer func() {
		_ = get.Body.Close()
	}()
	if errParseIp != nil {
		return "", errParseIp
	}

	ip := string(all)
	return ip, nil
}

func retry(f func() (interface{}, error), retryTimes int) (interface{}, error) {
	if result, err := f(); err != nil {
		if retryTimes > 0 {
			return retry(f, retryTimes-1)
		}
		return result, err
	} else {
		return result, nil
	}

}
