package locker_client

import (
	"net/http"
	"sync"
	"testing"
)

var testUrl = "http://127.0.0.1:33000"

func TestMultiGoroutineWrite(_ *testing.T) {
	wg := sync.WaitGroup{}
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(index int) {
			db := CloudLockerClient{
				Client: http.DefaultClient,
				Url:    testUrl,
			}
			for j := 0; j < 1000; j++ {
				key := []byte{byte(index)}
				value := []byte{byte(j % 100)}
				_ = db.Set(key, value)
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
	db := CloudLockerClient{
		Client: http.DefaultClient,
		Url:    testUrl,
	}
	for i := 0; i < 100; i++ {
		v, _ := db.Get([]byte{byte(i)})
		if len(v) == 0 || v[0] != 99 {
			panic("v should be 99")
		}
	}
}

func TestBasics(t *testing.T) {
	db := CloudLockerClient{
		Client: http.DefaultClient,
		Url:    testUrl,
	}
	for i := 0; i < 100; i++ {
		_ = db.Set([]byte{byte(i)}, []byte{byte(i + 1)})
		v, _ := db.Get([]byte{byte(i)})
		if len(v) == 0 || v[0] != byte(i+1) {
			panic("value not set exactly")
		}
	}
}
