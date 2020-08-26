package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"gopkg.in/yaml.v2"
)

// Item represents an archive.org item
type Item struct {
	Identifier  string
	D1          string
	D2          string
	Dir         string
	LastUpdated int
	Root        string
}

// PagesMetadata provides the contents of the archive.yml file
type PagesMetadata struct {
	// Version  int    `yaml:"version"`
	ItemType string `yaml:"itemtype"`
	Root     string `yaml:"root"`
}

type ItemMetadata struct {
	Created     int        `json:"created"`
	D1          string     `json:"d1"`
	D2          string     `json:"d2"`
	Dir         string     `json:"dir"`
	Files       []ItemFile `json:"files"`
	LastUpdated int        `json:"item_last_updated"`
	Size        int        `json:"size"`
}

type ItemFile struct {
	Name  string `json:"name"`
	Mtime string `json:"mtime"`
	Size  string `json:"size"`
}

func GetItem(identifier string) (*Item, error) {
	var item *Item
	item = GetCacheItem(identifier)
	if item != nil {
		log.Printf("CACHE HIT: %s\n", identifier)
		return item, nil
	}
	log.Printf("CACHE MISS: %s\n", identifier)

	metadata, err := GetItemMetadata(identifier)
	if err != nil {
		return nil, err
	}
	log.Println(metadata)

	root, err := getArchivePagesRoot(*metadata)
	if err != nil {
		return nil, err
	}

	item = &Item{
		Identifier:  identifier,
		LastUpdated: metadata.LastUpdated,
		Root:        root,
		D1:          metadata.D1,
		D2:          metadata.D2,
		Dir:         metadata.Dir,
	}
	SetCacheItem(identifier, item)
	return item, nil
}

func getArchivePagesRoot(metadata ItemMetadata) (string, error) {

	url := fmt.Sprintf("https://%s%s/archive.yml", metadata.D1, metadata.Dir)
	log.Println(url)
	var pagesMetadata PagesMetadata
	err := readYAML(url, &pagesMetadata)
	//err := Test(pagesMetadata)
	if err != nil {

	}
	return pagesMetadata.Root, nil
}

// GetItemMetadata returns the metadata of an archive.org item.
func GetItemMetadata(identifier string) (*ItemMetadata, error) {
	var metadata ItemMetadata
	url := "https://archive.org/metadata/" + identifier
	err := readJSON(url, &metadata)
	if err != nil {
		log.Fatalf("Failed to fetch metadata: %s", err)
		return nil, err
	}
	return &metadata, nil
}

func readJSON(url string, out interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	err = json.Unmarshal(body, out)
	if err != nil {
		return err
	}
	return nil
}

func readYAML(url string, out interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	log.Printf("%s\n", body)

	err = yaml.Unmarshal(body, out)
	if err != nil {
		return err
	}
	return nil
}

func Test(out interface{}) error {
	body := []byte("version: 1\nitemtype: website\nroot: www.aaronsw.com.zip")
	err := yaml.Unmarshal(body, out)
	if err != nil {
		return err
	}
	return nil
}
