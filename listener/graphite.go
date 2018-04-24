package listener

import (
	"bufio"
	"errors"
	"io"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/ovh/fossil/core"
	log "github.com/sirupsen/logrus"
)

const nanosPerSec = 1000000000
const nanosPerMilli = 1000000

// Writer interface which is used to save graphite datapoints
type Writer interface {
	Write(*core.GTS)
}

// Graphite is a Graphite socket who parse to sensision format
type Graphite struct {
	Writer Writer
	Listen string
	Parse  bool
}

// NewGraphite return a new Graphite initialized with his output chan
func NewGraphite(listen string, writer Writer, p bool) *Graphite {
	return &Graphite{
		Listen: listen,
		Writer: writer,
		Parse:  p,
	}
}

// OpenTCPServer opens the Graphite TCP input format and starts processing data.
func (g *Graphite) OpenTCPServer() error {
	ln, err := net.Listen("tcp", g.Listen)
	if err != nil {
		return err
	}

	log.Infof("Listen on %s", g.Listen)

	for {
		conn, err := ln.Accept()

		if opErr, ok := err.(*net.OpError); ok && !opErr.Temporary() {
			log.WithFields(log.Fields{
				"error": opErr,
			}).Debug("Graphite TCP listener closed")
			continue
		}

		if err != nil {
			log.WithFields(log.Fields{
				"error": err.Error(),
			}).Warn("Error has occurred while accepting the TCP connection")
			continue
		}

		go g.handleTCPConnection(conn)
	}
}

// handleTCPConnection services an individual TCP connection for the Graphite input
func (g *Graphite) handleTCPConnection(conn net.Conn) {
	defer conn.Close()

	var metric string
	reader := bufio.NewReader(conn)
	for {
		buf, _, err := reader.ReadLine()
		if err == io.EOF {
			return
		}
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
				"buf":   buf,
			}).Warn("unable to read TCP payload")
			return
		}

		metric = strings.TrimSpace(string(buf))
		datapoint, err := g.parseLine(metric)
		if err != nil {
			log.WithFields(log.Fields{
				"error":  err,
				"metric": metric,
			}).Info("unable to parse line")
			continue
		}

		log.Debug(datapoint)
		g.Writer.Write(datapoint)
	}
}

func (g *Graphite) parseLine(metric string) (*core.GTS, error) {
	split := strings.Split(metric, " ")

	// Minimum expected format is : 'metrics value timestamp'
	if len(split) < 3 {
		return nil, errors.New("Bad metric format")
	}

	ts, err := strconv.ParseInt(split[2], 10, 64)
	if err != nil {
		return nil, errors.New("Bad metric part: timestamp")
	}

	var value interface{}
	skip := false

	// try to convert the string into a float64
	if strings.Contains(split[1], ".") {
		number, err := strconv.ParseFloat(split[1], 64)
		if err == nil {
			skip = true
			value = number
		}
	}

	// try to convert the string into an integer
	if !skip {
		number, err := strconv.ParseInt(split[1], 10, 64)
		if err == nil {
			skip = true
			value = number
		}
	}

	// try to convert the string into a boolean
	if !skip {
		if strings.ToLower(split[1]) == "true" {
			skip = true
			value = true
		} else if strings.ToLower(split[1]) == "false" {
			skip = true
			value = false
		}
	}

	// assume that the value is a string
	if !skip {
		value = split[1]
	}

	dp := &core.GTS{
		Ts:     int64toTime(ts).UnixNano() / 1000,
		Value:  value,
		Labels: make(map[string]string),
	}

	// Check if there are tags
	if strings.Contains(split[0], ";") {
		subSplit := strings.Split(split[0], ";")
		dp.Name = subSplit[0]

		// If no tags, but auto fill enabled, we map the hierarchy for later by label processing purpose
		if g.Parse {
			classPart := strings.Split(subSplit[0], ".")
			for idx, part := range classPart {
				dp.Labels[strconv.Itoa(idx)] = part
			}
		}

		// Parse tags
		for _, v := range subSplit[1:] {
			tagSplit := strings.Split(v, "=")
			dp.Labels[tagSplit[0]] = tagSplit[1]
		}

	} else {
		dp.Name = split[0]

		// If no tags, but auto fill enabled, we map the hierarchy for later by label processing purpose
		if g.Parse {
			classPart := strings.Split(split[0], ".")
			for idx, part := range classPart {
				dp.Labels[strconv.Itoa(idx)] = part
			}
		}
	}

	return dp, nil
}

// int64toTime Convert an int expressed either in seconds or milliseconds into a Time object
func int64toTime(timestamp int64) time.Time {
	if timestamp == 0 {
		return time.Now()
	}

	timeNanos := timestamp
	// If less than 2^32, assume it's in seconds
	// (in millis that would be Thu Feb 19 18:02:47 CET 1970)
	if timeNanos < 0xFFFFFFFF {
		timeNanos *= nanosPerSec
	} else {
		timeNanos *= nanosPerMilli
	}

	return time.Unix(timeNanos/nanosPerSec, timeNanos%nanosPerSec)
}
