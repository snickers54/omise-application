package rot128

func Rot128(buf []byte) {
	for idx := range buf {
		buf[idx] += 128
	}
}
