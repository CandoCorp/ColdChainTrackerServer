package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func init() {
	gin.SetMode(gin.ReleaseMode)
}

func TestGetAuthToken(t *testing.T){
	auth := getAuthToken(
		EVENT_HUB_NAME,
		EVENT_HUB_ROUTE,
		EVENT_HUB_PRIMARY_KEY_NAME,
		EVENT_HUB_PRIMARY_KEY )
	auth2 := "https%3A%2F%2Fiot-button.servicebus.windows.net%2Fcoldchaintrack%2Fmessages&sig=mF1t1jmWb9aJcicIYLzqbzEBKqYqDn/1Mw8Qtv4kHPw%3D&se=1550272832&skn=WifiManageSharedAccessKey"

	assert.Equal(t, auth, auth2)
}

func main(){
	// Current epoch time
	fmt.Printf("Current epoch time is:\t\t\t%d\n\n", currentEpochTime())

	// Convert from human readable date to epoch
	humanReadable := time.Now()
	fmt.Printf("Human readable time is:\t\t\t%s\n", humanReadable)
	fmt.Printf("Human readable to epoch time is:\t%d\n\n", humanReadableToEpoch(humanReadable))


	// Convert from epoch to human readable date
	epoch := currentEpochTime()
	fmt.Printf("Epoch to human readable time is:\t%s\n", epochToHumanReadable(epoch))

}

func currentEpochTime() int64 {
	return time.Now().Unix()
}

func humanReadableToEpoch(date time.Time) int64 {
	return date.Unix()
}

func epochToHumanReadable(epoch int64) time.Time {
	return time.Unix(epoch, 0)
}