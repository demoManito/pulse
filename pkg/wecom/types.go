package wecom

import "encoding/xml"

// Config 企业微信应用配置。
type Config struct {
	CorpID              string     `yaml:"corp_id"`
	CorpSecret          string     `yaml:"corp_secret"`
	CallbackToken       string     `yaml:"callback_token"`
	CallbackEncodingAES string     `yaml:"callback_encoding_aes_key"`
	AgentID             int        `yaml:"agent_id"`
	Sync                SyncConfig `yaml:"sync"`
}

// SyncConfig 文档同步配置。
type SyncConfig struct {
	// WikiID 要同步的知识库 ID，为空则同步所有文档。
	WikiID string `yaml:"wiki_id"`
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

// DocInfo 文档信息。
type DocInfo struct {
	DocID      string `json:"docid"`
	DocName    string `json:"doc_name"`
	DocType    int    `json:"doc_type"`
	CreateTime int64  `json:"create_time"`
	ModifyTime int64  `json:"modify_time"`
	// WikiFileID 知识库内文件标识，格式：<wiki_id>_<file_id>，非知识库文档为空。
	WikiFileID string `json:"wiki_fileid,omitempty"`
	// SpaceID 所属空间 ID。
	SpaceID string `json:"spaceid,omitempty"`
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
