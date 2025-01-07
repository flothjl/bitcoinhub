package plugins

import (
	"fmt"
	"os"

	"github.com/jritsema/gotoolbox"
)

const bitaxeName = "BitAxe"

type BitAxeStatus struct {
	Power                  float64 `json:"power,omitempty"`
	Voltage                float64 `json:"voltage,omitempty"`
	Current                float64 `json:"current,omitempty"`
	Temp                   int     `json:"temp,omitempty"`
	VrTemp                 int     `json:"vrTemp,omitempty"`
	HashRate               float64 `json:"hashRate,omitempty"`
	BestDiff               string  `json:"bestDiff,omitempty"`
	BestSessionDiff        string  `json:"bestSessionDiff,omitempty"`
	StratumDiff            int     `json:"stratumDiff,omitempty"`
	IsUsingFallbackStratum int     `json:"isUsingFallbackStratum,omitempty"`
	FreeHeap               int     `json:"freeHeap,omitempty"`
	CoreVoltage            int     `json:"coreVoltage,omitempty"`
	CoreVoltageActual      int     `json:"coreVoltageActual,omitempty"`
	Frequency              int     `json:"frequency,omitempty"`
	Ssid                   string  `json:"ssid,omitempty"`
	MacAddr                string  `json:"macAddr,omitempty"`
	Hostname               string  `json:"hostname,omitempty"`
	WifiStatus             string  `json:"wifiStatus,omitempty"`
	SharesAccepted         int     `json:"sharesAccepted,omitempty"`
	SharesRejected         int     `json:"sharesRejected,omitempty"`
	UptimeSeconds          int     `json:"uptimeSeconds,omitempty"`
	AsicCount              int     `json:"asicCount,omitempty"`
	SmallCoreCount         int     `json:"smallCoreCount,omitempty"`
	ASICModel              string  `json:"ASICModel,omitempty"`
	StratumURL             string  `json:"stratumURL,omitempty"`
	FallbackStratumURL     string  `json:"fallbackStratumURL,omitempty"`
	StratumPort            int     `json:"stratumPort,omitempty"`
	FallbackStratumPort    int     `json:"fallbackStratumPort,omitempty"`
	StratumUser            string  `json:"stratumUser,omitempty"`
	FallbackStratumUser    string  `json:"fallbackStratumUser,omitempty"`
	Version                string  `json:"version,omitempty"`
	BoardVersion           string  `json:"boardVersion,omitempty"`
	RunningPartition       string  `json:"runningPartition,omitempty"`
	Flipscreen             int     `json:"flipscreen,omitempty"`
	OverheatMode           int     `json:"overheat_mode,omitempty"`
	Invertscreen           int     `json:"invertscreen,omitempty"`
	Invertfanpolarity      int     `json:"invertfanpolarity,omitempty"`
	Autofanspeed           int     `json:"autofanspeed,omitempty"`
	Fanspeed               int     `json:"fanspeed,omitempty"`
	Fanrpm                 int     `json:"fanrpm,omitempty"`
}

func buildBitAxeProviderData() (ProviderData, error) {
	url := os.Getenv("BITAXE_URL")
	var data ProviderData
	data.ProviderName = bitaxeName
	// get bitaxe data
	bitAxeData := &BitAxeStatus{}
	err := gotoolbox.HttpGetJSON(url, bitAxeData)
	if err != nil {
		return data, err
	}
	items := []ProviderDataItem{}

	items = append(items, ProviderDataItem{"Power", fmt.Sprintf("%.2f mW", bitAxeData.Power)})
	items = append(items, ProviderDataItem{"Voltage", fmt.Sprintf("%.2f mV", bitAxeData.Voltage)})
	items = append(items, ProviderDataItem{"Current", fmt.Sprintf("%.2f mA", bitAxeData.Current)})
	items = append(items, ProviderDataItem{"Temp", fmt.Sprintf("%d", bitAxeData.Temp)})
	items = append(items, ProviderDataItem{"Hash Rate", fmt.Sprintf("%.2f GH/s", bitAxeData.HashRate)})
	items = append(items, ProviderDataItem{"Best Difficulty", bitAxeData.BestDiff})
	items = append(items, ProviderDataItem{"Best Session Difficulty", bitAxeData.BestSessionDiff})

	return ProviderData{bitaxeName, items}, nil
}

type BitAxePlugin struct{}

func (p BitAxePlugin) Render() (ProviderData, error) {
	d, err := buildBitAxeProviderData()
	return d, err
}

func (p BitAxePlugin) GetName() string {
	return bitaxeName
}
