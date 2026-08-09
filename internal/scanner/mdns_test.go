package scanner

import (
	"testing"
)

// dnsName encodes a name as length-prefixed labels + terminator (no compression).
func dnsName(name string) []byte { return encodeDNSName(name) }

func u16(v int) []byte { return []byte{byte(v >> 8), byte(v)} }

// record builds one uncompressed resource record.
func record(name string, rtype int, rdata []byte) []byte {
	var r []byte
	r = append(r, dnsName(name)...)
	r = append(r, u16(rtype)...)  // type
	r = append(r, u16(0x8001)...) // class (cache-flush bit; parser ignores class)
	r = append(r, 0, 0, 0, 60)    // ttl
	r = append(r, u16(len(rdata))...)
	r = append(r, rdata...)
	return r
}

func TestParseAndBuildMDNS_FullChain(t *testing.T) {
	// A response advertising an IPP printer:
	//   PTR  _ipp._tcp.local           -> Office._ipp._tcp.local
	//   SRV  Office._ipp._tcp.local     -> printbox.local
	//   A    printbox.local            -> 192.168.1.50
	var msg []byte
	msg = append(msg, 0, 0, 0x84, 0) // id, flags (response)
	msg = append(msg, u16(0)...)     // questions
	msg = append(msg, u16(3)...)     // answers
	msg = append(msg, u16(0)...)     // authority
	msg = append(msg, u16(0)...)     // additional

	msg = append(msg, record("_ipp._tcp.local", typePTR, dnsName("Office._ipp._tcp.local"))...)
	srvData := append([]byte{0, 0, 0, 0, 0x1F, 0x90}, dnsName("printbox.local")...) // prio,weight,port + target
	msg = append(msg, record("Office._ipp._tcp.local", typeSRV, srvData)...)
	msg = append(msg, record("printbox.local", typeA, []byte{192, 168, 1, 50})...)

	aRecords := map[string]string{}
	srv := map[string]string{}
	var ptrs []ptrEntry
	parseMDNSInto(msg, aRecords, srv, &ptrs)

	if aRecords["printbox.local"] != "192.168.1.50" {
		t.Errorf("A record wrong: %v", aRecords)
	}
	if srv["Office._ipp._tcp.local"] != "printbox.local" {
		t.Errorf("SRV record wrong: %v", srv)
	}
	if len(ptrs) != 1 || ptrs[0].service != "_ipp._tcp.local" {
		t.Errorf("PTR record wrong: %v", ptrs)
	}

	devices := buildDevicesFromMDNS(aRecords, srv, ptrs)
	if len(devices) != 1 {
		t.Fatalf("expected 1 device, got %d", len(devices))
	}
	d := devices[0]
	if d.IP != "192.168.1.50" || d.Hostname != "printbox" || d.Type != TypePrinter {
		t.Errorf("device wrong: %+v", d)
	}
}

func TestDecodeDNSName_Compression(t *testing.T) {
	// "local" at offset 12, then "printbox" + pointer to offset 12.
	msg := make([]byte, 12)
	localOff := len(msg)
	msg = append(msg, dnsName("local")...) // "local" root at localOff
	nameOff := len(msg)
	msg = append(msg, 8)
	msg = append(msg, []byte("printbox")...)
	msg = append(msg, 0xC0, byte(localOff)) // pointer to "local"

	name, next, err := decodeDNSName(msg, nameOff)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if name != "printbox.local" {
		t.Errorf("name = %q, want printbox.local", name)
	}
	if next != len(msg) {
		t.Errorf("next = %d, want %d", next, len(msg))
	}
}

func TestDecodeDNSName_Malformed(t *testing.T) {
	if _, _, err := decodeDNSName([]byte{0x05, 'a'}, 0); err == nil {
		t.Error("expected error on truncated label")
	}
	// A pointer that points to itself must not loop forever.
	if _, _, err := decodeDNSName([]byte{0xC0, 0x00}, 0); err == nil {
		t.Error("expected error on self-referential pointer")
	}
}

func TestMDNSTypeForService(t *testing.T) {
	if mdnsTypeForService("_airplay._tcp.local") != TypeMediaPlayer {
		t.Error("airplay should be media player")
	}
	if mdnsTypeForService("_ssh._tcp.local") != TypeServer {
		t.Error("ssh should be server")
	}
	if mdnsTypeForService("_unknown._tcp.local") != "" {
		t.Error("unknown service should be empty")
	}
}

func TestBuildMDNSQuery(t *testing.T) {
	q := buildMDNSQuery("_ipp._tcp.local", typePTR)
	if q[5] != 1 {
		t.Errorf("questions = %d, want 1", q[5])
	}
	// tail: type PTR, class with QU bit.
	tail := q[len(q)-4:]
	if tail[0] != 0 || tail[1] != typePTR || tail[2] != 0x80 || tail[3] != 0x01 {
		t.Errorf("tail = % x, want 00 0c 80 01", tail)
	}
}
