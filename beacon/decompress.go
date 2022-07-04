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
			i = i + endChar + 1
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

// func decompressBucketLogNumber(input string) uint64 {
// 	return 999
// }

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

// Case 1

// Compressed
// "00.3bl._.3e..2v..1t.k000D.3x..5s..8n.P.2c..6u..2h.*7*0H.4c..2m.2"

// Decompressed
// 0: 0
// 1: 4305
// 2: 62
// 3: 122
// 4: 103
// 5: 65
// 6: 20
// 7: 0
// 8: 0
// 9: 0
// 10: 39
// 11: 141
// 12: 208
// 13: 311
// 14: 51
// 15: 84
// 16: 246
// 17: 89
// 18: 0
// 19: 0
// 20: 0
// 21: 0
// 22: 0
// 23: 0
// 24: 0
// 25: 43
// 26: 156
// 27: 94
// 28: 2

// Case 2

// Compressed
// "00.43..54..28..3z..29."

// Decompressed
// 0: 0
// 1: 147
// 2: 184
// 3: 80
// 4: 143
// 5: 81

// Case 3

// Compressed
// 0*j*0td

// Decompressed
// 0: 0
// 1: 0
// 2: 0
// 3: 0
// 4: 0
// 5: 0
// 6: 0
// 7: 0
// 8: 0
// 9: 0
// 10: 0
// 11: 0
// 12: 0
// 13: 0
// 14: 0
// 15: 0
// 16: 0
// 17: 0
// 18: 0
// 19: 29
// 20: 13

// Case 4

// Compressed
// 000.40..28..4x..54..5b..45.m*8*0.6c..27..3h..b7..e8.k*9*0O.5f..3u..3s..2l.lS.5v..3a..38..3q..3a..24.

// Decompressed
// 0: 0
// 1: 0
// 2: 144
// 3: 80
// 4: 177
// 5: 184
// 6: 191
// 7: 149
// 8: 22
// 9: 0
// 10: 0
// 11: 0
// 12: 0
// 13: 0
// 14: 0
// 15: 0
// 16: 0
// 17: 228
// 18: 79
// 19: 125
// 20: 403
// 21: 512
// 22: 20
// 23: 0
// 24: 0
// 25: 0
// 26: 0
// 27: 0
// 28: 0
// 29: 0
// 30: 0
// 31: 0
// 32: 50
// 33: 195
// 34: 138
// 35: 136
// 36: 93
// 37: 21
// 38: 54
// 39: 211
// 40: 118
// 41: 116
// 42: 134
// 43: 118
// 44: 76
