package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"time"
)

var client = &http.Client{}
var req *http.Request
var stores = map[string]string{
	"0070014200142": "Mataró",
	"0070025700257": "Ciutat Vella - Barcelona",
	"0070189601896": "Sabadell",
	"0070062300623": "Mollet del vallés",
	"0070054700547": "Sant Adrià de Besos",
	"0070008600086": "Badalona",
	"0070187701877": "La Maquinista - Barcelona",
	"0070011400114": "L'Illa Diagonal - Barcelona",
	"0070140601406": "Gràcia",
	"0070032800328": "Gran Via 2 - L'Hospitalet",
}
var availabilityIcon = map[string]string{
	"noStock": "🚫",
	"inStock": "✅",
}

type availabilityResponse struct {
	ResponseTO responseTO `json:"responseTO"`
}

type responseTO struct {
	StoreAvailabilities []availability `json:"data"`
}

type availability struct {
	AvailabilityInfo string `json:"availabilityInfo"`
	StoreID          string `json:"storeId"`
}

var location, _ = time.LoadLocation("Europe/Rome")

func main() {
	req = buildRequest()

	if req == nil {
		fmt.Println("Can't build request")
		os.Exit(0)
	}

	doRequest()

	ctrlC := make(chan os.Signal, 1)
	signal.Notify(ctrlC, os.Interrupt)

	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			fmt.Println("Running update")
			doRequest()
		case <-ctrlC:
			ticker.Stop()
			fmt.Println("\nFinished with CTRL + C")
			return
		}
	}
}

func doRequest() {
	t := time.Now().In(location)

	fmt.Printf("NOW: %s\n", t)

	if req == nil {
		fmt.Println("No request built!")
		return
	}

	resp, err := client.Do(req)

	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	defer resp.Body.Close()
	respBody, _ := ioutil.ReadAll(resp.Body)

	var availabilityRes availabilityResponse

	err = json.Unmarshal(respBody, &availabilityRes)

	if err != nil {
		fmt.Println(err)
		return
	}

	for _, t := range availabilityRes.ResponseTO.StoreAvailabilities {
		storeName := stores[t.StoreID]
		availIcon := availabilityIcon[t.AvailabilityInfo]

		fmt.Printf("[%s] %s\n", availIcon, storeName)
	}
}

func buildRequest() *http.Request {
	req, err := http.NewRequest("GET", getURLWithParams(), nil)

	if err != nil {
		fmt.Println(err)
		return nil
	}

	return req
}

func getURLWithParams() string {
	base, err := url.Parse("https://www.decathlon.es/es/ajax/rest/model/com/decathlon/cube/commerce/inventory/InventoryActor/getStoreAvailability")
	if err != nil {
		fmt.Println(err)
		return ""
	}

	params := url.Values{}

	params.Add("storeIds", "0070014200142,0070008600086,0070062300623,0070054700547,0070187701877,0070025700257,0070140601406,0070189601896,0070011400114,0070032800328")
	params.Add("skuId", "2524420")
	params.Add("modelId", "8491831")
	params.Add("displayStoreDetails", "false")

	base.RawQuery = params.Encode()

	return base.String()
}
