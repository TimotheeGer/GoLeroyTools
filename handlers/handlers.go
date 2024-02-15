package handlers

import (
	"GoLeroyTools/leroyTools/models"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"golang.org/x/net/html"
)

func RetrieveURLsFromJSON() ([]models.LeroymerlinURL, error) {

	fmt.Println("Go RetrieveURLsFromJSON!")

	data, err := os.ReadFile("url.json")
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}

	var urls []models.LeroymerlinURL
	err = json.Unmarshal(data, &urls)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}

	return urls, nil
}

func MakeRequest(url models.LeroymerlinURL) (*http.Response, string, error) {

	fmt.Println("Go MakeRequest!")
	file, _ := os.Open("config.json")
	defer file.Close()

	decoder := json.NewDecoder(file)
	config := models.Config{}
	err := decoder.Decode(&config)
	if err != nil {
		return nil, "", err
	}

	req, err := http.NewRequest("GET", url.QuestionURL, nil)
	if err != nil {
		return nil, "", err
	}
	for key, value := range config.Headers {
		req.Header.Add(key, value)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	return resp, url.ArticleURL, nil
}

func ProcessURLs(urls []models.LeroymerlinURL) {
	for _, url := range urls {
		resp, article, err := MakeRequest(url)
		if err != nil {
			log.Println("Erreur lors de la requête à l'URL", url, ":", err)
			continue
		}
		defer resp.Body.Close()

		fmt.Printf("response status: %v, for this item: %v\n", resp.Status, article)

		var reader io.ReadCloser
		switch resp.Header.Get("Content-Encoding") {
		case "gzip":
			reader, err = gzip.NewReader(resp.Body)
			if err != nil {
				log.Fatalln(err)
			}
			defer reader.Close()
		default:
			reader = resp.Body
		}

		doc, err := html.Parse(reader)
		if err != nil {
			log.Fatalln(err)
		}

		dateSubmit := FindSubmissionDates(doc)
		fmt.Println("dateSubmit: ", dateSubmit)
	}
}

func FindSubmissionDates(n *html.Node) []string {
	var dates []string

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "p" {
			for _, a := range n.Attr {
				if a.Key == "class" && a.Val == "a-user__submissionDate" {
					if n.FirstChild != nil {
						dates = append(dates, n.FirstChild.Data)
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(n)

	return dates
}
