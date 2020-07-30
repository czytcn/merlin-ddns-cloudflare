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

	var client http.Client
	client.Timeout = 100 * time.Second
	get, errGetIp := client.Get(config.O.App.GetIpFromUrl)
	if errGetIp != nil {
		NotifyDdnsState(false)
		fmt.Println(errGetIp)
		return
	}
	all, errParseIp := ioutil.ReadAll(get.Body)
	if errParseIp != nil {
		fmt.Println(errParseIp.Error())
		NotifyDdnsState(false)
		return
	}
	defer get.Body.Close()
	ip := string(all)
	for _, r := range recs {
		fmt.Println(r.Name)
		err := api.UpdateDNSRecord(zoneID, r.ID, cloudflare.DNSRecord{Content: ip, Type: "A"})
		if err != nil {
			NotifyDdnsState(false)
			fmt.Println(err.Error())
			return
		}
	}

	NotifyDdnsState(true)
}

func NotifyDdnsState(success bool) {
	if success {
		exec.Command("/sbin/ddns_custom_updated", "1")
	} else {
		exec.Command("/sbin/ddns_custom_updated", "0")
	}
}
