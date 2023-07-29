package tlsc

import (
	tls_client "github.com/bogdanfinn/tls-client"
)

type TLSClient struct {
	tls_client.HttpClient
}

func NewTLSClient() (client tls_client.HttpClient, err error) {
	jar := tls_client.NewCookieJar()

	options := []tls_client.HttpClientOption{
		tls_client.WithTimeoutSeconds(360),
		tls_client.WithClientProfile(tls_client.Firefox_110),
		tls_client.WithNotFollowRedirects(),
		tls_client.WithCookieJar(jar), // create cookieJar instance and pass it as argument
		// Disable SSL verification
		tls_client.WithInsecureSkipVerify(),
	}

	return tls_client.NewHttpClient(tls_client.NewNoopLogger(), options...)

}
