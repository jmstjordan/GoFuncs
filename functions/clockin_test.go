package gecko

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"
)

func TestClockIn(t *testing.T) {

	clock := Clock{
		UserId:    "ZkIJ6FtcBGRBeppt8N5Hdl7f9x33",
		ClockType: "In",
	}
	tests := []struct {
		body Clock
		want int
	}{
		{body: clock, want: 201},
	}

	for _, test := range tests {
		body, _ := json.Marshal(test.body)
		req := httptest.NewRequest("POST", "/", bytes.NewReader(body))
		req.Header.Add("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		ClockIn(rr, req)

		if got := rr.Code; got != test.want {
			t.Errorf("ClockInHTTP(%q) = %q, want %q", test.body, got, test.want)
		}
	}
}
