package common

import "testing"

func TestGenRandomString(t *testing.T) {
	testLength := 5
	randomString := GenRandomString(testLength)

	if len(randomString) != testLength {
		t.Errorf("length of generated random string should be %d instead of %d", testLength, len(randomString))
	}

}
