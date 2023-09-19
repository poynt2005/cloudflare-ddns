package utils

import (
	"net/http"
	"strings"
)

func GetExternalIp() (string, error) {
	request, err := http.NewRequest("GET", "https://ident.me", nil)
	request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/117.0")
	if err != nil {
		return "", err
	}

	resp, err := http.DefaultClient.Do(request)

	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBuffers := []byte{}
	for {
		buf := make([]byte, 2048)
		readBytes, err := resp.Body.Read(buf)
		if err != nil {
			break
		}
		respBuffers = append(respBuffers, buf[:readBytes]...)
	}

	return strings.Trim(string(respBuffers), " "), nil
}
