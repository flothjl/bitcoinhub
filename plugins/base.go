package plugins

import (
	"errors"
)

type ProviderDataItem struct {
	Label string
	Value string
}

type ProviderData struct {
	ProviderName string
	Data         []ProviderDataItem
}

type Plugin interface {
	Render() (ProviderData, error)
	GetName() string
}

type PluginManager struct {
	plugins []Plugin
}

func (pm *PluginManager) Register(plugin Plugin) {
	pm.plugins = append(pm.plugins, plugin)
}

func (pm *PluginManager) RenderAll() ([]ProviderData, error) {
	var pData []ProviderData

	for _, plugin := range pm.plugins {
		d, err := plugin.Render()
		if err != nil {
			return pData, err
		}
		pData = append(pData, d)
	}

	return pData, nil
}

func (pm *PluginManager) FindPluginByName(name string) (Plugin, error) {
	for _, p := range pm.plugins {
		if p.GetName() == name {
			return p, nil
		}
	}
	return nil, errors.New("Unable to find plugin")
}
