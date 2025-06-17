package vueling

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
	"database/sql"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-resty/resty/v2"
    "github.com/MarcOrfilaCarreras/skyfare/internal/model"
	"github.com/MarcOrfilaCarreras/skyfare/internal/repository"
)

type VuelingScraper struct {
	client       *resty.Client
	profileID    string
	accessToken  string
	airportRepo  *repository.AirportRepository
	flightRepo  *repository.FlightsRepository
}

func NewVuelingScraper(db *sql.DB) *VuelingScraper {
	client := resty.New().
		SetTimeout(15 * time.Second).
		SetRetryCount(3).
		SetHeader("User-Agent", "Mozilla/5.0 (compatible)").
		SetHeaders(HEADERS)

	vs := &VuelingScraper{
		client: client,
		airportRepo: repository.NewAirportRepository(db),
		flightRepo: repository.NewFlightsRepository(db),
	}
	if err := vs.init(); err != nil {
		fmt.Printf("Failed to initialize Vueling scraper: %v\n", err)
	}
	return vs
}

func (v *VuelingScraper) init() error {
	profileID, err := v.getProfileID()
	if err != nil {
		return err
	}
	v.profileID = profileID

	token, err := v.getAccessToken()
	if err != nil {
		return err
	}
	v.accessToken = token

	return nil
}

func (v *VuelingScraper) GetAirports() ([]model.Airport, error) {
	resp, err := v.client.R().
		Get(TICKETS_SERVICE_ASSETS_STATIONS_URL)
	if err != nil {
		return nil, err
	}

	if !isJSONResponse(resp) {
		return nil, errors.New("invalid content type")
	}

	var raw []struct {
		StationCode    string `json:"stationCode"`
		FullName       string `json:"fullName"`
		LocationDetails struct {
			CountryCode string `json:"countryCode"`
			Coordinates struct {
				Lat string `json:"latitude"`
				Lng string `json:"longitude"`
			} `json:"coordinates"`
		} `json:"locationDetails"`
	}

	if err := json.Unmarshal(resp.Body(), &raw); err != nil {
		return nil, err
	}

	var airports []model.Airport
	for _, item := range raw {
		airport := model.Airport{
			Code:    item.StationCode,
			Name:    item.FullName,
			Country: item.LocationDetails.CountryCode,
			Coordinates: model.Coordinates{
				Lat: item.LocationDetails.Coordinates.Lat,
				Lng: item.LocationDetails.Coordinates.Lng,
			},
		}
		airports = append(airports, airport)

		if err := v.airportRepo.InsertAirport(airport); err != nil {
			fmt.Printf("Failed to insert airport %s: %v\n", airport.Code, err)
		}
	}
	return airports, nil
}

func (v *VuelingScraper) GetAirportRoutes(code string) ([]model.Route, error) {
	url := fmt.Sprintf("%s/%s", AMS_SERVICE_RES_MARKETS_BYORIGIN_URL, code)

	resp, err := v.client.R().
		SetHeader("Authorization", "Bearer "+v.accessToken).
		Get(url)
	if err != nil {
		return nil, err
	}

	if !isJSONResponse(resp) {
		return nil, errors.New("invalid content type")
	}

	var raw []struct {
		ToCode    string `json:"toCode"`
		Connection string `json:"connection"`
	}

	if err := json.Unmarshal(resp.Body(), &raw); err != nil {
		return nil, err
	}

	var routes []model.Route
	for _, item := range raw {
		routes = append(routes, model.Route{
			Code: item.ToCode,
			Connection: item.Connection,
		})
	}
	return routes, nil
}

func (v *VuelingScraper) GetRoute(codeFrom string, codeTo string, month int, year int, currency string) ([]model.Flight, error) {
	const layout = "2006-01-02T15:04:05"

	body := map[string]any{
		"originCode":      codeFrom,
		"destinationCode": codeTo,
		"year":            year,
		"month":           month,
		"currencyCode":    currency,
		"monthsRange":     9,
	}

	resp, err := v.client.R().
		SetHeader("Authorization", "Bearer "+v.accessToken).
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		Post(AMS_SERVICE_AVY_AVAILABILITY_FLIGHTS_URL)
	if err != nil {
		return nil, err
	}

	if !isJSONResponse(resp) {
		return nil, errors.New("invalid content type")
	}

	var raw []struct {
		Date string   `json:"departureDate"`
		Price         *float64 `json:"price,omitempty"`
		Promotion  bool     `json:"promotion"`
	}

	if err := json.Unmarshal(resp.Body(), &raw); err != nil {
		return nil, err
	}

	var flights []model.Flight
	for _, item := range raw {
		price := 0.0
		if item.Price != nil {
			price = *item.Price
		}

		parsedDate, err := time.Parse(layout, item.Date)
		if err != nil {
			return nil, fmt.Errorf("invalid date format: %w", err)
		}

		flight := model.Flight{
			Date:			parsedDate,
			Price:			price,
			Promotion:		item.Promotion,
		}

		if err := v.flightRepo.InsertFlight(codeFrom, codeTo, &flight); err != nil {
			fmt.Printf("failed to insert flight: %v\n", err)
		}

		flights = append(flights, flight)
	}
	return flights, nil
}

func (v *VuelingScraper) getProfileID() (string, error) {
	bookingURL := TICKETS_SERVICE_BOOKING_URL

	resp, err := v.client.R().Get(bookingURL)
	if err != nil {
		return "", err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(resp.String()))
	if err != nil {
		return "", err
	}

	var chunkHref string
	doc.Find("link[href]").EachWithBreak(func(i int, s *goquery.Selection) bool {
		href, exists := s.Attr("href")
		if exists && chunkFileRegexp.MatchString(href) {
			chunkHref = href
			return false
		}
		return true
	})

	if chunkHref == "" {
		return "", errors.New("chunk file not found")
	}

	chunkResp, err := v.client.R().Get(TICKETS_BASE_URL + "/" + chunkHref)
	if err != nil {
		return "", err
	}

	match := profileIDRegexp.FindStringSubmatch(chunkResp.String())
	if len(match) < 2 {
		return "", errors.New("profile ID not found")
	}

	return match[1], nil
}

func (v *VuelingScraper) getAccessToken() (string, error) {
	data := map[string]string{"profileId": v.profileID}

	resp, err := v.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(data).
		Post(AMS_SERVICE_ASM_AUTH_URL)
	if err != nil {
		return "", err
	}

	var result struct {
		TokenType      string `json:"tokenType"`
		AccessToken    string `json:"accessToken"`
		Expiration     int    `json:"expiration"`
		UserType       string `json:"userType"`
		HasActiveSession bool `json:"hasActiveSession"`
	}

	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return "", err
	}

	if result.AccessToken == "" {
		return "", errors.New("access token missing in response")
	}
	return result.AccessToken, nil
}

func isJSONResponse(resp *resty.Response) bool {
	return strings.Contains(resp.Header().Get("Content-Type"), "application/json")
}

var (
	chunkFileRegexp  = regexp.MustCompile(`chunk-[A-Z0-9]+\.js`)
	profileIDRegexp  = regexp.MustCompile(`profileId:"([a-f0-9-]{36})"`)
)
