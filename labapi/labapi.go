package labapi

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"github.com/byuoitav/common/log"
	"github.com/byuoitav/common/nerr"
)

const labAPIURL = "http://saasta.byu.edu/noauth/ea/eaLabTrack.php"

type labRequest struct {
	LabID      string    `json:"lab_id"`
	Time       time.Time `json:"time"`
	BYUID      string    `json:"byu_id"`
	Action     string    `json:"action"`
	DeviceID   string    `json:"device_id"`
	DeviceType string    `json:"device_type"`
}

// LogAttendance logs the given BYUID's attendance to the current lab
func LogAttendance(byuID string) error {

	req, err := json.Marshal(&labRequest{
		LabID:      "100",
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
	if res.StatusCode != 200 {
		// TODO: What is the current requirement for a non-network related error? Should we cache here as well?
		err = nerr.Createf("Internal", "Got non-200 status code back from the lab attendance API: %+v", res)
		log.L.Error(err)
		return err
	}

	// Successful call
	return nil

}
