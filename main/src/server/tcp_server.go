package server

import (
	"bufio"
	"bytes"
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/robfig/go-cache"
	"github.com/scando1993/ColdChainTrackerServer/main/src/database"
	"github.com/scando1993/ColdChainTrackerServer/main/src/models"
	"math/rand"
	"net"
	"net/http"

	"strconv"
	"strings"
	"time"
)

var TCP_Port = "65432"
var EVENT_HUB_HOST = "iot-button.servicebus.windows.net"
var EVENT_HUB_PORT = ""
var EVENT_HUB_ROUTE = "coldchaintrack/messages"
var EVENT_HUB_AUTHORIZATION string = "SharedAccessSignature sr=https%3a%2f%2fiot-button.servicebus.windows.net%2fcoldchaintrack%2fmessages&sig=T21GiLWj5psgBT0sURW%2b15c3%2baVLm%2bsuEZoErhBGTbM%3d&se=1549555760&skn=RootManageSharedAccessKey"
var CONFIG_DATA bool = true

const MIN = 1
const MAX = 100

func random() int {
	return rand.Intn(MAX-MIN) + MIN
}

var (
	httpClient *http.Client
	routeCache *cache.Cache
)

const (
	MaxIdleConnections int = 20
	RequestTimeout     int = 300
)

// init HTTPClient
func init() {
	httpClient = createHTTPClient()
	routeCache = cache.New(5*time.Minute, 10*time.Minute)
}

// createHTTPClient for connection re-use
func createHTTPClient() *http.Client {
	client := &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: MaxIdleConnections,
		},
		Timeout: time.Duration(RequestTimeout) * time.Second,
	}

	return client
}

func sendToEventHub(data models.RawSensorData){
	url := "https://" + EVENT_HUB_HOST +  "/" + EVENT_HUB_ROUTE
	logger.Log.Debug(url)
	bPayload, err := json.Marshal(data)
	if err != nil {
		err = errors.Wrap(err, "problem marshaling data")
		//aChan <- a{err: err}
		return
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bPayload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", EVENT_HUB_AUTHORIZATION)
	req.Header.Set("Host", EVENT_HUB_HOST)
	logger.Log.Debug(req)
	resp, err := httpClient.Do(req)
	if err != nil {
		err = errors.Wrap(err, "problem posting payload")
		//aChan <- a{err: err}
		return
	}
	logger.Log.Debugf("Request with data sent to event hub with response: %", resp.StatusCode)
	logger.Log.Debug(resp)
	defer resp.Body.Close()
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func getDeviceFromFamily(device string) (family string){
	families,_ := database.GetFamilies()
	for _, _family := range families{
		d, _err := database.Open(_family, true)
		if _err != nil {
			logger.Log.Warn(_err)
			return
		}
		defer d.Close()
		s, _err := d.GetDevices()
		if _err != nil {
			logger.Log.Warn(_err)
			return
		}
		if stringInSlice(device, s) {
			family = _family
			return
		}
	}

	return
}

func bindStringData(rawData string, data *models.RawSensorData) (err error){
	s := strings.Split(rawData, ",")
	if len(s) != 5 {
		err = errors.Wrap(err,"Too few or too many data")
		return
	}

	data.SensorId = s[0]
	time, _err := strconv.ParseUint(s[1],10, 32)
	if _err != nil{
		err = errors.Wrap(_err,"Error in time parsing")
		return
	}
	data.Timestamp = time

	t, _err := strconv.ParseUint(s[2], 10, 16)
	if _err != nil{
		err = errors.Wrap(_err,"Error in temperature parsing")
		return
	}
	data.Temperature = int(t)

	mac1, _err := getMACAddress(s[3][:12])
	if _err != nil{
		err = errors.Wrap(_err,"Error in mac1 parsing")
		return
	}

	mac2, _err := getMACAddress(s[3][12:24])
	if _err != nil{
		err = errors.Wrap(_err,"Error in mac2 parsing")
		return
	}

	mac3, _err := getMACAddress(s[3][24:36])
	if _err != nil{
		err = errors.Wrap(_err,"Error in mac3 parsing")
		return
	}

	_rss1, _err := strconv.ParseUint(s[3][36:38], 16 ,8)
	if _err != nil{
		err = errors.Wrap(_err,"Error in rss1 parsing")
		return
	}
	rss1 := int(_rss1) - 256

	_rss2, _err := strconv.ParseUint(s[3][38:40], 16 ,8)
	if _err != nil{
		err = errors.Wrap(_err,"Error in rss2 parsing")
		return
	}
	rss2 := int(_rss2) - 256

	_rss3, _err := strconv.ParseUint(s[3][40:42], 16 ,8)
	if _err != nil{
		err = errors.Wrap(_err,"Error in rss3 parsing")
		return
	}
	rss3 := int(_rss3) - 256

	r1 := models.Router{
		Mac:mac1,
		Rssi:rss1,
	}

	r2 := models.Router{
		Mac:mac2,
		Rssi:rss2,
	}

	r3 := models.Router{
		Mac:mac3,
		Rssi:rss3,
	}

	wifi := make([]models.Router, 0)
	wifi = append(wifi, r1)
	wifi = append(wifi, r2)
	wifi = append(wifi, r3)

	data.Wifi = wifi

	b, _err := strconv.ParseUint(s[4], 16, 16)
	if _err != nil{
		err = errors.Wrap(_err,"Error in batery parsing")
		return
	}

	data.Battery = int(b)

	data.RawSensorData = true

	data.Family = getDeviceFromFamily(data.SensorId)

	return
}

func getMACAddress (s string) (mac string, err error){
	var _mac string = s[0:2]

	for i := 2; i < len(s); i = i + 2 {
		_mac += ":"
		_mac += s[i:i+2]
	}
	mac = strings.ToLower(_mac)
	logger.Log.Debug(mac)

	//hw,_err := net.ParseMAC(mac)
	//
	//if _err != nil{
	//	logger.Log.Warn(_err)
	//	err = _err
	//	return
	//}

	//mac = hw.String()
	//err = nil
	return
}

func handleConnection(c net.Conn) {
	logger.Log.Infof("Serving %s", c.RemoteAddr().String())
	for {
		buff := bufio.NewScanner(c)
		err := buff.Scan()
		netData := buff.Text()
		//netData, err := bufio.NewReader(c).ReadString('\n')
		if !err  {
			logger.Log.Warn("No more scanned data")
			break
		}

		//temp := strings.TrimSpace(string(netData))
		temp := string(netData)
		logger.Log.Debug(temp)

		switch temp {
			case "end":{
				logger.Log.Debug("Enter close")
				goto close
			}
			case "doneAll":{
				if CONFIG_DATA {
					result := "wait"
					_, _err := c.Write([]byte(string(result)))
					if _err != nil {
						logger.Log.Warn(_err)
					}
					CONFIG_DATA = false
				}else{
					result := "ok"
					_, _err := c.Write([]byte(string(result)))
					if _err != nil {
						logger.Log.Warn(_err)
					}
					goto close
				}
			}
			case "dataTime":{
				logger.Log.Debug("Send dataTime")
				result := "=2"
				_, _err := c.Write([]byte(string(result)))
				if _err != nil {
					logger.Log.Warn(_err)
				}
			}
			case "sendTime":{
				logger.Log.Debug("Send sendTime")
				result := "=2"
				_, _err := c.Write([]byte(string(result)))
				if _err != nil {
					logger.Log.Warn(_err)
				}
			}
			case "ssid":{
				logger.Log.Debug("Send ssid")
				result := `=IoTPacificsoft`
				_, _err := c.Write([]byte(string(result)))
				if _err != nil {
					logger.Log.Warn(_err)
				}
			}
			case "pass":{
				logger.Log.Debug("Send pass")
				result := `=IoTPacificsoft`
				_, _err := c.Write([]byte(string(result)))
				if _err != nil {
					logger.Log.Warn(_err)
				}
			}
			default:{
				var d models.RawSensorData
				__err := bindStringData(temp, &d)
				if __err != nil{
					logger.Log.Warn(__err)
					result := "fail"
					_, _err := c.Write([]byte(string(result)))
					if _err != nil {
						logger.Log.Warn(_err)
					}
				}else{
					result := "ok"
					_, _err := c.Write([]byte(string(result)))
					go sendToEventHub(d)
					if _err != nil {
						logger.Log.Warn(_err)
					}
				}
			}
		}
	}

	close:
		logger.Log.Info("Connection Closed")

		err := c.Close()
		if err != nil {
			logger.Log.Warn(err)
		}
}

func Tcp_Run() (err error) {

	logger.Log.Debug("TCP Server initiated")

	PORT := ":" + TCP_Port
	l, _err := net.Listen("tcp", PORT)

	if _err != nil {
		logger.Log.Warn(_err)
		err = _err
		return
	}

	defer l.Close()
	rand.Seed(time.Now().Unix())

	for {
		c, _err2 := l.Accept()
		if _err2 != nil {
			logger.Log.Warn(_err2)
			err = _err2
			return
		}
		go handleConnection(c)
	}
}
