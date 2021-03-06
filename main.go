package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

// The Gelf data structure
type Gelf struct {
	// mandatory fields
	Version      string  `json:"version"`
	Host         string  `json:"host"`
	ShortMessage string  `json:"short_message"`
	FullMessage  string  `json:"full_message"`
	Timestamp    float64 `json:"timestamp"`
	Level        int     `json:"level"`
	// additional fields
	LogType   string `json:"_logType"` // _logType will show as logType in Graylog
	SourceEnv string `json:"_source_env"`
	Type      string `json:"_type"`
	MessageId int    `json:"_messageId"`
	DateTime  string `json:"_dateTime"`
}

// Strng implements the Stringer interface
func (g Gelf) String() string {
	buf, err := json.Marshal(g)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%s", string(buf))
}

// GrayLogConn stores the connectionto a graylog server
type GrayLogConn struct {
	GrayLogServer string
	GrayLogPort   int
	conn          net.Conn
}

// connect to a graylog server
func (g *GrayLogConn) connect(graylogServer string, graylogPort int) error {
	serverAddr := graylogServer + ":" + strconv.Itoa(graylogPort)
	var conn net.Conn
	conn, err := net.Dial(protocol, serverAddr)
	if err != nil {
		return err
	}
	g.conn = conn
	return nil
}

// close graylog connection
func (g *GrayLogConn) close() {
	g.conn.Close()
}

// sendToGrayLog sends the gelf structure to the graylog server
func (g *Gelf) sendToGrayLog(conn net.Conn) error {

	buf, err := json.Marshal(g)
	if err != nil {
		return err
	}
	if verbose {
		fmt.Printf("%s\n", string(buf))
	}
	buf = append(buf, byte(0))                              // GELF over TCP needs a nullbyte at the end
	err = conn.SetDeadline(time.Now().Add(5 * time.Second)) // must be set before every call to resert the timeout
	if err != nil {
		return err
	}

	_, err = conn.Write(buf)
	if err != nil {
		return err
	}
	return nil
}

// commandline options
var (
	graylogServer     string
	graylogPort       int
	verbose           bool
	count             int
	sleepMilliseconds int
	logType           string
	sourceEnv         string
	protocol          string
)

const RFC3339Milli = "2006-01-02T15:04:05.999Z07:00"

func main() {
	var gelf Gelf
	var gconn GrayLogConn

	flag.StringVar(&graylogServer, "graylog", "localhost", "The Graylog server")
	flag.StringVar(&graylogServer, "g", "localhost", "The Graylog server (shorthand)")
	flag.IntVar(&graylogPort, "port", 12201, "The port of the Graylog server")
	flag.IntVar(&graylogPort, "p", 12201, "The port of the Graylog server (shorthand)")
	flag.BoolVar(&verbose, "verbose", false, "Be verbose")
	flag.BoolVar(&verbose, "v", false, "Be verbose")
	flag.IntVar(&count, "count", 1, "Number of messages to send")
	flag.IntVar(&count, "c", 1, "Number of messages to send (shorthand)")
	flag.IntVar(&sleepMilliseconds, "s", 0, "Sleeptime in milliseconds between sends")
	flag.StringVar(&logType, "logtype", "APP", "The logtype (APP or EVENT)")
	flag.StringVar(&logType, "t", "APP", "The logtype (APP or EVENT) (shorthand)")
	flag.StringVar(&sourceEnv, "sourceenv", "dev", "the source_env field in the message")
	flag.StringVar(&sourceEnv, "e", "dev", "the source_env field in the message (shorthand)")
	flag.StringVar(&protocol, "protocol", "tcp", "the protocol used")
	flag.StringVar(&protocol, "P", "tcp", "the protocol used (shorthand)")

	flag.Parse()
	//fmt.Printf("len: %d, value: %v\n", len(flag.Args()), flag.Args())

	if len(flag.Args()) > 0 {
		message := strings.Join(flag.Args(), " ")
		gelf.ShortMessage = message
		gelf.FullMessage = message
	} else {
		// read from stdin
		reader := bufio.NewReader(os.Stdin)
		message, _ := reader.ReadString('\n')
		message = strings.TrimRight(message, "\r\n")
		gelf.ShortMessage = message
		gelf.FullMessage = message
	}
	gelf.Version = "1.1"
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}
	gelf.Host = hostname
	gelf.Level = 6
	gelf.LogType = logType
	gelf.SourceEnv = sourceEnv
	if logType == "APP" {
		gelf.Type = "applog-gelftest"

	} else {
		gelf.Type = "eventlog-gelftest"

	}
	err = gconn.connect(graylogServer, graylogPort)
	if err != nil {
		log.Fatal(err)
	}
	defer gconn.close()

	for i := 0; i < count; i++ {
		if i > 0 {
			time.Sleep(time.Duration(sleepMilliseconds) * time.Millisecond)
		}
		gelf.MessageId = i + 1
		t := time.Now()
		gelf.Timestamp = math.Round(float64(t.UnixNano())/1e9*1000.0) / 1000.0 // round milliseconds
		gelf.DateTime = t.Format(RFC3339Milli)
		err := gelf.sendToGrayLog(gconn.conn)
		if err != nil {
			log.Fatal(err)
		}
	}
}
