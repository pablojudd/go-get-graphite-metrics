package GraphiteData

import (
	"crypto/tls"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
)

// GraphiteMetrics struct
type GraphiteMetrics struct {
	Target     string           `json:"target"`
	Datapoints [][2]NullFloat64 `json:"datapoints"`
}

// NullFloat64 struct
type NullFloat64 struct {
	sql.NullFloat64
}

// GetMetrics from Graphite
func GetMetrics(url string, target string, dur int) (data []GraphiteMetrics, err error) {

	// Attempt to connect to Graphite and perform query
	tr := &http.Transport{
		// Accept insecure key
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	query := fmt.Sprintf("%s/render?target=%s&from=-%dmin&format=json", url, target, dur)
	resp, err := client.Get(query)
	if err != nil {
		return data, fmt.Errorf("Graphite Server Error: %s", err)
	}
	defer resp.Body.Close()

	// Decode JSON response from Graphite
	err = json.NewDecoder(resp.Body).Decode(&data)

	// Return decoded data
	return data, err

}

// UnmarshalJSON correctly de-serializes a NullFloat64 from JSON
func (n *NullFloat64) UnmarshalJSON(b []byte) error {
	var i interface{}
	if err := json.Unmarshal(b, &i); err != nil {
		return err
	}
	return n.Scan(i)
}
