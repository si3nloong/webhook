package dto

type Item struct {
	Item interface{} `json:"item"`
}

type Items struct {
	Items []interface{} `json:"items"`
	Count int64         `json:"count"`
	Size  int           `json:"size"`
	Links struct {
		Self     string `json:"self"`
		Previous string `json:"previous"`
		Next     string `json:"next"`
	} `json:"_links"`
}
