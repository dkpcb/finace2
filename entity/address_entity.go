package entity

type Address struct {
	PostalCode       string  `json:"postal_code"`
	HitCount         int     `json:"hit_count"`
	CommonAddress    string  `json:"address"`
	TokyoStaDistance float64 `json:"tokyo_sta_distance"`
}

func NewAddress(postalCode string, hitCount int, commonAddress string, tokyoStaDistance float64) *Address {
	return &Address{
		PostalCode:       postalCode,
		HitCount:         hitCount,
		CommonAddress:    commonAddress,
		TokyoStaDistance: tokyoStaDistance,
	}
}
