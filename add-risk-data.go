package main

/*
This API will update ticket payload with weather
risk and fruad risk data

*/
import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

const (
	//FORMS
	STORM_FORM_ID = 114093996871
	TV_FORM_ID    = 114093998312
	//FIELDS
	WIND_SPEED_FIELD_ID   = 114100596852
	TV_MODEL_FIELD_ID     = 114099896612
	CLAIM_TYPE_FIELD_ID   = 114099964311 // TV , Storm Surge
	WEATHER_RISK_FIELD_ID = 114100712171
	FRAUD_RISK_FIELD_ID   = 114100657672
	PHONE_FIELD_ID        = 114100658992
	EMAIL_FIELD_ID        = 114100659172
)

func Handler(w http.ResponseWriter, r *http.Request) {

	var composedResult = ComposedResult{}
	err := json.NewDecoder(r.Body).Decode(&composedResult)
	if err != nil {
		createErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	// populate ticket payload with custom field
	ticketDetails := TicketDetails{}
	var claimType, riskDesc, windSpeed, tvModel string
	if composedResult.TranformedData.Status == 200 {
		ticketDetails = composedResult.TranformedData.TicketDetails
		if strings.Contains(ticketDetails.Ticket.Subject, "Storm") {
			claimType = "Storm Surge"
			ticketDetails.Ticket.TicketFormID = STORM_FORM_ID
		} else {
			claimType = "TV"
			ticketDetails.Ticket.TicketFormID = TV_FORM_ID
		}
	}

	if composedResult.WeatherRisk.Status == 200 {
		weatherRisk := composedResult.WeatherRisk
		riskDesc = strconv.Itoa(weatherRisk.RiskScore) + " : " + weatherRisk.Description
	}
	if composedResult.WeatherData.Status == 200 {
		weatherData := composedResult.WeatherData
		windSpeed = weatherData.History.Dailysummary[0].Maxwspdm
	}
	if (TVClaimData{}) != composedResult.TranformedData.TVClaimData {
		tvModel = composedResult.TranformedData.TVClaimData.TVModelNo
	}

	customFields := []CustomFields{
		CustomFields{WIND_SPEED_FIELD_ID, windSpeed},
		CustomFields{CLAIM_TYPE_FIELD_ID, claimType},
		CustomFields{WEATHER_RISK_FIELD_ID, riskDesc},
		CustomFields{PHONE_FIELD_ID, ticketDetails.Ticket.Requester.Phone},
		CustomFields{EMAIL_FIELD_ID, ticketDetails.Ticket.Requester.Email},
		CustomFields{TV_MODEL_FIELD_ID, tvModel},
	}

	ticketDetails.Ticket.CustomFields = customFields
	ticketDetails.Status = 200

	//marshal to JSON
	ticketDetailsJSON, err := json.Marshal(ticketDetails)
	if err != nil {
		createErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.Write([]byte(string(ticketDetailsJSON)))

}

// func main() {
// 	println("staritng app..")
// 	http.HandleFunc("/", Handler)
// 	http.ListenAndServe(":8084", nil)
// }

func createErrorResponse(w http.ResponseWriter, message string, status int) {
	errorJSON, _ := json.Marshal(&Error{
		Status:  status,
		Message: message})
	//Send custom error message to caller
	w.WriteHeader(status)
	w.Header().Set("content-type", "application/json")
	w.Write([]byte(errorJSON))
}

type Error struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type ComposedResult struct {
	TranformedData TranformedData `json:"tranformed-data"`
	WeatherData    WeatherData    `json:"weather-data"`
	WeatherRisk    struct {
		Status      int    `json:"status"`
		Description string `json:"description"`
		RiskScore   int    `json:"riskScore"`
	} `json:"weather-risk"`
}

type WeatherData struct {
	Status  int `json:"status"`
	History struct {
		Dailysummary []struct {
			Fog          string `json:"fog"`
			Maxpressurem string `json:"maxpressurem"`
			Maxtempm     string `json:"maxtempm"`
			Maxwspdm     string `json:"maxwspdm"`
			Minpressurem string `json:"minpressurem"`
			Mintempm     string `json:"mintempm"`
			Minwspdm     string `json:"minwspdm"`
			Rain         string `json:"rain"`
			Tornado      string `json:"tornado"`
		} `json:"dailysummary"`
	} `json:"history"`
	Response struct {
		Version string `json:"version"`
	} `json:"response"`
}
type CustomFields struct {
	ID    int64  `json:"id"`
	Value string `json:"value"`
}

type TranformedData struct {
	Status          int             `json:"status,omitempty"`
	TicketDetails   TicketDetails   `json:"ticket_details,omitempty"`
	WeatherAPIInput WeatherAPIInput `json:"weather_api_input,omitempty"`
	TVClaimData     TVClaimData     `json:"tv_claim_data,omitempty"`
	StromClaimData  StromClaimData  `json:"storm_claim_data,omitempty"`
}

type TVClaimData struct {
	TVPrice         string `json:"tv_price,omitempty"`
	CrimeRef        string `json:"crime_ref,omitempty"`
	IncidentDate    string `json:"incident_date,omitempty"`
	TVModelNo       string `json:"tv_model_no,omitempty"`
	TVMake          string `json:"tv_make,omitempty"`
	TVSerialNo      string `json:"tv_serial_no,omitempty"`
	DamageImageURL1 string `json:"damage_image_url_1,omitempty"`
	DamageImageURL2 string `json:"damage_image_url_2,omitempty"`
	TVReceiptImage  string `json:"tv_reciept_image_url"`
	DamageVideoURL  string `json:"damage_video_url,omitempty"`
}

type StromClaimData struct {
	IncidentPlace       string `json:"incident_place,omitempty"`
	IncidentDate        string `json:"incident_date,omitempty"`
	DamageImageURL1     string `json:"damage_image_url_1,omitempty"`
	DamageImageURL2     string `json:"damage_image_url_2,omitempty"`
	RepairEstimateImage string `json:"estimate_image_url,omitempty"`
	DamageVideoURL      string `json:"damage_video_url,omitempty"`
}

type TicketDetails struct {
	Status int `json:"status"`
	Ticket struct {
		Type     string `json:"type"`
		Subject  string `json:"subject"`
		Priority string `json:"priority"`
		Status   string `json:"status"`
		Comment  struct {
			HTMLBody string   `json:"html_body"`
			Uploads  []string `json:"uploads,omitempty"`
		} `json:"comment"`
		CustomFields []CustomFields `json:"custom_fields,omitempty"`
		Requester    struct {
			LocaleID     int    `json:"locale_id"`
			Name         string `json:"name"`
			Email        string `json:"email"`
			Phone        string `json:"phone"`
			PolicyNumber string `json:"policy_number"`
		} `json:"requester"`
		TicketFormID int64 `json:"ticket_form_id"`
	} `json:"ticket"`
}

type WeatherAPIInput struct {
	City    string `json:"city,omitempty"`
	Country string `json:"country,omitempty"`
	Date    string `json:"date,omitempty"` //YYYYMMDD
}
