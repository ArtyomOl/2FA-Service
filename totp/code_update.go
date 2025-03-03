package totp

import (
	"fmt"
	"time"
)

func UpdateAllCodes() {
	ticker := time.NewTicker(5 * time.Second)
	for range ticker.C {
		fmt.Println("Timer fired!")
	}
}
