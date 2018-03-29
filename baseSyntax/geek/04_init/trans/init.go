package trans //一个目录中只能包含一个package且目录名称要和包名称保持一致

import "math"

var Pi float64

func init() {
	Pi = 4 * math.Atan(1) // init() function computes Pi
}