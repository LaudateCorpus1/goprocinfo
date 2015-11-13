package linux

import (
	"io/ioutil"
	"reflect"
	"strconv"
	"strings"
)

//Tcp: RtoAlgorithm RtoMin RtoMax MaxConn ActiveOpens PassiveOpens AttemptFails EstabResets CurrEstab InSegs OutSegs RetransSegs InErrs OutRsts InCsumErrors
type Snmp struct {
	TCP SnmpTCP
}

type SnmpTCP struct {
	// Tcp
	RtoAlgorithm uint64
	RtoMin       uint64
	RtoMax       uint64
	MaxConn      uint64
	ActiveOpens  uint64
	PassiveOpens uint64
	AttemptFails uint64
	EstabResets  uint64
	CurrEstab    uint64
	InSegs       uint64
	OutSegs      uint64
	RetransSegs  uint64
	InErrs       uint64
	OutRsts      uint64
	InCsumErrors uint64
}

func ReadSNMP(path string) (*Snmp, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(data), "\n")

	// Maps a netstat metric to its value (i.e. SyncookiesSent --> 0)
	statMap := make(map[string]string)

	// patterns
	// TcpExt: SyncookiesSent SyncookiesRecv SyncookiesFailed... <-- header
	// TcpExt: 0 0 1764... <-- values
	for i := 1; i < len(lines); i = i + 2 {
		headers := strings.Fields(lines[i-1][strings.Index(lines[i-1], ":")+1:])
		values := strings.Fields(lines[i][strings.Index(lines[i], ":")+1:])
		for j, header := range headers {
			statMap[header] = values[j]
		}
	}

	out := Snmp{}
	elem := reflect.ValueOf(&out.TCP).Elem()
	typeOfElem := elem.Type()

	for i := 0; i < elem.NumField(); i++ {
		if val, ok := statMap[typeOfElem.Field(i).Name]; ok {
			parsedVal, _ := strconv.ParseUint(val, 10, 64)
			elem.Field(i).SetUint(parsedVal)
		}
	}

	return &out, nil
}
