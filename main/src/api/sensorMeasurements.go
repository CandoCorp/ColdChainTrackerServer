package api

import (
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/scando1993/ColdChainTrackerServer/main/src/models"
	"io/ioutil"
	"net/http"
)

var HOST = "104.209.223.100/chaintrack/auth"
var TRACKING_SEARCH_ALL = "api/tracking/trackAll"
var TRACKING_ALL_FAMILIES = "api/tracking/getAllFamilies"
var TRACKING_ALL_DEVICE = "api/tracking/getAllDevice"
var TRACKING_ALL = "api/tracking/all"
var TRACKING_LAST_TEMPERATURE = "api/tracking/lastTemperature"
var TRACKING_GROUP_BY = "api/tracking/groupby"
var CONFIG_DATA bool = true

func GetTrackAllFromWebApp(device string, family string) (data string, err error){
	url := "http://" + HOST +  "/" + TRACKING_SEARCH_ALL

	logger.Log.Debug(url)
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Content-Type", "application/json")
	q := req.URL.Query()
	q.Add("name", device)
	q.Add("family", family)
	req.URL.RawQuery = q.Encode()
	logger.Log.Debug(req)
	resp, _err := httpClient.Do(req)
	if _err != nil {
		_err := errors.Wrap(err, "problem posting requeriment")
		err = _err
		return
	}
	contents, _err := ioutil.ReadAll(resp.Body)
	if _err != nil {
		err = errors.Wrap(err, "problem posting payload")
		return
	}
	data = string(contents)
	logger.Log.Debugf("Request with data sent to event hub with response: %", resp.StatusCode)
	logger.Log.Debug(resp)
	defer resp.Body.Close()

	return data, err
}

func GetAllFamiliesFromWebApp(device string, family string)(data string, err error){
	url := "http://" + HOST +  "/" + TRACKING_ALL_FAMILIES
	logger.Log.Debug(url)
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Content-Type", "application/json")
	logger.Log.Debug(req)
	resp, _err := httpClient.Do(req)
	if _err != nil {
		_err := errors.Wrap(err, "problem posting requeriment")
		err = _err
		return
	}
	contents, _err := ioutil.ReadAll(resp.Body)
	if _err != nil {
		err = errors.Wrap(err, "problem posting payload")
		return
	}
	data = string(contents)
	logger.Log.Debugf("Request with data sent to event hub with response: %", resp.StatusCode)
	logger.Log.Debug(resp)
	defer resp.Body.Close()

	return data, err
}

func GetAllDevicesFromWebApp(family string) (data string, err error){
	url := "http://" + HOST +  "/" + TRACKING_ALL_DEVICE
	logger.Log.Debug(url)
	if err != nil {
		err = errors.Wrap(err, "problem marshaling data")
		//aChan <- a{err: err}
		return
	}
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Content-Type", "application/json")
	q := req.URL.Query()
	q.Add("family", family)
	req.URL.RawQuery = q.Encode()
	logger.Log.Debug(req)
	resp, err := httpClient.Do(req)
	if err != nil {
		err = errors.Wrap(err, "problem posting payload")
		//aChan <- a{err: err}
		return
	}
	contents, _err := ioutil.ReadAll(resp.Body)
	if _err != nil {
		err = errors.Wrap(err, "problem posting payload")
		return
	}
	data = string(contents)
	logger.Log.Debugf("Request with data sent to event hub with response: %", resp.StatusCode)
	logger.Log.Debug(resp)
	defer resp.Body.Close()
	return
}

func GetDataGroupByFamilyFromWebApp(family string) (data []models.WebApiData, err error){
	url := "http://" + HOST +  "/" + TRACKING_GROUP_BY
	logger.Log.Debug(url)
	if err != nil {
		err = errors.Wrap(err, "problem marshaling data")
		//aChan <- a{err: err}
		return
	}
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Content-Type", "application/json")
	q := req.URL.Query()
	q.Add("family", family)
	req.URL.RawQuery = q.Encode()
	logger.Log.Debug(req)
	resp, _err := httpClient.Do(req)
	if _err != nil {
		err = errors.Wrap(_err, "problem posting payload")
		//aChan <- a{err: err}
		return
	}

	var _data []models.WebApiData
	_data = make([]models.WebApiData,0)
	contents, _err := ioutil.ReadAll(resp.Body)
	if _err != nil {
		err = errors.Wrap(_err, "problem posting payload")
		return
	}
	//logger.Log.Debug(string(contents))
	_err = json.Unmarshal(contents, &_data)
	if _err != nil{
		err = errors.Wrap(_err, "problem posting payload")
		return
	}
	data = _data
	logger.Log.Debugf("Request with data sent to web api with response: %", resp.StatusCode)
	logger.Log.Debug(resp)
	defer resp.Body.Close()

	return data, err
}

func GetLastTemperatureFromWebApp(device string, family string) (data string, err error){
	url := "http://" + HOST +  "/" + TRACKING_LAST_TEMPERATURE
	logger.Log.Debug(url)
	if err != nil {
		err = errors.Wrap(err, "problem marshaling data")
		//aChan <- a{err: err}
		return
	}
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Content-Type", "application/json")
	q := req.URL.Query()
	q.Add("name", device)
	q.Add("family", family)
	req.URL.RawQuery = q.Encode()

	logger.Log.Debug(req)
	resp, err := httpClient.Do(req)
	if err != nil {
		err = errors.Wrap(err, "problem posting payload")
		//aChan <- a{err: err}
		return
	}
	contents, _err := ioutil.ReadAll(resp.Body)
	if _err != nil {
		err = errors.Wrap(err, "problem posting payload")
		return
	}
	data = string(contents)
	logger.Log.Debugf("Request with data sent to web api with response: %", resp.StatusCode)
	logger.Log.Debug(resp)
	defer resp.Body.Close()
	return
}

func GetAllFromWebApp() (data string, err error){
	url := "http://" + HOST +  "/" + TRACKING_ALL
	logger.Log.Debug(url)
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Content-Type", "application/json")
	logger.Log.Debug(req)
	resp, _err := httpClient.Do(req)
	if _err != nil {
		_err = errors.Wrap(err, "problem posting payload")
		//aChan <- a{err: err}
		return
	}
	contents, _err := ioutil.ReadAll(resp.Body)
	if _err != nil {
		err = errors.Wrap(err, "problem posting payload")
		return
	}
	data = string(contents)
	logger.Log.Debugf("Request with data sent to web api with response: %", resp.StatusCode)
	logger.Log.Debug(resp)
	defer resp.Body.Close()
	return
}