package dto

type Item struct {
	Item interface{} `json:"item"`
}

type Items struct {
	Items []interface{} `json:"items"`
	Size  int           `json:"size"`
	Links struct {
		Self     string `json:"self"`
		Previous string `json:"previous"`
		Next     string `json:"next"`
	} `json:"_links"`
}
