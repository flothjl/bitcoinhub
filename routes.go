package main

import (
	"log"
	"net/http"

	"github.com/flothjl/bitcoinhub/plugins"
	"github.com/jritsema/gotoolbox/web"
)

func index(r *http.Request) *web.Response {
	pm := &plugins.PluginManager{}
	pm.Register(plugins.BitAxePlugin{})
	pm.Register(plugins.RaspiblitzPlugin{})

	data, err := pm.RenderAll()
	if err != nil {
		log.Printf("Error building plugins: %v", err)
	}

	return web.HTML(http.StatusOK, html, "index.html", data, nil)
}
