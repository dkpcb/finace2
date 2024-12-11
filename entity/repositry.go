package entity

type AddressRepository interface {
	FetchAddressData(postalCode string) ([]AddressLocation, error)
}

type AddressLocation struct {
	Prefecture string  `json:"prefecture"`
	City       string  `json:"city"`
	Town       string  `json:"town"`
	Lat        float64 `json:"y,string"`
	Lon        float64 `json:"x,string"`
}
