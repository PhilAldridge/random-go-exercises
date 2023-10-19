package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"
	"sync"
)

func main() {
	http.HandleFunc("/price/", func(w http.ResponseWriter, r *http.Request) {
		begin := time.Now()
		cryptoName := strings.SplitN(r.URL.Path, "/", 3)[2]
		data, err := querySinglePrice(cryptoName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"response": data,
			"took":     time.Since(begin).String(),
		})
	})

	http.HandleFunc("/allprices/", func(w http.ResponseWriter, r *http.Request) {
		begin := time.Now()
		var sliceOfResponses []singleCryptoPrice
		listOfCrypto := [4]string{"BTC", "ETH", "DOGE", "SOL"}
		for _, crypto := range listOfCrypto {
			data, err := querySinglePrice(crypto)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			sliceOfResponses = append(sliceOfResponses, data)
		}
                json.NewEncoder(w).Encode(map[string]interface{}{
			"response": sliceOfResponses,
			"took":     time.Since(begin).String(),
		})
	})

	http.HandleFunc("/allprices-concurrent/", func(w http.ResponseWriter, r *http.Request) {
		begin := time.Now()
		listOfCrypto := [4]string{"BTC", "ETH", "DOGE", "SOL"}
		var responses sync.Map
		wg := sync.WaitGroup{}
		for idx, crypto := range listOfCrypto {
			wg.Add(1)
			go func(crypto string, idx int) {
				data, err := querySinglePrice(crypto)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				responses.Store(idx, data)
				wg.Done()
			}(crypto, idx)
		}
		wg.Wait()
		var sliceOfResponses []singleCryptoPrice
		for i := 0; i < 4; i++ {
			data, _ := responses.Load(i)
			sliceOfResponses = append(sliceOfResponses, data.(singleCryptoPrice))
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"response": sliceOfResponses,
			"took":     time.Since(begin).String(),
		})
	})

	http.ListenAndServe(":8080", nil)
}

func querySinglePrice(cryptoName string) (singleCryptoPrice, error) {
	api_key := "7FDEA7F3-F76E-402C-902C-940C0BD7F2E0"
	c := http.Client{Timeout: time.Duration(3) * time.Second}
	req, err := http.NewRequest("GET", "https://rest.coinapi.io/v1/exchangerate/"+cryptoName+"/USD", nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("X-CoinAPI-Key", api_key)
	resp, err := c.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	var cryptoPricePlaceholder singleCryptoPrice
	err = json.NewDecoder(resp.Body).Decode(&cryptoPricePlaceholder)
	if err != nil {
		return singleCryptoPrice{}, err
	}
	return cryptoPricePlaceholder, nil
}

type singleCryptoPrice struct {
	ID          string  `json:"asset_id_base"`
	Currency    string  `json:"asset_id_quote"`
	MarketPrice float64 `json:"rate"`
}
