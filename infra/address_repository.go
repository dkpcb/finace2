package infra

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dkpcb/finatext_kadai_2/entity"
)

type AddressRepository struct {
	ExternalAPI string
}

// 新しい AddressRepository を作成
func NewAddressRepository(externalAPI string) *AddressRepository {
	return &AddressRepository{ExternalAPI: externalAPI}
}

// 外部APIを呼び出して住所データを取得
func (r *AddressRepository) FetchAddressData(postalCode string) ([]entity.AddressLocation, error) {
	// 外部APIのURLを組み立て
	url := fmt.Sprintf("%s%s", r.ExternalAPI, postalCode)

	// HTTPリクエストを送信
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to call external API: %w", err)
	}
	defer resp.Body.Close()

	// HTTPステータスコードをチェック
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// レスポンスをデコード
	var result struct {
		Response struct {
			Location []entity.AddressLocation `json:"location"`
		} `json:"response"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode API response: %w", err)
	}

	return result.Response.Location, nil
}

//aaaaaa

// package infra

// import (
// 	"encoding/json"
// 	"fmt"
// 	"net/http"

// 	"github.com/dkpcb/finatext_kadai_2/entity"
// 	"github.com/dkpcb/finatext_kadai_2/repository"
// )

// type AddressRepositoryImpl struct {
// 	APIEndpoint string
// }

// func NewAddressRepository(apiEndpoint string) repository.AddressRepository {
// 	return &AddressRepositoryImpl{APIEndpoint: apiEndpoint}
// }

// func (r *AddressRepositoryImpl) FetchAddressData(postalCode string) ([]entity.AddressLocation, error) {
// 	url := fmt.Sprintf("%s%s", r.APIEndpoint, postalCode)

// 	resp, err := http.Get(url)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to fetch address data: %w", err)
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusOK {
// 		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
// 	}

// 	var response struct {
// 		Response struct {
// 			Location []entity.AddressLocation `json:"location"`
// 		} `json:"response"`
// 	}
// 	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
// 		return nil, fmt.Errorf("failed to decode address data: %w", err)
// 	}

// 	return response.Response.Location, nil
// }

// package infra

// import (
// 	"encoding/json"
// 	"fmt"
// 	"net/http"

// 	"github.com/dkpcb/finatext_kadai_2/domain"
// )

// type ExternalAPIResponse struct {
// 	Response struct {
// 		Location []struct {
// 			Prefecture string  `json:"prefecture"`
// 			City       string  `json:"city"`
// 			Town       string  `json:"town"`
// 			X          float64 `json:"x,string"`
// 			Y          float64 `json:"y,string"`
// 		} `json:"location"`
// 	} `json:"response"`
// }

// type ExternalAPIClient struct {
// 	BaseURL string
// }

// func NewExternalAPIClient(baseURL string) *ExternalAPIClient {
// 	return &ExternalAPIClient{BaseURL: baseURL}
// }

// func (client *ExternalAPIClient) FetchAddressData(postalCode string) ([]domain.AddressLocation, error) {
// 	url := fmt.Sprintf("%s?method=searchByPostal&postal=%s", client.BaseURL, postalCode)

// 	resp, err := http.Get(url)
// 	if err != nil || resp.StatusCode != http.StatusOK {
// 		return nil, fmt.Errorf("failed to fetch address data")
// 	}
// 	defer resp.Body.Close()

// 	var apiResponse ExternalAPIResponse
// 	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
// 		return nil, fmt.Errorf("failed to parse API response")
// 	}

// 	locations := make([]domain.AddressLocation, len(apiResponse.Response.Location))
// 	for i, loc := range apiResponse.Response.Location {
// 		locations[i] = domain.AddressLocation{
// 			Prefecture: loc.Prefecture,
// 			City:       loc.City,
// 			Town:       loc.Town,
// 			Lat:        loc.Y,
// 			Lon:        loc.X,
// 		}
// 	}
// 	return locations, nil
// }
