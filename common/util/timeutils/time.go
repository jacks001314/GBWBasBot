package timeutils

import "time"

func GetNowTimeMS() uint64 {

	return uint64((time.Now().UnixNano()) / (1000 * 1000))
}
