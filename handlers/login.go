package handlers

import (
	"fmt"
	"time"

	"github.com/byuoitav/central-event-system/messenger"
	"github.com/byuoitav/common/log"
	"github.com/byuoitav/common/nerr"
	"github.com/byuoitav/common/v2/events"
	"github.com/byuoitav/lab-attendance/labapi"
	"github.com/byuoitav/wso2services/wso2requests"
	"github.com/labstack/echo"
)

// PersonsResponse represents the response structure given by the Persons v3 API
type PersonsResponse struct {
	Basic PersonsBasicFieldSet `json:"basic"`
}

// PersonsBasicFieldSet represents the "basic" field set returned by the Persons v3 API
type PersonsBasicFieldSet struct {
	Name  UAPIField `json:"preferred_name"`
	NetID UAPIField `json:"net_id"`
}

// UAPIField represents any generic field returned in a BYU University API response
type UAPIField struct {
	Value       string `json:"value"`
	APIType     string `json:"api_type"`
	Description string `json:"description"`
}

// Login will check the given BYUID for validity and then call the Lab-Attendan/ce API to log the user's attendance
func Login(m *messenger.Messenger, i events.BasicDeviceInfo) func(echo.Context) error {

	return func(ctx echo.Context) error {

		byuID := ctx.Param("byuID")
		res := PersonsResponse{}

		event := events.Event{
			GeneratingSystem: i.DeviceID,
			Timestamp:        time.Now(),
			TargetDevice:     i,
			AffectedRoom:     i.BasicRoomInfo,
			Key:              "Login",
			Value:            "False",
		}

		// Call Persons v3 to validate the BYUID and get the name of the user
		err := wso2requests.MakeWSO2Request("GET", fmt.Sprintf("https://api.byu.edu/byuapi/persons/v3/%s", byuID), nil, &res)
		if err != nil {

			// TODO: Check for network errors and cache the login attempt

			err = nerr.Createf("Internal", "Error while attempting to validate the BYU ID %s: %s", byuID, err)
			log.L.Error(err)
			m.SendEvent(event)
			return err
		}

		log.L.Debugf("Successfully validated BYU ID %s: %s (%s)\n", byuID, res.Basic.Name.Value, res.Basic.NetID.Value)

		err2 := labapi.LogAttendance(byuID)
		if err2 != nil {
			err = nerr.Createf("Internal", "Error while attemtping to log attendance to lab for BYU ID %s: %s", byuID, err)
			log.L.Error(err)
			m.SendEvent(event)
			return err
		}

		// Set the login event to true and add extra information
		event.Value = "True"
		event.User = res.Basic.NetID.Value
		event.Data = res.Basic.Name.Value

		log.L.Debugf("Sending event: %+v\n", event)
		m.SendEvent(event)

		return nil
	}
}
