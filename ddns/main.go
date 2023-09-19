package main

import (
	"ddns/cloudflare"
	"ddns/utils"
	"fmt"
	"os"
)

var CF_DATA_CTX map[string]string

func init() {
	fmt.Println("[Info  ]  Initializing...")
	cfData, err := utils.GetCfContext()
	if err != nil {
		fmt.Printf("[Error ]  %s\n", err.Error())
		fmt.Println("[Error ]  Program encountered a fatal error, aborting...")
		os.Exit(-1)
	}
	CF_DATA_CTX = cfData
}

func main() {
	fmt.Println("[Info  ]  This project aim to create a record on cloudflare of your domain and your machine's ip automatic.")
	fmt.Println("[Info  ]  Trying to get external ip of your machine...")

	ipAddr, err := utils.GetExternalIp()

	if err != nil {
		fmt.Printf("[Error ]  Cannot get external ip, reason: %s\n", err.Error())
		fmt.Println("[Error ]  Program encountered a fatal error, aborting...")
		os.Exit(-1)
	}

	fmt.Printf("[Info  ]  Get external ip: %s\n", ipAddr)

	cfApi, err := cloudflare.NewApi(
		CF_DATA_CTX["apiKey"],
		CF_DATA_CTX["zoneName"],
		CF_DATA_CTX["subDomainName"],
	)

	if err != nil {
		fmt.Printf("[Error ]  Initialize cloudflare api failed, reason: %s\n", err.Error())
		fmt.Println("[Error ]  Program encountered a fatal error, aborting...")
		os.Exit(-1)
	}

	if err := cfApi.RenewDns(ipAddr); err != nil {
		fmt.Printf("[Error ]  Cannot renew the dns record, reason: %s\n", err.Error())
		fmt.Println("[Error ]  Program encountered a fatal error, aborting...")
		os.Exit(-1)
	}

	fmt.Printf("[Info  ]  Dns record of %s.%s map to %s successfully\n", CF_DATA_CTX["subDomainName"], CF_DATA_CTX["zoneName"], ipAddr)
}
