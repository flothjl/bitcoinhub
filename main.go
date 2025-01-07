package main

import (
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/flothjl/bitcoinhub/plugins"
	"github.com/joho/godotenv"
	"github.com/jritsema/gotoolbox/web"
)

var (
	//go:embed all:templates/*
	templateFS embed.FS

	//go:embed css/*.css
	css embed.FS

	// parsed templates
	html *template.Template

	pm *plugins.PluginManager
)

func init() {
	pm = &plugins.PluginManager{}
	pm.Register(plugins.BitAxePlugin{})
	pm.Register(plugins.RaspiblitzPlugin{})
}

func index(r *http.Request) *web.Response {
	data, err := pm.RenderAll()
	if err != nil {
		log.Printf("Error building plugins: %v", err)
	}

	return web.HTML(http.StatusOK, html, "index.html", data, nil)
}

func refresh(r *http.Request) *web.Response {
	name := r.PathValue("name")
	p, err := pm.FindPluginByName(name)
	if err != nil {
		log.Printf("%v", err)
		return web.HTML(http.StatusInternalServerError, html, "terminal.html", nil, nil)
	}

	data, _ := p.Render()

	return web.HTML(http.StatusOK, html, "terminal.html", data, nil)
}

func main() {
	// exit process immediately upon sigterm
	handleSigTerms()
	var err error

	err = godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// parse templates
	html, err = web.TemplateParseFSRecursive(templateFS, ".html", true, nil)
	if err != nil {
		panic(err)
	}

	router := http.NewServeMux()
	router.Handle("GET /css/output.css", http.FileServer(http.FS(css)))
	router.Handle("GET /css/main.css", http.FileServer(http.FS(css)))
	// refresh
	router.Handle("GET /refresh/{name}", web.Action(refresh))
	// home
	router.Handle("GET /", web.Action(index))
	router.Handle("GET /index.html", web.Action(index))

	// logging/tracing
	nextRequestID := func() string {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	logger := log.New(os.Stdout, "http: ", log.LstdFlags)
	middleware := tracing(nextRequestID)(logging(logger)(router))

	port := GetEnvWithDefault("PORT", "8080")
	logger.Println("listening on http://localhost:" + port)
	if err := http.ListenAndServe(":"+port, middleware); err != nil {
		logger.Println("http.ListenAndServe():", err)
		os.Exit(1)
	}
}

func handleSigTerms() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("received SIGTERM, exiting")
		os.Exit(1)
	}()
}
