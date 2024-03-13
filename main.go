package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-yaml/yaml"
	"github.com/gorilla/mux"
)

type Config struct {
	Port  string `yaml:"port"`
	Links []Link `yaml:"links"`
}

type Link struct {
	Name string `yaml:"name"`
	URL  string `yaml:"url"`
	Type int    `yaml:"type"`
}

var (
	configPath string
	example    bool
	config     Config
	shortLinks map[string]Link
)

func init() {
	flag.StringVar(&configPath, "c", "", "Path to config file")
	flag.BoolVar(&example, "i", false, "Generate example config file")
}

func main() {
	flag.Parse()

	if example {
		generateExampleConfig()
		return
	}

	if configPath == "" {
		flag.PrintDefaults()
		return
	}

	loadConfig()

	shortLinks = make(map[string]Link)
	for _, link := range config.Links {
		shortLinks[link.Name] = link
	}

	r := mux.NewRouter()
	r.HandleFunc("/{name:.*}", redirectShortLink)
	http.Handle("/", r)

	port := ":" + config.Port
	fmt.Printf("Short links Server is running on port %s...\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

func redirectShortLink(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	link, ok := shortLinks[name]
	if !ok {
		http.Error(w, "Short link not found", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, link.URL, link.Type)
}

func loadConfig() {
	yamlFile, err := os.ReadFile(configPath)
	if err != nil {
		fmt.Printf("Error reading YAML file: %s\n", err)
		os.Exit(1)
	}

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		fmt.Printf("Error parsing YAML file: %s\n", err)
		os.Exit(1)
	}
}

func generateExampleConfig() {
	exampleConfig := `port: 8080
links:
  - name: "abc"
    url: "http://example.cc"
    type: 302
  - name: "abcd"
    url: "http://example1.cc"
    type: 302
`
	if configPath == "" {
		configPath = "example_config.yml"
		return
	}
	err := os.WriteFile(configPath, []byte(exampleConfig), 0644)
	if err != nil {
		fmt.Printf("Error generating example config file: %s\n", err)
		os.Exit(1)
	}
	fmt.Println("Example config file generated: " + configPath)
}
