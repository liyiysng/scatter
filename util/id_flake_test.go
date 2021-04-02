package util

import (
	"testing"
	"time"
)

func TestIDFlake(t *testing.T) {
	n, err := NewNode(500)
	if err != nil {
		t.Fatal(err)
	}

	id := n.Generate()

	tm := time.Date(2021, 4, 2, 17, 0, 0, 0, time.Local)

	t.Log(tm.UnixNano() / int64(time.Millisecond))

	t.Log(id, time.Now().UnixNano()/int64(time.Millisecond), id.GetUnixMS(), id.GetNodeID(), id.GetStep(), GetUnixFromEpochTime(id.GetEpochMS()), id.GetEpochMS())
}
