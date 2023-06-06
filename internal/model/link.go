package model

type Link struct {
	Url string `json:"url"`
}

type LinkDB struct {
	OriginalUrl    string `db:"original_url"`
	AbbreviatedUrl string `db:"abbreviated_url"`
}
