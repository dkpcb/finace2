package service

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"unicode/utf8"

	"github.com/dkpcb/finatext_kadai_2/entity"
	"github.com/dkpcb/finatext_kadai_2/util"
)

type AddressService struct {
	Repo        entity.AddressRepository
	ExternalAPI string
}

func NewAddressService(repo entity.AddressRepository, externalAPI string) *AddressService {
	return &AddressService{
		Repo:        repo,
		ExternalAPI: externalAPI,
	}
}

func (s *AddressService) GetAddress(postalCode string) (*entity.Address, error) {
	log.Printf("Starting GetAddress for postalCode: %s", postalCode)

	// 外部 API からデータを取得
	locations, err := s.Repo.FetchAddressData(postalCode)
	if err != nil {
		log.Printf("Failed to fetch address data for postalCode: %s, error: %v", postalCode, err)
		return nil, err
	}

	// データが空の場合
	if len(locations) == 0 {
		log.Printf("No address data found for postalCode: %s", postalCode)
		return nil, nil
	}

	log.Printf("Fetched %d locations for postalCode: %s", len(locations), postalCode)

	// 共通の住所を組み立てる
	commonPrefecture := locations[0].Prefecture
	commonCity := locations[0].City
	commonTown := extractCommonTown(locations)
	commonAddress := commonPrefecture + commonCity + commonTown

	log.Printf("Constructed common address: %s", commonAddress)

	// 東京駅からの最大距離を計算
	var maxDistance float64
	for _, loc := range locations {
		distance := util.CalculateDistance(util.TokyoStationLat, util.TokyoStationLon, loc.Lat, loc.Lon)
		log.Printf("Calculated distance for location (%f, %f): %f km", loc.Lat, loc.Lon, distance)

		if distance > maxDistance {
			maxDistance = distance
			log.Printf("Updated max distance: %f km", maxDistance)
		}
	}

	log.Printf("Final max distance for postalCode %s: %f km", postalCode, maxDistance)

	address := entity.NewAddress(postalCode, len(locations), commonAddress, maxDistance)
	log.Printf("Constructed address entity: %+v", address)

	return address, nil
}

func (s *AddressService) fetchAddressFromAPI(postalCode string) ([]entity.AddressLocation, error) {
	url := fmt.Sprintf("%s%s", s.ExternalAPI, postalCode)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to call external API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

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

// 共通の町名を抽出
func extractCommonTown(locations []entity.AddressLocation) string {
	if len(locations) == 0 {
		return ""
	}

	// 最初の町名を基準に共通部分を抽出
	commonTown := locations[0].Town
	for _, loc := range locations {
		commonTown = commonPrefix(commonTown, loc.Town)
		if commonTown == "" {
			break
		}
	}

	// 共通部分をUTF-8として正しく処理
	if !utf8.ValidString(commonTown) {
		return strings.TrimRight(commonTown, "\uFFFD") // 不正文字（�）を取り除く
	}

	return strings.TrimSpace(commonTown)
}

// 2つの文字列の共通部分を抽出
func commonPrefix(a, b string) string {
	minLength := len(a)
	if len(b) < minLength {
		minLength = len(b)
	}

	for i := 0; i < minLength; i++ {
		if a[i] != b[i] {
			return a[:i]
		}
	}
	return a[:minLength]
}
