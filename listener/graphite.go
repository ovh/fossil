package listener

import (
	"bufio"
	"errors"
	"io"
	"net"
	"strconv"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/ovh/fossil/core"
)

// Graphite is a Graphite socket who parse to sensision format
type Graphite struct {
	Output chan *core.GTS
	Listen string
}

// NewGraphite return a new Graphite initialized with his output chan
func NewGraphite(listen string) *Graphite {
	return &Graphite{
		Output: make(chan *core.GTS),
		Listen: listen,
	}
}

// OpenTCPServer opens the Graphite TCP input format and starts processing data.
func (g *Graphite) OpenTCPServer() error {
	ln, err := net.Listen("tcp", g.Listen)
	if err != nil {
		return err
	}
	log.Info("Listen on", g.Listen)

	go func() {
		for {
			conn, err := ln.Accept()

			if opErr, ok := err.(*net.OpError); ok && !opErr.Temporary() {
				log.Debug("graphite TCP listener closed")
				continue
			}
			if err != nil {
				log.WithFields(log.Fields{
					"error": err.Error(),
				}).Warn("error accepting TCP connection")
				continue
			}
			go g.handleTCPConnection(conn)
		}
	}()
	return nil
}

// handleTCPConnection services an individual TCP connection for the Graphite input
func (g *Graphite) handleTCPConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	var metric string

	for {
		buf, _, err := reader.ReadLine()
		if err == io.EOF {
			log.Info("EOF")
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

		log.Info(datapoint)
		// if nothing eat the chan, prevent to handle next line
		g.Output <- datapoint
	}
}

func (g *Graphite) parseLine(metric string) (*core.GTS, error) {

	split := strings.Split(metric, " ")
	// From metrics with love
	if len(split) < 3 {
		return nil, errors.New("Bad metric format")
	}

	ts, err := strconv.ParseInt(split[2], 10, 64)
	if err != nil {
		return nil, errors.New("Bad metric part: timestamp")
	}

	dp := &core.GTS{
		Ts:     int64toTime(ts).UnixNano() / 1000,
		Name:   split[0],
		Value:  split[1],
		Labels: make(map[string]string),
	}

	classPart := strings.Split(split[0], ".")

	i := 0
	for _, part := range classPart {
		dp.Labels[strconv.Itoa(i)] = part
		i++
	}

	return dp, nil
}

// int64toTime Convert an int expressed either in seconds or milliseconds into a Time object
func int64toTime(timestamp int64) time.Time {
	if timestamp == 0 {
		return time.Now()
	}
	const nanosPerSec = 1000000000
	const nanosPerMilli = 1000000

	timeNanos := timestamp
	if timeNanos < 0xFFFFFFFF {
		// If less than 2^32, assume it's in seconds
		// (in millis that would be Thu Feb 19 18:02:47 CET 1970)
		timeNanos *= nanosPerSec
	} else {
		timeNanos *= nanosPerMilli
	}

	return time.Unix(timeNanos/nanosPerSec, timeNanos%nanosPerSec)
}
