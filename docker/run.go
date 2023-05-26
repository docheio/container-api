package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"golang.org/x/net/http2"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
)

func main() {
	request, _ := http.NewRequest("GET", "https://www.minecraft.net/en-us/download/server/bedrock", nil)
	request.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.198 Safari/537.36")
	request.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	request.Header.Add("Accept-Language", "en-us")
	request.Header.Add("Connection", "keep-alive")
	request.Header.Add("Host", "www.google.com")
	client := &http.Client{}
	caCert, _ := ioutil.ReadFile("/etc/ssl/cert.pem")
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	client.Transport = &http2.Transport{
		TLSClientConfig: &tls.Config{
			RootCAs: caCertPool,
		},
	}
	response, err := client.Do(request)
	if err != nil {
		return
	}
	defer response.Body.Close()
	r := regexp.MustCompile(`https://minecraft.azureedge.net/bin-linux/bedrock-server-.+\.zip`)
	b, _ := io.ReadAll(response.Body)
	newLink := r.FindString(string(b))
	fmt.Println(newLink)
}
