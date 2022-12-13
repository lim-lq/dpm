package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/lim-lq/dpm/core"
	"github.com/lim-lq/dpm/metadata"
	"github.com/lim-lq/dpm/models"
)

var Api = "http://localhost:8890/api"

func Pretty(data interface{}) string {
	jsonBytes, _ := json.MarshalIndent(data, "", "    ")
	return string(jsonBytes)
}

type HttpClient struct {
	cli    *http.Client
	url    string
	method string
	data   []byte
	req    *http.Request
}

func GetHttpClient() *HttpClient {
	return &HttpClient{cli: &http.Client{}}
}

func (h *HttpClient) Request(url string) *HttpClient {
	h.url = url
	return h
}

func (h *HttpClient) Get(url string) *HttpClient {
	h.method = "GET"
	h.url = url
	return h
}

func (h *HttpClient) Post(url string, data []byte) *HttpClient {
	h.method = "POST"
	h.data = data
	h.url = url
	return h
}

func (h *HttpClient) Put(url string, data []byte) *HttpClient {
	h.method = "PUT"
	h.data = data
	h.url = url
	return h
}

func (h *HttpClient) Do() ([]byte, error) {
	var req *http.Request
	var err error
	if len(h.data) > 0 {
		req, err = http.NewRequest(h.method, h.url, bytes.NewReader(h.data))
	} else {
		req, err = http.NewRequest(h.method, h.url, nil)
	}
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJEUE0iLCJleHAiOjE2NzM1MTEzOTV9.j8fcphDd3bPQK_7onJiOUMIgPLV1Vuz2m9cZobf9GDE")
	h.req = req
	resp, err := h.cli.Do(h.req)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(resp.Body)
}

func encryptPass(pass string) (string, error) {
	// 获取pubkey
	client := GetHttpClient()
	resp, err := client.Get(fmt.Sprintf("%s/common/publickey", Api)).Do()
	if err != nil {
		return "", err
	}
	data := metadata.Response{}
	err = json.Unmarshal(resp, &data)
	if err != nil {
		return "", err
	}
	pubStr := data.Info.(string)
	return core.GetRsaClient().EncryptByPubkey(pubStr, pass)
}

func TestCreateAccount(t *testing.T) {
	cipher, err := encryptPass("67890")
	if err != nil {
		t.Fatalf("encrypt password error: %v", err)
	}
	t.Log(cipher)
	account := models.Account{
		Username: "lcq",
		Password: cipher,
	}
	client := GetHttpClient()
	bodyBytes, err := json.Marshal(account)
	if err != nil {
		t.Fatalf("json marshal error: %v", err)
	}
	respBytes, err := client.Post(fmt.Sprintf("%s/accounts", Api), bodyBytes).Do()
	if err != nil {
		t.Fatalf("Create account error: %v", err)
	}
	resp := metadata.Response{}
	err = json.Unmarshal(respBytes, &resp)
	if err != nil {
		t.Fatalf("Umarshal response error: %v", err)
	}
	t.Logf("Response info: %s", resp.Info)
}

func TestUpdateAccount(t *testing.T) {
	account := metadata.MapStr{"id": "1", "is_admin": true}
	httpCli := GetHttpClient()
	bodyBytes, err := json.Marshal(account)
	if err != nil {
		t.Fatalf("json marshal error: %v", err)
	}
	respBytes, err := httpCli.Put(fmt.Sprintf("%s/accounts/0", Api), bodyBytes).Do()
	if err != nil {
		t.Fatalf("Update account error: %v", err)
	}
	resp := metadata.Response{}
	err = json.Unmarshal(respBytes, &resp)
	if err != nil {
		t.Fatalf("json unmarshal error: %v", err)
	}
	t.Logf("Response info: %+v", resp)
}

func TestQueryAccount(t *testing.T) {
	cond := metadata.Condition{
		Page:    metadata.Page{Limit: 10, Offset: 0},
		Filters: metadata.Filters{"email": map[string]string{"$regex": "licongqing"}},
	}
	bodyBytes, err := json.Marshal(&cond)
	if err != nil {
		t.Fatalf("json marshal error: %v", err)
	}
	respBytes, err := GetHttpClient().Post(fmt.Sprintf("%s/accounts/search", Api), bodyBytes).Do()
	if err != nil {
		t.Fatalf("Get account error: %v", err)
	}
	resp := metadata.Response{}
	err = json.Unmarshal(respBytes, &resp)
	if err != nil {
		t.Fatalf("Json unmarshal error: %v", err)
	}
	t.Logf("Account List: %s", Pretty(resp))
}

func TestChangePassword(t *testing.T) {
	httpCli := GetHttpClient()
	cipher, err := encryptPass("67890")
	if err != nil {
		t.Fatalf("Encrypt cipher error: %v", err)
	}
	data := metadata.MapStr{"cipher": cipher}
	bodyBytes, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("Json marshal error: %v", err)
	}
	respBytes, err := httpCli.Put(fmt.Sprintf("%s/accounts/1/changepass", Api), bodyBytes).Do()
	if err != nil {
		t.Fatalf("Change Password error: %v", err)
	}
	resp := metadata.Response{}
	err = json.Unmarshal(respBytes, &resp)
	if err != nil {
		t.Fatalf("Json unmarshal error: %v", err)
	}
	t.Log(Pretty(resp))
}
