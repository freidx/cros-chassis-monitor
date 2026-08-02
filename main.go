package main

import (
	"fmt"
	"os"
	"time"
	"unsafe"

	"golang.org/x/sys/unix"
)

type ec_payload struct {
	Version, Command, Outsize, Insize, Result uint32
	Data                                      [4]byte
}

func main() {
	// If not running as root, exit
	if os.Geteuid() != 0 {
		fmt.Println("Error: Run as root")
		os.Exit(1)
	}

	// Open file for sending requests
	f, err := os.OpenFile("/dev/cros_ec", os.O_RDWR, 0)
	if err != nil {
		fmt.Printf("Error opening /dev/cros_ec: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	const cmd_mac = (3 << 30) | (20 << 16) | (0xEC << 8) | 0

	for {
		cmd := ec_payload{Command: 0x3E0F, Insize: 1}
		_, _, errno := unix.Syscall(
			unix.SYS_IOCTL,
			f.Fd(),
			// RW     Size        EC
			cmd_mac,
			uintptr(unsafe.Pointer(&cmd)),
		)

		if errno != 0 {
			// note: consider notifying externally
			fmt.Printf("Request failed: %v\n", errno)
			continue
		}

		if cmd.Data[0] == 1 {
			unix.Sync()
			// Reboot if chassis is opened.
			if err := unix.Reboot(unix.LINUX_REBOOT_CMD_POWER_OFF); err != nil {
				fmt.Printf("Response failed: %v\n", err)
				// note: consider notifying externally / trying other shutoff attempts
				continue
			}
		}

		time.Sleep(2 * time.Second)
	}
}
