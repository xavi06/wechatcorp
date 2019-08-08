package wechatcorp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	baseURL = "https://qyapi.weixin.qq.com"
)

// API struct
type API struct {
	CorpID     string
	CorpSecret string
	AgentID    int
	Client     *http.Client
}

// NewAPI func
func NewAPI(corpID string, corpSecret string, agentID int) *API {
	return &API{corpID, corpSecret, agentID, &http.Client{}}
}

// HTTPDo func
func (api *API) HTTPDo(url string, method string, data []byte) (body []byte, err error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(data))
	if err != nil {
		return
	}
	//client := &http.Client{}
	response, err := api.Client.Do(req)
	if err != nil {
		return
	}
	defer response.Body.Close()
	body, err = ioutil.ReadAll(response.Body)
	return
}

// GetToken func
func (api *API) GetToken() (string, error) {
	url := fmt.Sprintf("%s/cgi-bin/gettoken?corpid=%s&corpsecret=%s", baseURL, api.CorpID, api.CorpSecret)
	body, err := api.HTTPDo(url, "GET", nil)
	if err != nil {
		return "", err
	}
	var data TokenInfo
	err = json.Unmarshal(body, &data)
	if err != nil {
		return "", err
	}
	return data.AccessToken, nil

}

// SendText func
func (api *API) SendText(token, content, user string) error {
	/*
		var postText *PostText
		postText = new(PostText)
		postText.AgentID = api.AgentID
		postText.ToUser = user
		postText.MsgType = "text"
		postText.Safe = 0
	*/
	/*
		var text *Text
		text = new(Text)
	*/
	var postText map[string]interface{}
	postText = make(map[string]interface{})
	postText["agentid"] = api.AgentID
	postText["touser"] = user
	postText["msgtype"] = "text"
	postText["safe"] = 0
	var text map[string]string
	text = make(map[string]string)
	text["content"] = content
	//text.Content = content
	postText["text"] = text
	//fmt.Println(postText)
	postData, err := json.Marshal(postText)
	if err != nil {
		return err
	}
	//
	url := fmt.Sprintf("%s/cgi-bin/message/send?access_token=%s", baseURL, token)
	body, err := api.HTTPDo(url, "POST", postData)
	if err != nil {
		return err
	}
	var respData ResponseSend
	err = json.Unmarshal(body, &respData)
	if err != nil {
		return err
	}
	fmt.Println(respData)
	return nil
}
