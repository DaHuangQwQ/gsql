package valuer

import "testing"

func Test_unsafeValue_SetColumns(t *testing.T) {
	testSetColumns(t, NewUnsafeValue)
}
