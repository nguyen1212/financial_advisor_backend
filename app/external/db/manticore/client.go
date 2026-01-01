// Package manticore define client to interact with manticore server
package manticore

import (
	"fmt"
	"sync"

	"github.com/financial_advisor/app/config"
	manticoresearch "github.com/manticoresoftware/manticoresearch-go"
)

var (
	once      sync.Once
	singleton *manticoresearch.APIClient
)

func InitClient(config *config.Config) {
	once.Do(func() {
		cfg := manticoresearch.NewConfiguration()
		cfg.Servers[0].URL = fmt.Sprintf("http://%s:%s", config.ManticoreHost, config.ManticorePort)

		println("Manticore Search URL:", cfg.Servers[0].URL, config.ManticoreHost) // Debug line to print the URL
		singleton = manticoresearch.NewAPIClient(cfg)
	})
}

func Get() *manticoresearch.APIClient {
	return singleton
}
