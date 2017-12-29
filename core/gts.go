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

	// value
	switch gts.Value.(type) {
	case bool:
		if gts.Value.(bool) {
			sensision += "T"
		} else {
			sensision += "F"
		}

	case float64:
		sensision += fmt.Sprintf("%f", gts.Value.(float64))

	case int64:
		sensision += fmt.Sprintf("%d", gts.Value.(int64))

	case float32:
		sensision += fmt.Sprintf("%f", gts.Value.(float32))

	case int:
		sensision += fmt.Sprintf("%d", gts.Value.(int))

	case string:
		sensision += fmt.Sprintf("'%s'", url.QueryEscape(gts.Value.(string)))

	default:
		// Other types: just output their default format
		strVal := fmt.Sprintf("%v", gts.Value)
		sensision += url.QueryEscape(strVal)
	}

	// According to https://github.com/golang/go/blob/master/src/fmt/print.go#L1148
	sensision += "\n"
	return []byte(sensision)
}
