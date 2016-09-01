package main

import "fmt"

func main() {
	nums := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	input := make([]*job, 0)
	for _, n := range nums {
		input = append(input, &job{num: n})
	}

	out, fail := NewManager(2).SetRetry(4).Manage(input).Report()

	fmt.Printf("Output is a %d-element slice.\n", len(out))
	for _, r := range out {
		fmt.Println(r.num)
	}

	fmt.Printf("Fail is a %d-element slice.\n", len(fail))
	for _, f := range fail {
		fmt.Println(f.num)
	}
}
