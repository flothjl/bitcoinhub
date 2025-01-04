package plugins

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
