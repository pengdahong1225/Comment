package IpQuery

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// http://ip-api.com/json/110.242.68.66?fields=status,message,continent,country,region,regionName,city,district,zip,lat,lon,timezone,query&lang=zh-CN

const (
	baseUrl = "http://ip-api.com/json"
	fields  = "status,message,continent,country,region,regionName,city,district,zip,lat,lon,timezone,query"
	lang    = "zh-CN"
)

func QueryIpGeolocation(ip string) (*IPInfoResp, error) {
	// 创建一个HTTP客户端
	client := &http.Client{}

	// 构建请求URL
	url := fmt.Sprintf("%s/%s?fileds=%s&lang=%s", baseUrl, ip, fields, lang)
	fmt.Println(url)

	// 发送GET请求
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 解析响应体为JSON
	ipInfo := &IPInfoResp{}
	err = json.NewDecoder(resp.Body).Decode(&ipInfo)
	if err != nil {
		return nil, err
	}

	if ipInfo.Status != "success" {
		return nil, fmt.Errorf("ipInfo.Status != success")
	}

	return ipInfo, nil
}
