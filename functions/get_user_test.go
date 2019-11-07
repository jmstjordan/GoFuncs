package gecko

import (
	"encoding/json"
	"io"
	"log"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetUser(t *testing.T) {
	expectedUser := User{
		FirstName:    "Jan-Michael",
		LastName:     "Jordan",
		Email:        "jmstjordan@gmail.com",
		EmployeeType: "Adminstrator",
		UserId:       "ZkIJ6FtcBGRBeppt8N5Hdl7f9x33",
	}
	tests := []struct {
		uid  string
		want User
	}{
		{uid: "ZkIJ6FtcBGRBeppt8N5Hdl7f9x33", want: expectedUser},
	}

	for _, test := range tests {
		req := httptest.NewRequest("GET", "/?user="+test.uid, nil)
		req.Header.Add("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		GetUser(rr, req)

		var usr User
		dec := json.NewDecoder(strings.NewReader(rr.Body.String()))
		if err := dec.Decode(&usr); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
		if got := usr; got != test.want {
			t.Errorf("GetUserHTTP(%q) = %q, want %q", test.uid, got, test.want)
		}
	}
}
