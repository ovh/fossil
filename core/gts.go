package core

import (
	"fmt"
	"net/url"
)

// GTS is a Warp10 representation of GeoTimeSerie
type GTS struct {
	Ts     int64
	Name   string
	Labels map[string]string
	Value  interface{}
}

// Encode a GTS to the Sensision format
// TS/LAT:LON/ELEV NAME{LABELS} VALUE
func (gts *GTS) Encode() []byte {
	sensision := fmt.Sprintf("%d// %s{", gts.Ts, url.QueryEscape(gts.Name))

	sep := ""
	for k, v := range gts.Labels {
		sensision += sep + url.QueryEscape(k) + "=" + url.QueryEscape(v)
		sep = ","
	}
	sensision += "} "
	sensision += fmt.Sprintf("%s", gts.Value)
	sensision += "\r\n"

	return []byte(sensision)
}
