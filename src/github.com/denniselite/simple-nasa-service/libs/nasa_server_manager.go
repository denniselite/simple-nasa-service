package libs

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"github.com/denniselite/simple-nasa-service/structs"
	"strconv"
)

const (
	APIKeyQueryParam = "api_key"
	pathNeo          = "/neo/browse"
)

type NasaServerManager struct {
	Config structs.NasaServerManager
}

// Получение информации о сервере
func (nsm *NasaServerManager) GetNEOInfo(page int) (neos structs.NasaResponse, err error) {

	params := make(map[string]string)
	params[APIKeyQueryParam] = nsm.Config.APIKey
	params["page"] = strconv.Itoa(page)
	params["size"] = "100"

	response, err := nsm.doRequest(http.MethodGet, pathNeo, params)
	if err != nil {
		return
	}
	defer (*response).Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &neos)
	return
}

func (nsm *NasaServerManager) doRequest(method string, path string, requestData map[string]string) (res *http.Response, err error) {

	requestUrl := fmt.Sprintf("%s%s%s", nsm.Config.EndPoint, nsm.Config.APIVersion, path)
	request, err := http.NewRequest(method, requestUrl, nil)
	if err != nil {
		return
	}

	queryParams := request.URL.Query()
	for key, val := range requestData {
		queryParams.Add(key, val)
	}
	request.URL.RawQuery = queryParams.Encode()

	var httpClient http.Client
	res, err = httpClient.Do(request)
	return
}
