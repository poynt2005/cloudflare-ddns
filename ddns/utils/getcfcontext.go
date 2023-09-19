package utils

import (
	"errors"
	"os"
	"strings"
)

func GetCfContext() (map[string]string, error) {
	retVal := map[string]string{
		"apiKey":        "",
		"zoneName":      "",
		"subDomainName": "",
	}

	for _, environKV := range os.Environ() {
		environSplit := strings.Split(environKV, "=")

		switch environSplit[0] {
		case "CLOUDFLARE_DNS_API_KEY":
			retVal["apiKey"] = environSplit[1]
			break
		case "CLOUDFLARE_DNS_ZONE_NAME":
			retVal["zoneName"] = environSplit[1]
			break
		case "CLOUDFLARE_DNS_SUBDOMAIN_NAME":
			retVal["subDomainName"] = environSplit[1]
			break
		}
	}

	if len(retVal["apiKey"]) == 0 || len(retVal["zoneName"]) == 0 || len(retVal["subDomainName"]) == 0 {
		return nil, errors.New("You must provide CLOUDFLARE_DNS_API_KEY, CLOUDFLARE_DNS_ZONE_NAME and CLOUDFLARE_DNS_SUBDOMAIN_NAME")
	}

	return retVal, nil
}
