package misc

import "math"

type AudioFrames [][]int

// ToMonoFrames converts stereo into mono frames by averaging each samples.
// Note that a stereo frame could have 2 samples in phase opposition which would lead
// to a zero value. This edge case isn't taken in consideration.
func ToMonoFrames(fs AudioFrames) AudioFrames {
	if fs == nil {
		return nil
	}

	mono := make(AudioFrames, len(fs))
	for i, f := range fs {
		mono[i] = []int{AvgInt(f...)}
	}
	return mono
}

// AvgInt averages the int values passed
func AvgInt(xs ...int) int {
	var output int
	for i := 0; i < len(xs); i++ {
		output += xs[i]
	}
	return output / len(xs)
}

// IeeeFloatToInt converts a 10 byte IEEE float into an int.
func IeeeFloatToInt(b [10]byte) int {
	var i uint32
	// Negative number
	if (b[0] & 0x80) == 1 {
		return 0
	}

	// Less than 1
	if b[0] <= 0x3F {
		return 1
	}

	// Too big
	if b[0] > 0x40 {
		return 67108864
	}

	// Still too big
	if b[0] == 0x40 && b[1] > 0x1C {
		return 800000000
	}

	i = (uint32(b[2]) << 23) | (uint32(b[3]) << 15) | (uint32(b[4]) << 7) | (uint32(b[5]) >> 1)
	i >>= (29 - uint32(b[1]))

	return int(i)
}

// IntToIeeeFloat converts an int into a 10 byte IEEE float.
func IntToIeeeFloat(i int) [10]byte {
	b := [10]byte{}
	num := float64(i)

	var sign int
	var expon int
	var fMant, fsMant float64
	var hiMant, loMant uint

	if num < 0 {
		sign = 0x8000
	} else {
		sign = 0
	}

	if num == 0 {
		expon = 0
		hiMant = 0
		loMant = 0
	} else {
		fMant, expon = math.Frexp(num)
		if (expon > 16384) || !(fMant < 1) { /* Infinity or NaN */
			expon = sign | 0x7FFF
			hiMant = 0
			loMant = 0 /* infinity */
		} else { /* Finite */
			expon += 16382
			if expon < 0 { /* denormalized */
				fMant = math.Ldexp(fMant, expon)
				expon = 0
			}
			expon |= sign
			fMant = math.Ldexp(fMant, 32)
			fsMant = math.Floor(fMant)
			hiMant = uint(fsMant)
			fMant = math.Ldexp(fMant-fsMant, 32)
			fsMant = math.Floor(fMant)
			loMant = uint(fsMant)
		}
	}

	b[0] = byte(expon >> 8)
	b[1] = byte(expon)
	b[2] = byte(hiMant >> 24)
	b[3] = byte(hiMant >> 16)
	b[4] = byte(hiMant >> 8)
	b[5] = byte(hiMant)
	b[6] = byte(loMant >> 24)
	b[7] = byte(loMant >> 16)
	b[8] = byte(loMant >> 8)
	b[9] = byte(loMant)

	return b
}

// Uint24to32 converts a 3 byte uint23 into a uint32
func Uint24to32(bytes []byte) uint32 {
	var output uint32
	output |= uint32(bytes[2]) << 0
	output |= uint32(bytes[1]) << 8
	output |= uint32(bytes[0]) << 16

	return output
}

// Uint32toUint24Bytes converts a uint32 into a 3 byte uint24 representation
func Uint32toUint24Bytes(n uint32) []byte {
	bytes := make([]byte, 3)
	bytes[0] = byte(n >> 16)
	bytes[1] = byte(n >> 8)
	bytes[2] = byte(n >> 0)

	return bytes
}
