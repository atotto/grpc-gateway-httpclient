package plugin

import (
	"testing"
)

func TestPattern(t *testing.T) {
	str, err := ParsePattern(`/message/{message_id}`)
	if err != nil {
		t.Fatal(err)
	}

	if `"/message/"+ req.MessageId+ ""` != str {
		t.Fatalf("invalid pattern: %s", str)
	}
}
