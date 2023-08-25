package models

type OrderPrimaryKey struct {
	Id string `json:"id"`
}

type Order struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Lat     string `json:"lat"`
	Long    string `json:"long"`
	Address string `json:"address"`
	Photo   string `json:"photo"`
}

type OrderCreate struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Lat     string `json:"lat"`
	Long    string `json:"long"`
	Address string `json:"address"`
	Photo   string `json:"photo"`
}

type OrderUpdate struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Lat     string `json:"lat"`
	Long    string `json:"long"`
	Address string `json:"address"`
	Photo   string `json:"photo"`
}
