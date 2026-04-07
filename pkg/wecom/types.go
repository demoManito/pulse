package wecom

import "encoding/xml"

// Config 企业微信应用配置。
type Config struct {
	CorpID              string `yaml:"corp_id"`
	CorpSecret          string `yaml:"corp_secret"`
	CallbackToken       string `yaml:"callback_token"`
	CallbackEncodingAES string `yaml:"callback_encoding_aes_key"`
	AgentID             int    `yaml:"agent_id"`
}

// TokenResponse /cgi-bin/gettoken 接口响应。
type TokenResponse struct {
	ErrCode     int    `json:"errcode"`
	ErrMsg      string `json:"errmsg"`
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

// CallbackEvent 解密后的企业微信回调事件消息。
type CallbackEvent struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   string   `xml:"ToUserName"`
	FromUserName string   `xml:"FromUserName"`
	CreateTime   int64    `xml:"CreateTime"`
	MsgType      string   `xml:"MsgType"`
	Event        string   `xml:"Event"`
	ChangeType   string   `xml:"ChangeType"`
}

// DocEvent 文档相关的回调事件。
type DocEvent struct {
	CallbackEvent
	DocID string `xml:"DocId"`
}

// DocListResponse /cgi-bin/wedoc/doc_list 接口响应。
type DocListResponse struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
	Data    struct {
		DocList []DocInfo `json:"doc_list"`
	} `json:"data"`
}

// DocInfo 文档列表中的文档信息。
type DocInfo struct {
	DocID      string `json:"docid"`
	DocName    string `json:"doc_name"`
	DocType    int    `json:"doc_type"`
	CreateTime int64  `json:"create_time"`
	ModifyTime int64  `json:"modify_time"`
}

// DocContentResponse /cgi-bin/wedoc/document/get 接口响应。
type DocContentResponse struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
	Data    struct {
		Doc struct {
			DocContent string `json:"doc_content"`
		} `json:"doc"`
	} `json:"data"`
}
