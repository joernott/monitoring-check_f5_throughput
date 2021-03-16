package checker

// Executes nagios check
import (
	"encoding/json"
	"io/ioutil"
	"math"
	"os"
	"time"

	"github.com/gosnmp/gosnmp"
	"github.com/olorin/nagiosplugin"
)

type jsonHistory struct {
	Timestamp      uint64 `json:"Timestamp"`
	ClientBytesIn  uint64 `json:"ClientBytesIn"`
	ClientBytesOut uint64 `json:"ClientBytesOut"`
	ServerBytesIn  uint64 `json:"ServerBytesIn"`
	ServerBytesOut uint64 `json:"ServerBytesOut"`
}

func NewHistory() jsonHistory {
	var h jsonHistory
	h.Timestamp = uint64(time.Now().Unix())
	h.ClientBytesIn = 0
	h.ClientBytesOut = 0
	h.ServerBytesIn = 0
	h.ServerBytesOut = 0
	return h
}

func getHistory(statsFile string) jsonHistory {
	h := NewHistory()
	f, err := os.Open(statsFile)
	if err != nil {
		return h
	}
	defer f.Close()
	raw, _ := ioutil.ReadAll(f)
	json.Unmarshal([]byte(raw), &h)
	return h
}

func putHistory(statsFile string, Timestamp uint64, ClientBytesIn uint64, ClientBytesOut uint64, ServerBytesIn uint64, ServerBytesOut uint64) error {
	var h jsonHistory

	h.Timestamp = Timestamp
	h.ClientBytesIn = ClientBytesIn
	h.ClientBytesOut = ClientBytesOut
	h.ServerBytesIn = ServerBytesIn
	h.ServerBytesOut = ServerBytesOut
	j, err := json.Marshal(h)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(statsFile, j, 0640)
	return err
}

func initData(statsFile string) (map[string]uint64, []string, []string) {
	data := make(map[string]uint64)
	h := getHistory(statsFile)
	oids := []string{
		".1.3.6.1.4.1.3375.2.1.1.2.1.3.0",
		".1.3.6.1.4.1.3375.2.1.1.2.1.5.0",
		".1.3.6.1.4.1.3375.2.1.1.2.1.10.0",
		".1.3.6.1.4.1.3375.2.1.1.2.1.12.0",
	}
	stats := []string{
		"ClientBytesIn",
		"ClientBytesOut",
		"ServerBytesIn",
		"ServerBytesOut",
	}
	data["last_time"] = h.Timestamp
	data["last_ClientBytesIn"] = h.ClientBytesIn
	data["last_ClientBytesOut"] = h.ClientBytesOut
	data["last_ServerBytesIn"] = h.ServerBytesIn
	data["last_ServerBytesOut"] = h.ServerBytesOut
	data["timestamp"] = uint64(time.Now().Unix())
	return data, oids, stats
}

func addResults(result *gosnmp.SnmpPacket, oids []string, stats []string, data map[string]uint64) map[string]uint64 {
	for _, variable := range result.Variables {
		var stat string
		for i, oid := range oids {
			if variable.Name == oid {
				stat = stats[i]
			}
		}
		data[stat] = gosnmp.ToBigInt(variable.Value).Uint64()
	}
	data["duration"] = data["timestamp"] - data["last_time"]
	for _, s := range stats {
		data["delta_"+s] = data[s] - data["last_"+s]
		if data["duration"] != 0 {
			data["throughput_"+s] = data["delta_"+s] * 8 / data["duration"]
		} else {
			data["throughput_"+s] = 0
		}
	}
	if data["duration"] != 0 {
		data["throughput_in"] = (data["delta_ClientBytesIn"] + data["delta_ClientBytesOut"]) * 8 / data["duration"]
		data["throughput_out"] = (data["delta_ServerBytesIn"] + data["delta_ServerBytesOut"]) * 8 / data["duration"]
	} else {
		data["throughput_in"] = 0
		data["throughput_out"] = 0

	}
	return data
}

func addPerfdata(check *nagiosplugin.Check, data map[string]uint64, stats []string, warn float64, crit float64) {
	for _, s := range stats {
		check.AddPerfDatum(s, "c", float64(data[s]), 0.0, math.Inf(1), 0.0, 0.0)
		check.AddPerfDatum("throughput_"+s, "", float64(data["throughput_"+s]), 0.0, math.Inf(1), 0.0, 0.0)
	}
	check.AddPerfDatum("throughput_in", "", float64(data["throughput_in"]), 0.0, math.Inf(1), warn, crit)
	check.AddPerfDatum("throughput_out", "", float64(data["throughput_out"]), 0.0, math.Inf(1), 0.0, 0.0)
}

func Check(
	host string,
	port uint16,
	community string,
	warningThreshold string,
	criticalThreshold string,
	statsFile string,
) {
	var warnRange *nagiosplugin.Range
	var critRange *nagiosplugin.Range
	var err error

	check := nagiosplugin.NewCheck()
	defer check.Finish()
	warn := 0.0
	crit := 0.0
	if warningThreshold != "" {
		warnRange, err = nagiosplugin.ParseRange(warningThreshold)
		if err != nil {
			check.AddResult(nagiosplugin.UNKNOWN, "error parsing warning range")
			return
		}
		warn = warnRange.End
	}
	if criticalThreshold != "" {
		critRange, err = nagiosplugin.ParseRange(criticalThreshold)
		if err != nil {
			check.AddResult(nagiosplugin.UNKNOWN, "error parsing critical range")
			return
		}
		crit = critRange.End
	}
	gosnmp.Default.Target = host
	gosnmp.Default.Port = port
	gosnmp.Default.Community = community
	err = gosnmp.Default.Connect()
	if err != nil {
		check.AddResult(nagiosplugin.UNKNOWN, "SNMP Connect() error: "+err.Error())
		return
	}
	defer gosnmp.Default.Conn.Close()
	data, oids, stats := initData(statsFile)
	result, err := gosnmp.Default.Get(oids)
	if err != nil {
		check.AddResult(nagiosplugin.UNKNOWN, "Get() error: "+err.Error())
	}
	data = addResults(result, oids, stats, data)
	addPerfdata(check, data, stats, warn, crit)

	err = putHistory(statsFile, data["timestamp"], data["ClientBytesIn"], data["ClientBytesOut"], data["ServerBytesIn"], data["ServerBytesOut"])
	if err != nil {
		check.AddResult(nagiosplugin.UNKNOWN, "Write Stats file "+statsFile+" error: %v"+err.Error())
	}
	if warningThreshold+criticalThreshold == "" {
		check.AddResult(nagiosplugin.OK, "Everything is fine")
	} else {
		if criticalThreshold != "" && critRange.Check(float64(data["throughput_in"])) {
			check.AddResult(nagiosplugin.CRITICAL, "Critical: Inbound traffic")
		} else {
			if warningThreshold != "" && warnRange.Check(float64(data["throughput_in"])) {
				check.AddResult(nagiosplugin.WARNING, "Warning: Inbound traffic")
			} else {
				check.AddResult(nagiosplugin.OK, "Everything is fine")
			}
		}

	}
}
