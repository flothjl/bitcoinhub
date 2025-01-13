package plugins

import (
	"fmt"
	"os"

	"github.com/jritsema/gotoolbox"
)

const name = "mempool"

type MempoolPriceData struct {
	Time int `json:"time,omitempty"`
	Usd  int `json:"USD,omitempty"`
	Eur  int `json:"EUR,omitempty"`
	Gbp  int `json:"GBP,omitempty"`
	Cad  int `json:"CAD,omitempty"`
	Chf  int `json:"CHF,omitempty"`
	Aud  int `json:"AUD,omitempty"`
	Jpy  int `json:"JPY,omitempty"`
}

type MempoolPlugin struct{}

func (p MempoolPlugin) Render() (ProviderData, error) {
	url := os.Getenv("MEMPOOL_URL")
	if url == "" {
		url = "https://mempool.space/api/v1/prices"
	}

	var data ProviderData
	data.ProviderName = name

	mempoolPriceData := &MempoolPriceData{}

	err := gotoolbox.HttpGetJSON(url, mempoolPriceData)
	if err != nil {
		return data, err
	}

	items := []ProviderDataItem{}

	items = append(items, ProviderDataItem{Label: "BTC Price ($)", Value: fmt.Sprintf("$ %d", mempoolPriceData.Usd)})

	return ProviderData{ProviderName: name, Data: items}, nil
}

func (p MempoolPlugin) GetName() string {
	return name
}
