package model

type Place struct {
	Name    string `json:"Name"`
	Address string `json:"Address"`
	Phone   string `json:"Phone"`

	Location struct {
		Lat float64 `json:"Lat"`
		Lon float64 `json:"Long"`
	} `json:"location"`
}

type Api struct {
	Name      string  `json:"name"`
	Total     int     `json:"total"`
	Places    []Place `json:"places"`
	Prev_page int     `json:"prev_page"`
	Next_page int     `json:"next_page"`
	Last_page int     `json:"last_page"`
}
