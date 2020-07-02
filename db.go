package locker_client

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type DB interface {
	Get(key []byte) ([]byte, error)
	Set(key, value []byte) error
	Delete(key []byte) error
}

var _ DB = CloudLockerClient{}

type CloudLockerClient struct {
	*http.Client
	Url string
}

//if key is not exist, return nil slice and error is nil
func (c CloudLockerClient) Get(key []byte) ([]byte, error) {
	resp, _ := http.Post(c.Url+"/get", "application/json", strings.NewReader(string(key)))
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return body, err
}

func (c CloudLockerClient) Set(key, value []byte) error {
	e := Entry{
		K: key,
		V: value,
	}
	b, _ := json.Marshal(e)
	_, err := http.Post(c.Url+"/set", "application/json", strings.NewReader(string(b)))
	return err
}

func (c CloudLockerClient) Delete(key []byte) error {
	_, err := http.Post(c.Url+"/delete", "application/json", strings.NewReader(string(key)))
	return err
}

type Entry struct {
	K HexBytes `json:"k"`
	V HexBytes `json:"v"`
}

type HexBytes []byte

func (bz HexBytes) MarshalJSON() ([]byte, error) {
	s := strings.ToUpper(hex.EncodeToString(bz))
	jbz := make([]byte, len(s)+2)
	jbz[0] = '"'
	copy(jbz[1:], []byte(s))
	jbz[len(jbz)-1] = '"'
	return jbz, nil
}

// This is the point of Bytes.
func (bz *HexBytes) UnmarshalJSON(data []byte) error {
	if len(data) < 2 || data[0] != '"' || data[len(data)-1] != '"' {
		return fmt.Errorf("invalid hex string: %s", data)
	}
	bz2, err := hex.DecodeString(string(data[1 : len(data)-1]))
	if err != nil {
		return err
	}
	*bz = bz2
	return nil
}
