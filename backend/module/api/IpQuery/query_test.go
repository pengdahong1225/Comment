package IpQuery

import (
	"encoding/json"
	"testing"
)

func TestQueryIpGeolocation(t *testing.T) {
	ip := "110.242.68.66"
	info, err := QueryIpGeolocation(ip)
	if err != nil {
		t.Error(err)
	} else {
		data, _ := json.Marshal(info)
		t.Log(string(data))
	}
}
