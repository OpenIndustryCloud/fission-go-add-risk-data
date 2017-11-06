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

	if r.Body == nil {
		http.Error(w, "Please send a valid JSON", 400)
		return
	}
	var composedResult ComposedResult
	err := json.NewDecoder(r.Body).Decode(&composedResult)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	// populate ticket payload with custom field
	ticketDetails := composedResult.TranformedData.TicketDetails
	var claimType = ""
	if strings.Contains(ticketDetails.Ticket.Subject, "Storm") {
		claimType = "Storm Surge"
		ticketDetails.Ticket.TicketFormID = STORM_FORM_ID
	} else {
		claimType = "TV"
		ticketDetails.Ticket.TicketFormID = TV_FORM_ID
	}
	weatherRisk := composedResult.WeatherRisk
	weatherData := composedResult.WeatherData
	riskDesc := strconv.Itoa(weatherRisk.RiskScore) + " : " + weatherRisk.Description

	customFields := []CustomFields{
		CustomFields{WIND_SPEED_FIELD_ID, weatherData.History.Dailysummary[0].Maxwspdm},
		CustomFields{CLAIM_TYPE_FIELD_ID, claimType},
		CustomFields{WEATHER_RISK_FIELD_ID, riskDesc},
		//CustomFields{FRAUD_RISK_FIELD_ID, ""},
		CustomFields{PHONE_FIELD_ID, ticketDetails.Ticket.Phone},
		CustomFields{EMAIL_FIELD_ID, ticketDetails.Ticket.Email}}

	ticketDetails.Ticket.CustomFields = customFields

	//marshal to JSON
	ticketDetailsJSON, err := json.Marshal(ticketDetails)
	if err != nil {
		http.Error(w, err.Error(), 400)
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

type ComposedResult struct {
	TranformedData TranformedData `json:"tranformedData"`
	WeatherData    WeatherData    `json:"weatherData"`
	WeatherRisk    struct {
		Description string `json:"description"`
		RiskScore   int    `json:"riskScore"`
	} `json:"weatherRisk"`
}

type WeatherData struct {
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
	TicketDetails struct {
		Ticket struct {
			Comment struct {
				HTMLBody string `json:"html_body"`
			} `json:"comment"`
			CustomFields []CustomFields `json:"custom_fields"`
			Requester    struct {
				LocaleID int    `json:"locale_id"`
				Name     string `json:"name"`
				Email    string `json:"email"`
			} `json:"requester"`
			Email        string `json:"email"`
			Phone        string `json:"phone"`
			Priority     string `json:"priority"`
			Status       string `json:"status"`
			Subject      string `json:"subject"`
			Type         string `json:"type"`
			TicketFormID int64  `json:"ticket_form_id"`
		} `json:"ticket"`
	} `json:"ticketDetails"`
	WeatherAPIInput struct {
		City    string `json:"city"`
		Country string `json:"country"`
		Date    string `json:"date"`
	} `json:"weatherAPIInput"`
}
