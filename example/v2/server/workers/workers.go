package workers

import (
	"fmt"
)

func Add(a int, b int) error {
	fmt.Printf("%d+%d=%d\n", a, b, a+b)
	return nil
}
