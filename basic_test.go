package wits

import (
	"fmt"
	"testing"
)

func TestBasicAuth(t *testing.T) {
	auth := BasicAuth("wendell", "1q1w1e1r")
	fmt.Println(auth)
}
