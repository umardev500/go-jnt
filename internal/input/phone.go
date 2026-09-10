package input

import (
	"bufio"
	"fmt"
	"regexp"
	"strings"
)

var indonesianPhoneRegex = regexp.MustCompile(`^62[0-9]{9,13}$`)

func ScanPhoneNumber(scanner *bufio.Scanner) string {
	for {
		fmt.Print("Enter your WhatsApp phone number: ")

		if !scanner.Scan() {
			continue
		}

		phone := strings.TrimSpace(scanner.Text())

		// Remove common formatting characters.
		phone = strings.ReplaceAll(phone, " ", "")
		phone = strings.ReplaceAll(phone, "-", "")
		phone = strings.ReplaceAll(phone, "(", "")
		phone = strings.ReplaceAll(phone, ")", "")

		// +628123456789 -> 628123456789
		if strings.HasPrefix(phone, "+62") {
			phone = "62" + phone[3:]
		}

		// 08123456789 -> 628123456789
		if strings.HasPrefix(phone, "08") {
			phone = "62" + phone[1:]
		}

		if !indonesianPhoneRegex.MatchString(phone) {
			fmt.Println("Invalid phone number.")
			fmt.Println("Example: 6281234567890")
			fmt.Println()
			continue
		}

		return phone
	}
}
