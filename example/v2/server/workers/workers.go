package workers

import (
	"fmt"
)

func Add(a int, b int) int {
	fmt.Printf("%d+%d=%d\n", a, b, a+b)
	//panic("pppp")
	return a + b
}
