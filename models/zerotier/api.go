package zerotier

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/astaxie/beego"
)

var (
	zerotierCtlApiRoot   = beego.AppConfig.String("zerotierctlapiroot")
	zerotierCtlAuthToken = beego.AppConfig.String("zerotierctlauthtoken")
	zerotierNetworkId    = beego.AppConfig.String("zerotiernetworkid")
)

func init() {
	zerotierCtlApiRoot = strings.TrimSuffix(zerotierCtlApiRoot, "/")
}

func zerotierCtlApi(subpath string) string {
	return zerotierCtlApiRoot + subpath + "?auth=" + zerotierCtlAuthToken
}

func AddMember(address string) bool {
	url := zerotierCtlApi("/network/" + zerotierNetworkId + "/member/" + address)
	body, _ := json.Marshal(map[string]interface{}{
		"authorized": true,
	})
	res, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err == nil && res.StatusCode == 200 {
		return true
	}
	return false
}

func RemoveMember(address string) bool {
	url := zerotierCtlApi("/network/" + zerotierNetworkId + "/member/" + address)

	// req, _ := http.NewRequest(http.MethodDelete, url, nil)
	// res, err := http.DefaultClient.Do(req)
	// if err == nil && res.StatusCode == 200 {
	// 	succeeded = true
	// }

	// Set field "authorized" to false instead of deleting it,
	// because the ZeroTier's API seems to have problem deleting members
	body, _ := json.Marshal(map[string]interface{}{
		"authorized": false,
	})
	res, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err == nil && res.StatusCode == 200 {
		return true
	}

	return false
}

func GetMember(address string) (member Member, err error) {
	url := zerotierCtlApi("/network/" + zerotierNetworkId + "/member/" + address)
	res, err := http.Get(url)
	if err != nil {
		return
	}
	dec := json.NewDecoder(res.Body)
	err = dec.Decode(&member)
	return
}
