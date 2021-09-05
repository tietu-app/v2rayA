package configure

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/tidwall/gjson"
	"github.com/v2rayA/v2rayA/core/serverObj"
	"github.com/v2rayA/v2rayA/core/vmessInfo"
)

type ServerRawV2 struct {
	ServerObj serverObj.ServerObj `json:"serverObj"`
	Latency   string              `json:"latency"`
}

type SubscriptionRawV2 struct {
	Remarks string        `json:"remarks,omitempty"`
	Address string        `json:"address"`
	Status  string        `json:"status"` //update time, error info, etc.
	Servers []ServerRawV2 `json:"servers"`
	Info    string        `json:"info"` // maybe include some info from provider
}

func Bytes2SubscriptionRaw2(b []byte) (*SubscriptionRawV2, error) {
	var s SubscriptionRawV2
	rawList := gjson.GetBytes(b, "servers").Array()
	for _, raw := range rawList {
		var obj serverObj.ServerObj
		obj, err := serverObj.New(raw.Get("serverObj.protocol").String())
		if err != nil {
			return nil, err
		}
		s.Servers = append(s.Servers, ServerRawV2{ServerObj: obj})
	}
	if err := jsoniter.Unmarshal(b, &s); err != nil {
		return nil, err
	}
	return &s, nil
}

func Bytes2ServerRaw2(b []byte) (*ServerRawV2, error) {
	var s ServerRawV2
	var obj serverObj.ServerObj
	obj, err := serverObj.New(gjson.GetBytes(b, "serverObj.protocol").String())
	if err != nil {
		return nil, err
	}
	s.ServerObj = obj
	if err := jsoniter.Unmarshal(b, &s); err != nil {
		return nil, err
	}
	return &s, nil
}

type ServerRaw struct {
	VmessInfo vmessInfo.VmessInfo `json:"vmessInfo"`
	Latency   string              `json:"latency"`
}

type SubscriptionRaw struct {
	Remarks string      `json:"remarks,omitempty"`
	Address string      `json:"address"`
	Status  string      `json:"status"` //update time, error info, etc.
	Servers []ServerRaw `json:"servers"`
	Info    string      `json:"info"` // maybe include some info from provider
}
