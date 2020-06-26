package main

import (
	"log"
	"os"
	"redis-proxy/src/api"
	"redis-proxy/src/proxy"
	"redis-proxy/src/resp"
	"time"
)

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func parameters() (*string, *string, *string) {
	apiPort := getEnv("PROXY_API_PORT", ":8080")
	proxyPort := getEnv("PROXY_PORT", ":9999")
	fileLocation := getEnv("CLIENTS_FILE_PATH", "./conf/init.json")

	resp.SENTINEL_ADDRESS = getEnv("SENTINEL_IP", "127.0.0.1") + ":" + getEnv("SENTINEL_PORT", "26379")
	resp.SENTINEL_NAME = getEnv("SENTINEL_MASTER_NAME", "mymaster")

	return &apiPort, &proxyPort, &fileLocation
}

func main() {
	log.Println("starting redis-sentinel proxy")

	apiPort, proxyPort, fileLocation := parameters()

	api := api.NewRouter(apiPort)
	proxy := proxy.NewProxy(proxyPort)
	sentinel := resp.NewSentinel(fileLocation)

	go api.Start()
	go proxy.Start()
	go sentinel.Start()

	for {
		time.Sleep(time.Second * 10)
	}
}
