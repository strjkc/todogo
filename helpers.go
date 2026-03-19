package main

import "errors"

func search(num int, lst []*Task) (int, error) {
	i := 0
	j := len(lst) - 1
	for i <= j {
		mid := (i + j) / 2
		if lst[mid].ID > num {
			j = mid - 1
		} else if lst[mid].ID < num {
			i = mid + 1
		} else {
			return mid, nil
		}
	}
	return -1, errors.New("index not found")
}
