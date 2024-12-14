package collector

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

type DiskUsage struct {
	Filesystem string
	Size       int64
	Used       int64
	Available  int64
	Capacity   int
	IUsed      int64
	IFree      int64
	PIUsed     int
	MountedOn  string
}

func GetDiskUsage() error {
	cmd := exec.Command("df", "-k")

	out, err := cmd.Output()
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}
	// fmt.Println(string(output)) // Get disk usage
	parseDiskUsage(string(out))
	return nil
}

func parseDiskUsage(data string) ([]DiskUsage, error) {
	var diskUsages []DiskUsage

	fmt.Println("data: ", data)
	lines := strings.Split(data, "\n")
	fmt.Println("# lines: ", lines)

	if len(lines) < 2 {
		return nil, fmt.Errorf("error: no disk info found")
	}

	// TODO: check that the fields are as expected
	// Filesystem     1024-blocks      Used Available Capacity iused      ifree %iused  Mounted on

	for _, line := range lines[1:] {
		fields := strings.Fields(line)

		// If there aren't 8 columns, error
		if len(fields) < 9 {
			return nil, fmt.Errorf("error: not finding 9 fields")
		}

		// Handle the case where the filesystem has a space in it
		// TODO: this is a hack, need to find a better way to handle this
		for len(fields) > 9 {
			fields[0] = fields[0] + " " + fields[1]
			fields = append(fields[:1], fields[2:]...)
		}

		disk := DiskUsage{
			Filesystem: fields[0],
			Size:       blocksToBytes(fields[1]),
			Used:       blocksToBytes(fields[2]),
			Available:  blocksToBytes(fields[3]),
			Capacity:   toPercent(fields[4]),
			IUsed:      blocksToBytes(fields[5]),
			IFree:      blocksToBytes(fields[6]),
			PIUsed:     toPercent(fields[7]),
			MountedOn:  fields[8],
		}
		diskUsages = append(diskUsages, disk)
	}
	return diskUsages, nil
}

// TODO: handle error
func blocksToBytes(blockStr string) int64 {
	block, err := strconv.ParseInt(blockStr, 10, 64)
	if err != nil {
		fmt.Println("Error:", err)
	}
	return 1024 * block
}

// TODO: handle error
func toPercent(str string) int {
	percent, err := strconv.Atoi(strings.TrimSuffix(str, "%"))
	if err != nil {
		fmt.Println("Error:", err)
	}
	return percent
}

func bytesToHumanReadable(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
