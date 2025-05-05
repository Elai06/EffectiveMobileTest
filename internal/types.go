package internal

type Person struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Surname     string `json:"surname"`
	Age         int    `json:"age"`
	Gender      string `json:"gender"`
	Nationality string `json:"nationality"`
}

type PersonFilter struct {
	Person Person
	AgeMax int `json:"age_max"`
	AgeMin int `json:"age_min"`
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

type Country struct {
	Country     string  `json:"country_id"`
	Probability float64 `json:"probability"`
}
