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
	"time"

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

		if len(dateSubmit) > 0 {
			// Votre date en string
			dateString := dateSubmit[0]

			// Supprimer le "Le " du début
			dateString = dateString[3:]

			// Définir le layout pour correspondre à votre date
			layout := "02 jan. 2006"

			// Utiliser time.Parse pour convertir la string en time.Time
			t, err := time.Parse(layout, dateString)
			if err != nil {
				fmt.Println(err)
			}

			// Obtenir la date actuelle
			now := time.Now()

			// Comparer les deux dates
			if t.After(now) {
				fmt.Println("La date est dans le futur")
			} else if t.Before(now) {
				fmt.Println("La date est dans le passé")
			} else {
				fmt.Println("La date est aujourd'hui")
			}
		}
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
