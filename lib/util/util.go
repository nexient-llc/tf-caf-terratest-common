package util

import (
	"fmt"
	"testing"
	"time"
)

func WaitUntilCloudResourceReady(t *testing.T, resourceName string, pauseIfFailInSec int, maxAttempts int, pollFunc func() (any, error), checkResult func(any) bool) {
	var counter int
	for {
		if counter > maxAttempts {
			t.Fatal("Timeout while waiting for " + resourceName + " initialized")
			break
		}
		res, err := pollFunc()
		if nil == err && checkResult(res) {
			break
		}
		counter++
		fmt.Println("waiting until " + resourceName + " initialized")
		time.Sleep(time.Duration(pauseIfFailInSec) * time.Second)
	}
}
