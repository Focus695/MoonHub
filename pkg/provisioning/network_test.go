package provisioning

import (
	"reflect"
	"testing"
)

func TestParseWifiListSSIDWithColons(t *testing.T) {
	stdout := "foo:bar:baz:75:WPA2:*\n"
	got := parseWifiList(stdout)
	want := []WifiNetwork{
		{SSID: "foo:bar:baz", Signal: 75, Security: "WPA2", InUse: true},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("parseWifiList() = %#v, want %#v", got, want)
	}
}

func TestParseWifiListSimple(t *testing.T) {
	stdout := "MyWifi:80:WPA2:\n"
	got := parseWifiList(stdout)
	if len(got) != 1 || got[0].SSID != "MyWifi" || got[0].Signal != 80 || got[0].Security != "WPA2" || got[0].InUse {
		t.Fatalf("unexpected: %#v", got)
	}
}
