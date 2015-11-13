package linux

import "testing"

func TestSNMP(t *testing.T) {
	stat, err := ReadSNMP("proc/3323/net/snmp")
	if err != nil {
		t.Fatal("stat read fail")
	}

	if stat.TCP.RetransSegs != 31369 {
		t.Fatal(stat.TCP.RetransSegs)
	}
}
