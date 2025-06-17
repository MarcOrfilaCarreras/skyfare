package vueling

const (
	TICKETS_BASE_URL = "https://tickets.vueling.com"
	AMS_BASE_URL     = "https://ams.vueling.com"

	TICKETS_SERVICE_BOOKING_URL             = TICKETS_BASE_URL + "/booking"
	TICKETS_SERVICE_BOOKING_FLIGHT_SEARCH_URL = TICKETS_SERVICE_BOOKING_URL + "/flightSearch"
	TICKETS_SERVICE_ASSETS_URL              = TICKETS_BASE_URL + "/assets"
	TICKETS_SERVICE_ASSETS_STATIONS_FILE    = "es-ES.json"
	TICKETS_SERVICE_ASSETS_STATIONS_URL     = TICKETS_SERVICE_ASSETS_URL + "/stations/" + TICKETS_SERVICE_ASSETS_STATIONS_FILE

	AMS_SERVICE_VERSION                     = "v1"
	AMS_SERVICE_ASM_BASE_URL                = AMS_BASE_URL + "/asm/" + AMS_SERVICE_VERSION
	AMS_SERVICE_ASM_AUTH_URL                = AMS_SERVICE_ASM_BASE_URL + "/Auth"

	AMS_SERVICE_RES_BASE_URL                = AMS_BASE_URL + "/res/" + AMS_SERVICE_VERSION
	AMS_SERVICE_RES_MARKETS_URL             = AMS_SERVICE_RES_BASE_URL + "/Markets"
	AMS_SERVICE_RES_MARKETS_BYORIGIN_URL    = AMS_SERVICE_RES_MARKETS_URL + "/ByOrigin"

	AMS_SERVICE_AVY_VERSION                 = "v2"
	AMS_SERVICE_AVY_BASE_URL                = AMS_BASE_URL + "/avy/" + AMS_SERVICE_AVY_VERSION
	AMS_SERVICE_AVY_AVAILABILITY_URL        = AMS_SERVICE_AVY_BASE_URL + "/AvailabilityServices"
	AMS_SERVICE_AVY_AVAILABILITY_FLIGHTS_URL = AMS_SERVICE_AVY_AVAILABILITY_URL + "/flightsSummary"
)

var HEADERS = map[string]string{
	"accept":            "application/json, text/plain, */*",
	"accept-language":   "es-ES,es;q=0.7",
	"content-type":      "application/json",
	"origin":            TICKETS_BASE_URL,
	"referer":           TICKETS_BASE_URL,
	"user-agent":        "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36",
	"sec-ch-ua":         `"Chromium";v="124", "Brave";v="124", "Not-A.Brand";v="99"`,
	"sec-ch-ua-mobile":  "?0",
	"sec-ch-ua-platform": `"Linux"`,
	"sec-fetch-dest":    "empty",
	"sec-fetch-mode":    "cors",
	"sec-fetch-site":    "same-site",
	"sec-gpc":           "1",
	"priority":          "u=1, i",
}
