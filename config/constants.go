package config

import "time"

const VERSION = "1.1.241"

const REQUEST_TIMEOUT = time.Duration(5 * time.Second)

const MINIMUM_NODE_COUNT = 6

const DEFAULT_NETWORK = "jalapeno"

func getJalapenoEndpoints() [10]string {
	return [10]string{
		"https://node2.litgateway.com:7370",
		"https://node2.litgateway.com:7371",
		"https://node2.litgateway.com:7372",
		"https://node2.litgateway.com:7373",
		"https://node2.litgateway.com:7374",
		"https://node2.litgateway.com:7375",
		"https://node2.litgateway.com:7376",
		"https://node2.litgateway.com:7377",
		"https://node2.litgateway.com:7378",
		"https://node2.litgateway.com:7379",
	}
}

func getSerranoEndpoints() [10]string {
	return [10]string{
		"https://serrano.litgateway.com:7370",
		"https://serrano.litgateway.com:7371",
		"https://serrano.litgateway.com:7372",
		"https://serrano.litgateway.com:7373",
		"https://serrano.litgateway.com:7374",
		"https://serrano.litgateway.com:7375",
		"https://serrano.litgateway.com:7376",
		"https://serrano.litgateway.com:7377",
		"https://serrano.litgateway.com:7378",
		"https://serrano.litgateway.com:7379",
	}
}

func getLocalHostEndpoints() [10]string {
	return [10]string{
		"http://localhost:7470",
		"http://localhost:7471",
		"http://localhost:7472",
		"http://localhost:7473",
		"http://localhost:7474",
		"http://localhost:7475",
		"http://localhost:7476",
		"http://localhost:7477",
		"http://localhost:7478",
		"http://localhost:7479",
	}
}

var NETWORKS = map[string][10]string{
	"jalapeno":  getJalapenoEndpoints(),
	"serrano":   getSerranoEndpoints(),
	"localhost": getLocalHostEndpoints(),
}
