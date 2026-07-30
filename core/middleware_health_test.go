package core

import (
	"net/http/httptest"
	"sync"
	"testing"
)

func TestSetReadyNoRace(t *testing.T) {
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(2)
		go func() {
			defer wg.Done()
			SetReady(true)
		}()
		go func() {
			defer wg.Done()
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/ready", nil)
			readyHandler().ServeHTTP(w, r)
		}()
	}
	wg.Wait()
}
