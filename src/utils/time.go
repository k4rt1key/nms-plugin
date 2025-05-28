package utils

import "time"

func GetIndianTime() time.Time {

	loc, _ := time.LoadLocation("Asia/Kolkata")

	return time.Now().In(loc)

}
