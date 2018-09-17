package main

import (
	"fmt"
	"net/http"
)

var rateCounter = make(map[int]int)
var rateBuckets = make([]int, 0)

func incRateCounter() {
	//rateCounter[int(time.Now().Unix()/60)]++
}

func rateCounterCleanup() {

	//keys := make([]int, 0)

	qwe := make(map[int]int)
	qwe[10] = 100
	qwe[20] = 100
	qwe[30] = 100
	qwe[40] = 100
	qwe[50] = 100
	qwe[60] = 100
	qwe[70] = 100
	qwe[80] = 100
	qwe[90] = 100
	qwe[100] = 100

	itemToDelete := 1
	itemItter := 0
	for key, value := range qwe {
		if itemItter == itemToDelete {
			break
		}

		fmt.Println("Key:", key, "Value:", value)

		delete(qwe, key)
		itemItter++
	}

	//fmt.Println(len(qwe), qwe[10])
	fmt.Println(qwe)
}

// RateCouterMiddleware ...
func RateCouterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
		// potential panic error
		go incRateCounter()
	})
}
