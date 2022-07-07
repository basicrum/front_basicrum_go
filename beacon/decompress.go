package beacon

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/martinlindhe/base36"
)

const COMPRESS_MODE_SMALL_NUMBERS = "0"
const COMPRESS_MODE_LARGE_NUMBERS = "1"
const COMPRESS_MODE_PERCENT = "2"

const LARGE_NUMBER_WRAP = "."

func DecompressBucketLog(data string) []uint64 {
	var out []uint64

	if len(data) == 0 {
		return out
	}

	var endChar int
	var num uint64

	// strip the type out
	logType := string(data[0])
	logData := data[1:]

	// decompress string
	repeat := uint64(1)

	logDataLen := len(logData)

	i := 0

	for i < logDataLen {
		fmt.Println("Scanning: " + strconv.Itoa(i))
		if string(logData[i]) == "*" {
			// this is a repeating number

			// up to the next * is the repeating count (base 36)
			endChar = strings.Index(logData[i+1:], "*")

			if endChar == -1 {
				i++
				continue
			}

			repeat = base36.Decode(logData[i+1 : i+endChar+1])

			// after is the number
			i = i + endChar + 2
			continue
		} else if string(logData[i]) == LARGE_NUMBER_WRAP {
			// this is a number larger than 63

			// up to the next wrap character is the number (base 36)
			endChar = strings.Index(logData[i+1:], LARGE_NUMBER_WRAP)

			if endChar == -1 {
				i++
				continue
			}

			num = base36.Decode(logData[i+1 : i+endChar+1])

			// move to this end char
			i = i + endChar + 1
		} else {
			if logType == COMPRESS_MODE_SMALL_NUMBERS {
				// this digit is a number from 0 to 63
				num = decompressBucketLogNumber(string(logData[i]))
			} else if logType == COMPRESS_MODE_LARGE_NUMBERS {
				// look for this digit to end at a comma

				endChar = strings.Index(logData[i:], ",")

				if endChar != -1 {
					// another index exists later, read up to that
					num = base36.Decode(logData[i : i+endChar])

					// move to this end char
					i = i + endChar
				} else {
					// this is the last number
					num = base36.Decode(string(logData[i]))

					// we're done
					i = logDataLen
				}
			} else if logType == COMPRESS_MODE_PERCENT {
				// check if this is 100
				if logData[i:i+2] == "__" {
					num = 100
				} else {
					convNum, _ := strconv.Atoi(logData[i : i+2])
					num = uint64(convNum)
				}

				// take two characters
				i++
			}
		}

		fmt.Println(num)

		out = append(out, num)

		j := uint64(1)

		for j < repeat {
			out = append(out, num)
			j++
		}

		repeat = 1

		i++
	}

	return out
}

func decompressBucketLogNumber(input string) uint64 {
	// if (!input || !input.charCodeAt) {
	// 	return 0;
	// }

	// convert to ASCII character codeDecompressBucketLog
	chr := uint64([]byte(input)[0])

	fmt.Println(input)

	if chr >= 48 && chr <= 57 {
		// 0 - 9
		return chr - 48
	} else if chr >= 97 && chr <= 122 {
		// a - z
		return chr - 97 + 10
	} else if chr >= 65 && chr <= 90 {
		// A - Z
		return (chr - 65) + 36
	} else if chr == 95 {
		// -
		return 62
	} else if chr == 45 {
		// _
		return 63
	} else {
		// unknown
		return 0
	}
}
