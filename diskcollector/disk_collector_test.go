package diskcollector

import (
	"testing"
)

func TestParseDiskUsage(t *testing.T) {
	passingTests := []struct {
		info     string
		input    string
		expected DiskUsage
	}{
		{
			info:  "valid input",
			input: "Filesystem     1024-blocks      Used Available Capacity iused      ifree %iused  Mounted on\n/dev/disk1s5s1   488245288  15057096 231339568     7%  502388 2313395680    0%   /",
			expected: DiskUsage{
				Filesystem: "/dev/disk1s5s1",
				Size:       488245288,
				Used:       15057096,
				Available:  231301444,
				Capacity:   7,
				IUsed:      502388,
				IFree:      2313014440,
				PIUsed:     0,
				MountedOn:  "/",
			},
		},
		{
			info:  "valid input with space in filesystem",
			input: "Filesystem     1024-blocks      Used Available Capacity iused      ifree %iused  Mounted on\nmap auto_home   488245288  15057096 231339568     7%  502388 2313395680    0%   /",
			expected: DiskUsage{
				Filesystem: "map auto_home",
				Size:       488245288,
				Used:       15057096,
				Available:  231301444,
				Capacity:   7,
				IUsed:      502388,
				IFree:      2313014440,
				PIUsed:     0,
				MountedOn:  "/",
			},
		},
	}

	failingTests := []struct {
		info     string
		input    string
		expected DiskUsage
	}{
		{
			info:     "too few fields",
			input:    "Filesystem     1024-blocks      Used Available Capacity iused      ifree %iused\n/dev/disk1s5s1   488245288  15057096 231339568     7%  502388 2313395680    0%",
			expected: DiskUsage{},
		},
		{
			info:     "no disk info found",
			input:    "Filesystem     1024-blocks      Used Available Capacity iused      ifree %iused  Mounted on",
			expected: DiskUsage{},
		},
		{
			info:     "empty input",
			input:    "",
			expected: DiskUsage{},
		},
	}

	for _, test := range passingTests {
		t.Logf("Test case: %s", test.info)
		result, err := parseDiskUsage(test.input)
		if err != nil {
			t.Errorf("parseDiskUsage(%q) returned error: %v", test.input, err)
		}
		if len(result) == 0 {
			t.Fatalf("parseDiskUsage(%q) returned empty", test.input)
		}
		if result[0] == test.expected {
			t.Errorf("parseDiskUsage(%q) = %v; want %v", test.input, result, test.expected)
		}
	}

	for _, test := range failingTests {
		t.Logf("Test case: %s", test.info)
		_, err := parseDiskUsage(test.input)
		if err == nil {
			t.Errorf("parseDiskUsage(%q) should have returned an error: %v", test.input, err)
		}
		// if len(result) != 0 {
		// 	t.Errorf("parseDiskUsage(%q) = should be empty %v", result)
		// }
	}
}

func TestBytesToHumanReadable(t *testing.T) {
	tests := []struct {
		info     string
		input    uint64
		expected string
	}{
		{
			info:     "bytes",
			input:    1024,
			expected: "1.0 KiB",
		},
		{
			info:     "kilobytes",
			input:    1048576,
			expected: "1.0 MiB",
		},
		{
			info:     "megabytes",
			input:    1073741824,
			expected: "1.0 GiB",
		},
		{
			info:     "gigabytes",
			input:    1099511627776,
			expected: "1.0 TiB",
		},
		{
			info:     "terabytes",
			input:    1125899906842624,
			expected: "1.0 PiB",
		},
	}

	for _, test := range tests {
		t.Logf("Test case: %s", test.info)
		result := bytesToHumanReadable(test.input)
		if result != test.expected {
			t.Errorf("bytesToHumanReadable(%d) = %s; want %s", test.input, result, test.expected)
		}
	}
}
