package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/go-yaml/yaml"
	"github.com/gorilla/mux"
)

type Config struct {
	Port  string       `yaml:"port"`
	Links []LinkConfig `yaml:"links"`
}

type URLGroup struct {
	Name string `yaml:"name"`
	Url  string `yaml:"url"`
}

type ReplaceGroup struct {
	QueryName string `yaml:"query_name"`
	Keyword   string `yaml:"keyword"`
}

type ReplaceKeyword struct {
	Keyword string
	Value   string
}

type LinkConfig struct {
	Name          string          `yaml:"name"`
	Type          int             `yaml:"type"`
	URL           *string         `yaml:"url"`
	URLs          *[]URLGroup     `yaml:"urls"`
	URLsQueryName *string         `yaml:"urls_query_name"`
	Replace       *[]ReplaceGroup `yaml:"replace"`
}

type LinkRule struct {
	Name          string             `yaml:"name"`
	Type          int                `yaml:"type"`
	IsGroup       bool               `yaml:"is_group"`
	URL           *string            `yaml:"url"`
	URLs          *map[string]string `yaml:"urls"`
	URLsQueryName *string            `yaml:"urls_query_name"`
	Replace       *[]ReplaceGroup    `yaml:"replace"`
}

var (
	configPath string
	example    bool
	config     Config
	shortLinks map[string]LinkRule
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

	// shortLinks = make(map[string]Link)
	// for _, link := range config.Links {
	// 	shortLinks[link.Name] = link
	// }
	generateRule()

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

	targetUrl := ""
	if link.IsGroup && shortLinks[name].URLs != nil {
		queryName := *link.URLsQueryName
		queryValue := r.URL.Query().Get(queryName)
		if queryValue == "" {
			http.Error(w, "Query parameter not found", http.StatusBadRequest)
			return
		}

		if shortLinks[name].URLs != nil && (*shortLinks[name].URLs)[queryValue] == "" {
			http.Error(w, "Query parameter not found", http.StatusBadRequest)
			return
		}

		url := (*shortLinks[name].URLs)[queryValue]
		if link.Replace != nil {
			url = replaceURL(url, r.URL.Query(), *link.Replace)
		}
		targetUrl = url
	} else {
		if link.URL == nil {
			http.Error(w, "URL not found", http.StatusBadRequest)
			return
		}
		targetUrl = *link.URL
		if link.Replace != nil {
			url := replaceURL(*link.URL, r.URL.Query(), *link.Replace)
			targetUrl = url
		}
	}

	http.Redirect(w, r, targetUrl, link.Type)
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

func generateRule() {
	shortLinks = make(map[string]LinkRule)
	for _, clink := range config.Links {
		rule := LinkRule{
			Name:    clink.Name,
			Type:    clink.Type,
			Replace: clink.Replace,
		}
		if clink.URLs != nil && clink.URLsQueryName != nil {
			rule.URLsQueryName = clink.URLsQueryName
			rule.IsGroup = true
			urls := make(map[string]string)
			for _, url := range *clink.URLs {
				urls[url.Name] = url.Url
			}
			rule.URLs = &urls
		} else {
			rule.URL = clink.URL
		}
		shortLinks[clink.Name] = rule
	}

	// jsonData, _ := json.Marshal(shortLinks)
	// fmt.Println("Generate rule success!", string(jsonData))
}

func replaceURL(url string, querys url.Values, replace []ReplaceGroup) string {
	for _, r := range replace {
		url = strings.ReplaceAll(url, r.Keyword, querys.Get(r.QueryName))
	}
	return url
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
	}
	err := os.WriteFile(configPath, []byte(exampleConfig), 0644)
	if err != nil {
		fmt.Printf("Error generating example config file: %s\n", err)
		os.Exit(1)
	}
	fmt.Println("Example config file generated: " + configPath)
}
