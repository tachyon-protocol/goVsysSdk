package goVsysSdk

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"
	"strconv"
	"sync"
)

const defaultTimeOut = 15 * time.Second
var gClient *http.Client
var gClientOnce sync.Once
func getHttpClient() *http.Client{
	gClientOnce.Do(func(){
		gClient = &http.Client{
			Timeout: defaultTimeOut,
		}
	})
	return gClient
}

func (a *VsysApi) httpPost(path string, data interface{}) (body []byte, err error) {
	url:=a.nodeAddress +path
	client := getHttpClient()
	d, err := json.Marshal(data)
	if err != nil {
		return []byte{}, err
	}
	resp, err := client.Post(url, "application/json", bytes.NewBuffer(d))
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}
	if err := getErrResp(resp, body); err != nil {
		return []byte{}, err
	}
	return body, nil
}

func (a *VsysApi) httpGet(path string) (body []byte, err error) {
	url:=a.nodeAddress +path
	httpReq,err:=http.NewRequest("GET",url,nil)
	if err!=nil{
		return nil,err
	}
	if a.req.ApiKey!=""{
		httpReq.Header.Set("api_key",a.req.ApiKey)
	}
	client := getHttpClient()
	resp, err := client.Do(httpReq)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}
	err = getErrResp(resp, body)
	if err != nil {
		return []byte{}, err
	}
	return body, nil
}

func getErrResp(resp *http.Response, body []byte) (err error) {
	if resp.StatusCode != 200 {
		errResp := CommonResp{}
		err := json.Unmarshal(body, &errResp)
		if err != nil {
			return errors.New("StatusCodeError ["+strconv.Itoa(int(resp.StatusCode))+"] ["+string(body)+"] "+err.Error())
		} else {
			return errors.New("hrw6sqkdv6 ["+strconv.Itoa(int(resp.StatusCode))+"] "+errResp.Message+" ["+string(body)+"]")
		}
	}
	return nil
}
