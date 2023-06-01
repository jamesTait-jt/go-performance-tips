package preallocate

func growSlice() {
	toAdd := []byte("123456")
	bs := []byte{}

	bs = append(bs, toAdd...)
	bs = append(bs, toAdd...)
	bs = append(bs, toAdd...)
}

func growSlicePreAllocate() {
	toAdd := []byte("123456")
	bs := make([]byte, 0, 18)

	bs = append(bs, toAdd...)
	bs = append(bs, toAdd...)
	bs = append(bs, toAdd...)
}