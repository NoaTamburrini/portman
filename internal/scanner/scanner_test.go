package scanner

import (
	"sort"
	"strconv"
	"testing"
)

const frenchNetstat = "\r\n" +
	"Connexions actives\r\n" +
	"\r\n" +
	"  Proto  Adresse locale         Adresse distante       État\r\n" +
	"  TCP    0.0.0.0:135            0.0.0.0:0              ÉCOUTE          1124\r\n" +
	"  TCP    0.0.0.0:3000           0.0.0.0:0              ÉCOUTE          8840\r\n" +
	"  TCP    127.0.0.1:3000         127.0.0.1:52011        ÉTABLIE         8840\r\n" +
	"  TCP    127.0.0.1:52011        127.0.0.1:3000         ÉTABLIE         9912\r\n" +
	"  TCP    127.0.0.1:52100        127.0.0.1:3000         TIME_WAIT       0\r\n" +
	"  TCP    [::]:135               [::]:0                 ÉCOUTE          1124\r\n" +
	"  UDP    0.0.0.0:5353           *:*                                    2210\r\n" +
	"  UDP    [::1]:1900             *:*                                    3320\r\n"

func TestParseWindowsOutputLocalized(t *testing.T) {
	ports, err := parseWindowsOutput(frenchNetstat)
	if err != nil {
		t.Fatal(err)
	}

	got := map[string]string{}
	for _, p := range ports {
		got[p.Protocol+":"+itoa(p.Number)+":"+itoa(p.PID)] = p.State
	}

	want := map[string]string{
		"tcp:135:1124":   "LISTEN",
		"tcp:3000:8840":  "LISTEN",
		"tcp:52011:9912": "ÉTABLIE",
		"udp:5353:2210":  "",
		"udp:1900:3320":  "",
	}

	if len(got) != len(want) {
		keys := make([]string, 0, len(got))
		for k := range got {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		t.Fatalf("got %d ports %v, want %d", len(got), keys, len(want))
	}
	for k, state := range want {
		if s, ok := got[k]; !ok || s != state {
			t.Errorf("%s: got state %q (present=%v), want %q", k, s, ok, state)
		}
	}

	if targets := FindKillTargets(ports, 3000); len(targets) != 1 || targets[0].PID != 8840 {
		t.Errorf("kill targets for 3000 = %+v, want only PID 8840", targets)
	}
}

func itoa(n int) string {
	return strconv.Itoa(n)
}
