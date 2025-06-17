package handlers

import (
	"github.com/byuoitav/common/v2/events"
	"github.com/byuoitav/lab-attendance/lab"
	"github.com/byuoitav/lab-attendance/messenger"
	"github.com/gin-gonic/gin"
)

// Login will check the given BYUID for validity and then call the Lab-Attendance API to log the user's attendance
func Login(m *messenger.Messenger, i events.BasicDeviceInfo, lab lab.Lab) func(*gin.Context) {

	return func(ctx *gin.Context) {

		byuID := ctx.Param("byuID")

		// TODO: Should we return a non 200 status code in any case?
		// 404 - BYU ID Invalid
		// 409 - Offline and no cache
		// 500 - Any non-cacheable error such as failure to marshal?
		lab.LogAttendanceForBYUID(byuID)

		ctx.Status(200)
	}
}
