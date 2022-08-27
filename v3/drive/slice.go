package drive

func rPush(l [][]byte, v []byte) [][]byte {
	return append(l, v)
}
func lPop(l [][]byte) ([]byte, [][]byte) {
	if len(l) == 0 {
		return nil, l
	}
	return l[0], l[1:]
}

func lPush(l [][]byte, v []byte) [][]byte {
	var n = make([][]byte, len(l)+1)
	copy(n[1:], l)
	n[0] = v
	return n
}
