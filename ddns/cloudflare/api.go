package cloudflare

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/http"
	"regexp"

	"github.com/bitly/go-simplejson"
)

type Api struct {
	dnsZoneName      string
	dnsSubDomainName string
	dnsApiKey        string

	zoneId      string
	subdomainId string
}

func NewApi(apiKey, zoneName, subDomainName string) (*Api, error) {
	instance := &Api{
		zoneName,
		subDomainName,
		apiKey,
		"",
		"",
	}

	if err := instance.getZoneId(); err != nil {
		return nil, err
	}

	if err := instance.findTargetSubDomain(); err != nil {
		return nil, err
	}

	return instance, nil
}

func (api *Api) parseResult(jsonResp *simplejson.Json) error {
	successType, _ := jsonResp.Get("success").Bool()

	if !successType {
		cfErrors, _ := jsonResp.Get("errors").Array()

		firstErrorType := cfErrors[0]
		firstError, ok := firstErrorType.(map[string]interface{})

		if !ok {
			return errors.New("Unknown Error, cannot type assert the first error map")
		}

		messageType, ok := firstError["message"]

		if !ok {
			return errors.New("Unknown Error, cannot retrieve message from error map")
		}

		message, ok := messageType.(string)

		if !ok {
			return errors.New("Unknown Error, cannot type assert message type to string")
		}

		return errors.New(fmt.Sprintf("Cloudflare Error, %s", message))
	}

	return nil
}

func (api *Api) getZoneId() error {
	/**
		curl --request GET \
	         --url https://api.cloudflare.com/client/v4/zones?page={page} \
	         --header 'Content-Type: application/json' \
	         --header 'Authorization: Bearer {apiKey}'
	*/

	var recurGetPage func(page int) (string, error)

	recurGetPage = func(page int) (string, error) {
		request, err := http.NewRequest("GET", fmt.Sprintf("https://api.cloudflare.com/client/v4/zones?page=%d", page), nil)

		if err != nil {
			return "", err
		}

		request.Header.Add("Content-Type", "application/json")
		request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", api.dnsApiKey))

		resp, err := http.DefaultClient.Do(request)

		if err != nil {
			return "", err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return "", errors.New(fmt.Sprintf("Error Return Code: %d", resp.StatusCode))
		}

		jsonResp, err := simplejson.NewFromReader(resp.Body)

		if err != nil {
			return "", err
		}

		if err := api.parseResult(jsonResp); err != nil {
			return "", err
		}

		resultArr, _ := jsonResp.Get("result").Array()

		if len(resultArr) == 0 {
			return "", errors.New(fmt.Sprintf("Target Zone name: %s not found", api.dnsZoneName))
		}

		for _, resultMapType := range resultArr {
			resultMap, ok := resultMapType.(map[string]interface{})

			if !ok {
				continue
			}

			nameType, ok := resultMap["name"]

			if !ok {
				continue
			}

			name, ok := nameType.(string)

			if !ok {
				continue
			}

			if name == api.dnsZoneName {
				idType, ok := resultMap["id"]

				if !ok {
					continue
				}

				id, ok := idType.(string)

				if !ok {
					continue
				}

				return id, nil
			}
		}

		return recurGetPage(page + 1)
	}

	zoneId, err := recurGetPage(1)

	if err != nil {
		return err
	}

	api.zoneId = zoneId

	return nil
}

func (api *Api) findTargetSubDomain() error {
	/**
	curl --request GET \
	 	 --url https://api.cloudflare.com/client/v4/zones/${DnsZoneID}/dns_records \
		 --header 'Content-Type: application/json' \
		 --header 'Authorization: Bearer ApiKey'
	*/
	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records", api.zoneId)
	request, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return err
	}

	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", api.dnsApiKey))

	resp, err := http.DefaultClient.Do(request)

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("Error Return Code: %d", resp.StatusCode))
	}

	jsonResp, err := simplejson.NewFromReader(resp.Body)

	if err != nil {
		return err
	}

	if err := api.parseResult(jsonResp); err != nil {
		return err
	}

	resultArr, _ := jsonResp.Get("result").Array()

	if len(resultArr) == 0 {
		// Subdomain not found
		return nil
	}

	domainPattern, _ := regexp.Compile(fmt.Sprintf("(\\.)?%s$", api.dnsZoneName))

	for _, subDomainMapType := range resultArr {
		subDomainMap, ok := subDomainMapType.(map[string]interface{})

		if !ok {
			continue
		}

		recordNameType, ok := subDomainMap["name"]

		if !ok {
			continue
		}

		recordName, ok := recordNameType.(string)

		if !ok {
			continue
		}

		if api.dnsSubDomainName == domainPattern.ReplaceAllString(recordName, "") {
			idType, ok := subDomainMap["id"]

			if !ok {
				continue
			}

			id, ok := idType.(string)

			if !ok {
				continue
			}

			// Subdomain found
			api.subdomainId = id
			return nil
		}
	}

	// Subdomain not found
	return nil
}

func (api *Api) RenewDns(ipAddr string) error {
	if net.ParseIP(ipAddr) == nil {
		return errors.New(fmt.Sprintf("Given address: %s, is not a valid IPv4 address", ipAddr))
	}

	var request *http.Request = nil

	reqBody := simplejson.New()

	reqBody.Set("content", ipAddr)
	reqBody.Set("name", api.dnsSubDomainName)
	reqBody.Set("proxied", false)
	reqBody.Set("type", "A")

	if len(api.subdomainId) == 0 {
		// Create
		/**
		curl --request POST \
		  --url https://api.cloudflare.com/client/v4/zones/zone_identifier/dns_records \
		  --header 'Content-Type: application/json' \
		  --header 'Authorization: Bearer ApiKey' \
		  --data '{
		      "content": "<YOUR_IP>",
		      "name": "<YOUR_MAIN_DOMAIN>",
		      "proxied": false,
		      "type": "A",
		      "comment": "Cloudflare DDNS project",
		  }'
		*/

		url := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records", api.zoneId)
		payload, _ := reqBody.MarshalJSON()
		req, err := http.NewRequest("POST", url, bytes.NewReader(payload))

		if err != nil {
			return err
		}

		request = req
	} else {
		// Update
		/**
		curl --request PUT \
		  --url https://api.cloudflare.com/client/v4/zones/${DnsZoneId}/dns_records/${DnsSubDomainId} \
		  --header 'Content-Type: application/json' \
		  --header 'Authorization: Bearer ApiKey' \
		  --data '{
		      "content": "<YOUR_IP>",
		      "name": "<YOUR_MAIN_DOMAIN>",
		      "proxied": false,
		      "type": "A",
		      "comment": "Cloudflare DDNS project",
		  }'
		*/

		url := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records/%s", api.zoneId, api.subdomainId)
		payload, _ := reqBody.MarshalJSON()
		req, err := http.NewRequest("PUT", url, bytes.NewReader(payload))

		if err != nil {
			return err
		}

		request = req
	}

	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", api.dnsApiKey))

	resp, err := http.DefaultClient.Do(request)

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// if resp.StatusCode != http.StatusOK {
	// 	return errors.New(fmt.Sprintf("Error Return Code: %d", resp.StatusCode))
	// }

	jsonResp, err := simplejson.NewFromReader(resp.Body)

	if err != nil {
		return err
	}

	if err := api.parseResult(jsonResp); err != nil {
		return err
	}

	return nil
}
