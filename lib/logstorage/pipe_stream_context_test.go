package logstorage

import (
	"testing"
)

func TestParsePipeStreamContextSuccess(t *testing.T) {
	f := func(pipeStr string) {
		t.Helper()
		expectParsePipeSuccess(t, pipeStr)
	}

	f(`stream_context before 5`)
	f(`stream_context after 10`)
	f(`stream_context after 0`)
	f(`stream_context before 10 after 20`)
	f(`stream_context after 1 time_window 2h30m`)
	f(`stream_context before 1 time_window 2h30m`)
	f(`stream_context before 1 after 3 time_window 2h30m`)
}

func TestParsePipeStreamContextFailure(t *testing.T) {
	f := func(pipeStr string) {
		t.Helper()
		expectParsePipeFailure(t, pipeStr)
	}

	f(`stream_context`)
	f(`stream_context before`)
	f(`stream_context after`)
	f(`stream_context before after`)
	f(`stream_context after before`)
	f(`stream_context before -4`)
	f(`stream_context after -4`)
	f(`stream_context time_window`)
	f(`stream_context before 3 time_window`)
	f(`stream_context before 3 time_window foobar`)
}

func TestPipeStreamContextUpdateNeededFields(t *testing.T) {
	f := func(s, allowFilters, denyFilters, allowFiltersExpected, denyFiltersExpected string) {
		t.Helper()
		expectPipeNeededFields(t, s, allowFilters, denyFilters, allowFiltersExpected, denyFiltersExpected)
	}

	// all the needed fields
	f("stream_context after 4 before 10", "*", "", "*", "")

	// plus unneeded fields
	f("stream_context before 10 after 4", "*", "f1,f2", "*", "f1,f2")
	f("stream_context after 4", "*", "_time,f1,_stream_id", "*", "f1")

	// needed fields
	f("stream_context before 3", "f1,f2", "", "_stream_id,_time,f1,f2", "")
	f("stream_context before 3", "_time,f1,_stream_id", "", "_stream_id,_time,f1", "")
}
