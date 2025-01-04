package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jritsema/gotoolbox"
	"github.com/jritsema/gotoolbox/web"
)

const (
	raspiblitzUrl = "http://192.168.69.65/api/system/hardware-info"
)

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

type RaspiBlitzStatus struct {
	CPUOverallPercent   float64   `json:"cpu_overall_percent,omitempty"`
	CPUPerCPUPercent    []float64 `json:"cpu_per_cpu_percent,omitempty"`
	VramTotalBytes      int64     `json:"vram_total_bytes,omitempty"`
	VramAvailableBytes  int64     `json:"vram_available_bytes,omitempty"`
	VramUsedBytes       int64     `json:"vram_used_bytes,omitempty"`
	VramUsagePercent    float64   `json:"vram_usage_percent,omitempty"`
	TemperaturesCelsius struct {
		SystemTemp float64 `json:"system_temp,omitempty"`
		Coretemp   []any   `json:"coretemp,omitempty"`
	} `json:"temperatures_celsius,omitempty"`
	BootTimeTimestamp float64 `json:"boot_time_timestamp,omitempty"`
	Networks          struct {
		InternetOnline       string `json:"internet_online,omitempty"`
		TorWebAddr           string `json:"tor_web_addr,omitempty"`
		InternetLocalip      string `json:"internet_localip,omitempty"`
		InternetLocaliprange string `json:"internet_localiprange,omitempty"`
	} `json:"networks,omitempty"`
	Disks []struct {
		Device              string  `json:"device,omitempty"`
		Mountpoint          string  `json:"mountpoint,omitempty"`
		FilesystemType      string  `json:"filesystem_type,omitempty"`
		PartitionTotalBytes int64   `json:"partition_total_bytes,omitempty"`
		PartitionUsedBytes  int64   `json:"partition_used_bytes,omitempty"`
		PartitionFreeBytes  int64   `json:"partition_free_bytes,omitempty"`
		PartitionPercent    float64 `json:"partition_percent,omitempty"`
	} `json:"disks,omitempty"`
}

type ProviderDataItem struct {
	Label string
	Value string
}

type ProviderData struct {
	ProviderName string
	Data         []ProviderDataItem
}

func buildBitAxeProviderData() ProviderData {
	url := os.Getenv("BITAXE_URL")
	const ProviderName = "BitAxe"
	var data ProviderData
	data.ProviderName = ProviderName
	// get bitaxe data
	bitAxeData := &BitAxeStatus{}
	err := gotoolbox.HttpGetJSON(url, bitAxeData)
	if err != nil {
		return data
	}
	items := []ProviderDataItem{}

	items = append(items, ProviderDataItem{"Power", fmt.Sprintf("%.2f mW", bitAxeData.Power)})
	items = append(items, ProviderDataItem{"Voltage", fmt.Sprintf("%.2f mV", bitAxeData.Voltage)})
	items = append(items, ProviderDataItem{"Current", fmt.Sprintf("%.2f mA", bitAxeData.Current)})
	items = append(items, ProviderDataItem{"Temp", fmt.Sprintf("%d", bitAxeData.Temp)})
	items = append(items, ProviderDataItem{"Hash Rate", fmt.Sprintf("%.2f GH/s", bitAxeData.HashRate)})
	items = append(items, ProviderDataItem{"Best Difficulty", bitAxeData.BestDiff})
	items = append(items, ProviderDataItem{"Best Session Difficulty", bitAxeData.BestSessionDiff})

	return ProviderData{ProviderName, items}
}

func raspiBlitzAuth() (string, error) {
	url := fmt.Sprintf("%s/%s", os.Getenv("RASPIBLITZ_URL"), "system/login")

	var jwt string

	password := os.Getenv("RASPIBLITZ_PASS")
	otp := os.Getenv("RASPIBLITZ_OTP")
	loginRequest := map[string]string{
		"password":          password,
		"one_time_password": otp,
	}

	err := gotoolbox.HttpPostJSON(url, loginRequest, &jwt, 200)
	return jwt, err
}

func buildRaspiblitzData() ProviderData {
	url := fmt.Sprintf("%s/%s", os.Getenv("RASPIBLITZ_URL"), "system/hardware-info")
	const ProviderName = "Raspiblitz"
	data := ProviderData{ProviderName, []ProviderDataItem{}}

	jwt, authErr := raspiBlitzAuth()
	if authErr != nil {
		log.Printf("failure to auth: %v", authErr)
		return data
	}

	raspiblitzData := &RaspiBlitzStatus{}
	err := HttpGetJSON(url, raspiblitzData, map[string]string{"Authorization": fmt.Sprintf("Bearer %s", jwt)})
	if err != nil {
		log.Printf("%v", err)
		return data
	}

	items := []ProviderDataItem{}
	items = append(items, ProviderDataItem{"CPU", fmt.Sprintf("%.2f %%", raspiblitzData.CPUOverallPercent)})
	items = append(items, ProviderDataItem{"RAM", fmt.Sprintf("%.2f %%", raspiblitzData.VramUsagePercent)})
	for i, v := range raspiblitzData.Disks {
		pd := ProviderDataItem{
			Label: fmt.Sprintf("Disk %d", i),
			Value: fmt.Sprintf(
				"%dGB/%dGB (%.2f %%)",
				v.PartitionUsedBytes/1000000000,
				v.PartitionTotalBytes/1000000000, float64(v.PartitionUsedBytes)/float64(v.PartitionTotalBytes)*100.00),
		}

		items = append(items, pd)
	}
	data.Data = items

	return data
}

func index(r *http.Request) *web.Response {
	var data []ProviderData

	data = append(data, buildBitAxeProviderData())
	data = append(data, buildRaspiblitzData())
	log.Printf("%v", data)

	return web.HTML(http.StatusOK, html, "index.html", data, nil)
}
