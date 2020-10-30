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
	api, err := cloudflare.New(config.Obj.Api.ApiToken, config.Obj.Api.Email)
	if err != nil {
		NotifyDdnsState(false)
		log.Fatal(err)
		return
	}

	zoneID, err := api.ZoneIDByName(config.Obj.Api.Domain)
	if err != nil {
		NotifyDdnsState(false)
		log.Fatal(err)
		return
	}
	// Fetch all records for a zone
	recs, err := api.DNSRecords(zoneID, cloudflare.DNSRecord{Name: config.Obj.Api.SubDomain, Type: "A"})
	if err != nil {
		NotifyDdnsState(false)
		log.Fatal(err)
		return
	}
	for _, r := range recs {
		var (
			ip  string
			err error
		)

		switch r.Type {
		case "A":

			ip, err = getIPV4()
			if err != nil || !config.Obj.App.UpdateIpv4DNSRecord {
				log.Fatal(err)
				continue
			}
		case "AAAA":
			ip, err = getIpV6()
			if err != nil || !config.Obj.App.UpdateIpv6DNSRecord {
				log.Fatal(err)
				continue
			}
		}
		if ip == "" {
			log.Fatal("Ip not fetch")
			continue
		}
		updateDnsErr := api.UpdateDNSRecord(zoneID, r.ID, cloudflare.DNSRecord{
			Type:    r.Type,
			Content: ip,
			TTL:     config.Obj.App.TTL,
			Proxied: config.Obj.App.EnableCloudflareProxied,
		})
		if updateDnsErr != nil {
			log.Fatal(updateDnsErr)
			NotifyDdnsState(false)
		} else {
			log.Println("update dns record ok,type:", r.Type)
			NotifyDdnsState(true)
		}
	}
}
func NotifyDdnsState(isSuccess bool) {
	if config.Obj.App.SendDdnsCustomUpdatedNotify {
		if isSuccess {
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

}

func getIp(echoUrl string) (string, error) {
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

func getIPV4() (string, error) {
	return getIp("https://api-ipv4.ip.sb/ip")
}

func getIpV6() (string, error) {
	return getIp("https://api-ipv6.ip.sb/ip")
}
