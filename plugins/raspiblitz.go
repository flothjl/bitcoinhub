package plugins

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jritsema/gotoolbox"
)

func HttpGetJSON(url string, result interface{}, headers map[string]string) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("cannot create request: %w", err)
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("cannot fetch URL %q: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected http GET status: %s", resp.Status)
	}

	err = json.NewDecoder(resp.Body).Decode(result)
	if err != nil {
		return fmt.Errorf("cannot decode JSON: %w", err)
	}

	return nil
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

func buildRaspiblitzData() (ProviderData, error) {
	url := fmt.Sprintf("%s/%s", os.Getenv("RASPIBLITZ_URL"), "system/hardware-info")
	const ProviderName = "Raspiblitz"
	data := ProviderData{ProviderName, []ProviderDataItem{}}

	jwt, authErr := raspiBlitzAuth()
	if authErr != nil {
		log.Printf("failure to auth: %v", authErr)
		return data, authErr
	}

	raspiblitzData := &RaspiBlitzStatus{}
	err := HttpGetJSON(url, raspiblitzData, map[string]string{"Authorization": fmt.Sprintf("Bearer %s", jwt)})
	if err != nil {
		log.Printf("%v", err)
		return data, err
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

	return data, nil
}

type RaspiblitzPlugin struct{}

func (p RaspiblitzPlugin) Render() (ProviderData, error) {
	return buildRaspiblitzData()
}
