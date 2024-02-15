package models

type LeroymerlinURL struct {
    QuestionURL string `json:"question_url"`
    ArticleURL  string `json:"article_url"`
}

type Config struct {
	Headers map[string]string `json:"headers"`
}