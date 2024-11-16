package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

// Main entry point to run the provider

func main() {
	// This tells Terraform how to use your provider.
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: provider,
	})

}

// Provider function returns the schema for the provider

func provider() *schema.Provider {

	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"url": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("HTTP_SERVER_URL", "http://localhost:9999"),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"article": resourceArticle(),
		},
		ConfigureFunc: configure, // Set the configure function to read the API URL
	}
}

// Configure function that is used to initialize the provider with API URL

func configure(d *schema.ResourceData) (interface{}, error) {

	apiURL := d.Get("url").(string)
	if apiURL == "" {
		return nil, fmt.Errorf("API URL is required")
	}

	log.Printf("Configured API URL: %s\n", apiURL)
	return apiURL, nil // Return the URL as 'meta'

}

// Resource schema for 'article'

func resourceArticle() *schema.Resource {

	return &schema.Resource{
		Create: resourceArticleCreate,
		Read:   resourceArticleRead,
		Update: resourceArticleUpdate,
		Delete: resourceArticleDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"heading": {
				Type:     schema.TypeString,
				Required: true,
			},

			"description": {
				Type:     schema.TypeString,
				Required: true,
			},
			"tags": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString}, // Ensure Elem is defined this way
				Required: true,
			},
		},
	}

}

// Helper function to send the HTTP request

func sendRequest(url string, method string, data interface{}) (*http.Response, error) {

	client := &http.Client{}
	var req *http.Request
	var err error

	if data != nil {

		jsonData, _ := json.Marshal(data)
		req, err = http.NewRequest(method, url, bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

	} else {
		req, err = http.NewRequest(method, url, nil)
	}

	if err != nil {
		return nil, err

	}

	return client.Do(req)

}

// Resource creation function

func resourceArticleCreate(d *schema.ResourceData, meta interface{}) error {

	// Check if tags is set and is a list of strings
	tags := d.Get("tags")
	if tags == nil {

		return fmt.Errorf("tags must be a non-empty list of strings")

	}
	tagsList, ok := tags.([]interface{})
	if !ok {

		return fmt.Errorf("tags should be a list of strings, but got %T", tags)

	}

	// Convert []interface{} to []string

	stringTags := make([]string, len(tagsList))

	for i, tag := range tagsList {
		strTag, ok := tag.(string)
		if !ok {
			return fmt.Errorf("each tag must be a string, but got %T", tag)

		}

		stringTags[i] = strTag

	}

	// Log the tags after conversion

	log.Printf("Received tags: %+v\n", stringTags)

	// Check the type of 'meta' to ensure it's a string

	baseURL, ok := meta.(string)

	if !ok {

		return fmt.Errorf("expected meta to be a string, but got %T", meta)

	}

	// Log the URL before sending the request

	log.Printf("Sending request to URL: %s/api/v1/article\n", baseURL)

	// Now safely construct the URL

	url := fmt.Sprintf("%s/api/v1/article", baseURL)

	// Prepare the article data

	article := map[string]interface{}{

		"heading": d.Get("heading").(string),

		"description": d.Get("description").(string),

		"tags": stringTags,
	}

	// Log the article data before sending the request

	log.Printf("Sending article data: %+v\n", article)

	// Send the request

	resp, err := sendRequest(url, "POST", article)

	if err != nil {

		return err

	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {

		return fmt.Errorf("failed to create article: %s", resp.Status)

	}

	var createdArticle map[string]interface{}

	if err := json.NewDecoder(resp.Body).Decode(&createdArticle); err != nil {

		return err

	}

	// Set the article ID in the Terraform state

	d.SetId(fmt.Sprintf("%v", createdArticle["id"]))

	return resourceArticleRead(d, meta)

}

// Resource read function (example)

func resourceArticleRead(d *schema.ResourceData, meta interface{}) error {

	url := fmt.Sprintf("%s/api/v1/article?id=%s", meta.(string), d.Id())
	resp, err := sendRequest(url, "GET", nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to read article: %s", resp.Status)
	}

	var article map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&article); err != nil {
		return err
	}

	d.Set("heading", article["heading"])

	d.Set("description", article["description"])

	d.Set("tags", article["tags"])

	return nil

}

// Resource update function (example)

func resourceArticleUpdate(d *schema.ResourceData, meta interface{}) error {

	url := fmt.Sprintf("%s/api/v1/article?id=%s", meta.(string), d.Id())

	article := map[string]interface{}{

		"heading":     d.Get("heading").(string),
		"description": d.Get("description").(string),
		"tags":        d.Get("tags").([]string),
	}

	resp, err := sendRequest(url, "PUT", article)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to update article: %s", resp.Status)
	}

	return resourceArticleRead(d, meta)

}

// Resource delete function (example)

func resourceArticleDelete(d *schema.ResourceData, meta interface{}) error {

	url := fmt.Sprintf("%s/api/v1/article?id=%s", meta.(string), d.Id())
	resp, err := sendRequest(url, "DELETE", nil)

	if err != nil {
		return err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {

		return fmt.Errorf("failed to delete article: %s", resp.Status)

	}
	// Remove from the Terraform state
	d.SetId("")

	return nil

}
