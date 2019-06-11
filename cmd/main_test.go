package main

import (
	"encoding/json"
	"testing"
	"time"
)

func TestJSONTime_MarshalJSON(t *testing.T) {
	tt := JSONTime(time.Now())
	bs, err := json.Marshal([]JSONTime{tt, tt})
	if err != nil {
		t.Error(err)
		return
	}
	s := string(bs)
	if len(s) == 0 {
		t.Error("empty string")
		return
	}
}
