package internal

func Lerp(c1, c2 uint32, param float64) (out uint32) {
	out = 0
	var ch1, ch2 float64
	ch1 = float64(c1 & 0xff)
	ch2 = float64(c2 & 0xff)
	out |= uint32((ch2-ch1)*param + (ch1))
	ch1 = float64((c1 >> 8) & 0xff)
	ch2 = float64((c2 >> 8) & 0xff)
	out |= uint32((ch2-ch1)*param+(ch1)) << 8
	ch1 = float64((c1 >> 16) & 0xff)
	ch2 = float64((c2 >> 16) & 0xff)
	out |= uint32((ch2-ch1)*param+(ch1)) << 16

	return
}

func Scale(c1 uint32, param float64) uint32 {
	return Lerp(c1, 0, param)
}
