package main

import "log"

func main() {
	res := Cop()
	log.Println(res)
}
func Cop() int {
	input := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 0}
	for _, k := range input {
		if k == 7 {
			return k
		}
		//return 0
	}
	return 0
}
