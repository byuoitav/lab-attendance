package lab

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/byuoitav/common/log"
	"github.com/byuoitav/common/nerr"
	"github.com/byuoitav/lab-attendance/cache"
	"github.com/byuoitav/lab-attendance/messenger"
	"github.com/byuoitav/wso2services/wso2requests"
)

const labAPIURL = "https://commtech.byu.edu/noauth/ea/eaLabTrack.php"

type labRequest struct {
	LabID      string    `json:"lab_id"`
	Time       time.Time `json:"time"`
	BYUID      string    `json:"byu_id"`
	Action     string    `json:"action"`
	DeviceID   string    `json:"device_id"`
	DeviceType string    `json:"device_type"`
}

// Lab represents the configuration for and functions that can be run on a Lab
type Lab struct {
	M     *messenger.Messenger
	ID    string
	Cache *cache.Cache
}

// personsQueryResponse represents the response structure given by the Persons v3 API when a query is requested
type personsQueryResponse struct {
	Values []personsResponse `json:"values"`
}

// personsResponse represents the response structure given by the Persons v3 API
type personsResponse struct {
	Basic personsBasicFieldSet `json:"basic"`
}

// personsBasicFieldSet represents the "basic" field set returned by the Persons v3 API
type personsBasicFieldSet struct {
	Name      uapiField `json:"preferred_name"`
	NetID     uapiField `json:"net_id"`
	BYUID     uapiField `json:"byu_id"`
	FirstName uapiField `json:"preferred_first_name"`
}

// uapiField represents any generic field returned in a BYU University API response
type uapiField struct {
	Value       string `json:"value"`
	APIType     string `json:"api_type"`
	Description string `json:"description"`
}

// LogAttendanceForCard validates the cardID, translates it into a BYUID and then logs the user's attendance in the given lab
func (l Lab) LogAttendanceForCard(cardID string) error {

	p, err := l.Cache.GetPersonByCardID(cardID)
	if err != nil {
		log.L.Debugf("Cache miss for Card ID %s\n", cardID)
		// Call Persons v3 to validate the BYUID and get the name of the user
		q := personsQueryResponse{}
		err2, res, _ := wso2requests.MakeWSO2RequestReturnResponse("GET", fmt.Sprintf("https://api.byu.edu:443/byuapi/persons/v4/?credentials.credential_type=SEOS_CARD&credentials.credential_id=%s", cardID), nil, &q)
		if err2 != nil {

			if err2.Type == "request-error" && res.StatusCode == http.StatusNotFound {
				l.M.SendLoginErrorEvent("ID Card is not associated to a valid Identity")
				return fmt.Errorf("No matching identity found for Card ID %s", cardID)
			}

			err2 = nerr.Createf("Internal", "Error while attempting to validate the Card ID %s: %s", cardID, err2)
			log.L.Error(err2)

			l.M.SendLoginErrorEvent("We are unable to validate your login at this time. Please try again later.")

			return err2
		}

		if len(q.Values) < 1 {
			l.M.SendLoginErrorEvent("ID Card is not associated to a valid Identity")
			return fmt.Errorf("No matching identity found for Card ID %s", cardID)
		}

		p.BYUID = q.Values[0].Basic.BYUID.Value
		p.Name = q.Values[0].Basic.Name.Value
		p.NetID = q.Values[0].Basic.NetID.Value
		p.FirstName = q.Values[0].Basic.FirstName.Value
		p.CardID = cardID

		l.Cache.SavePersonToCache(p)
	}

	log.L.Debugf("Successfully validated Card ID %s: %s (%s)\n", cardID, p.Name, p.NetID)

	err = l.logAttendance(p.BYUID)
	if err != nil {

		// TODO: Theoretically any failure here should cause a cache, not an error, so an offline event should be sent
		// We need to validate a couple of cases for cache. What if we get an error back from the API? for non 500s?

		err = nerr.Createf("Internal", "Error while attemtping to log attendance to lab for BYU ID %s: %s", p.BYUID, err)
		log.L.Error(err)

		l.M.SendLoginErrorEvent("We are unable to log your attendance at this time. Please try again later.")

		return err
	}

	// Send successful event
	l.M.SendLoginEvent(p)

	return nil

}

// LogAttendanceForBYUID validates the BYUID and then logs the user's attendance in the given lab
func (l Lab) LogAttendanceForBYUID(byuID string) error {

	p, err := l.Cache.GetPersonByBYUID(byuID)
	if err != nil {
		log.L.Debugf("Cache miss for BYU ID %s\n", byuID)
		r := personsResponse{}

		// Call Persons v3 to validate the BYUID and get the name of the user
		err2, res, _ := wso2requests.MakeWSO2RequestReturnResponse("GET", fmt.Sprintf("https://api.byu.edu/byuapi/persons/v4/%s", byuID), nil, &r)
		if err2 != nil {

			if err2.Type == "request-error" && res.StatusCode == http.StatusNotFound {
				l.M.SendLoginErrorEvent(fmt.Sprintf("BYUID: %s is not a valid BYUID.", byuID))
				return fmt.Errorf("No matching identity found for BYUID %s", byuID)
			}

			err2 = nerr.Createf("Internal", "Error while attempting to validate the BYU ID %s: %s", byuID, err2)
			log.L.Error(err2)

			l.M.SendLoginErrorEvent("We are unable to validate your login at this time. Please try again later.")

			return err2
		}

		p.BYUID = r.Basic.BYUID.Value
		p.NetID = r.Basic.NetID.Value
		p.FirstName = r.Basic.FirstName.Value
		p.Name = r.Basic.Name.Value

		l.Cache.SavePersonToCache(p)
	}

	log.L.Debugf("Successfully validated BYU ID %s: %s (%s)\n", byuID, p.Name, p.NetID)

	err = l.logAttendance(byuID)
	if err != nil {

		// TODO: Theoretically any failure here should cause a cache, not an error, so an offline event should be sent
		// We need to validate a couple of cases for cache. What if we get an error back from the API? for non 500s?

		err = nerr.Createf("Internal", "Error while attemtping to log attendance to lab for BYU ID %s: %s", byuID, err)
		log.L.Error(err)

		l.M.SendLoginErrorEvent("We are unable to log your attendance at this time. Please try again later.")

		return err
	}

	// Send successful event
	l.M.SendLoginEvent(p)

	return nil
}

// LogAttendance logs the given BYUID's attendance to the current lab
func (l Lab) logAttendance(byuID string) error {

	req, err := json.Marshal(&labRequest{
		LabID:      l.ID,
		Time:       time.Now(),
		BYUID:      byuID,
		Action:     "enter",
		DeviceID:   "dev",
		DeviceType: "LAP",
	})
	if err != nil {
		err = nerr.Createf("Internal", "Failed to marshal request for attendance for BYUID %s: %s", byuID, err)
		log.L.Error(err)
		return err
	}

	res, err := http.Post(labAPIURL, "application/json", bytes.NewBuffer(req))
	if err != nil {
		// TODO: check for network errors and cache if the network is down
		err = nerr.Createf("Internal", "Failed while trying to make request for attendance for BYUID %s: %s", byuID, err)
		log.L.Error(err)
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		// TODO: What is the current requirement for a non-network related error? Should we cache here as well?
		err = nerr.Createf("Internal", "Got non-200 status code back from the lab attendance API: %+v", res)
		log.L.Error(err)
		return err
	}

	log.L.Debugf("Successfully logged attendence with the lab API! Response: %s", res.Body)

	// Successful call
	return nil

}
