package logstorage

import (
	"testing"
)

func TestParsePipeUniqSuccess(t *testing.T) {
	f := func(pipeStr string) {
		t.Helper()
		expectParsePipeSuccess(t, pipeStr)
	}

	f(`uniq by (x)`)
	f(`uniq by (x) limit 10`)
	f(`uniq by (x, y)`)
	f(`uniq by (x, y) with hits`)
	f(`uniq by (x, y) limit 10`)
	f(`uniq by (x, y) with hits limit 10`)
}

func TestParsePipeUniqFailure(t *testing.T) {
	f := func(pipeStr string) {
		t.Helper()
		expectParsePipeFailure(t, pipeStr)
	}

	f(`uniq`)
	f(`uniq hits`)
	f(`uniq limit`)
	f(`uniq by ()`)
	f(`uniq by (*)`)
	f(`uniq by (a*)`)
	f(`uniq by`)
	f(`uniq by hits`)
	f(`uniq by foo bar`)
	f(`uniq foo bar`)
	f(`uniq by(x) limit`)
	f(`uniq by(x) limit foo`)
	f(`uniq by (x) with`)
}

func TestPipeUniq(t *testing.T) {
	f := func(pipeStr string, rows, rowsExpected [][]Field) {
		t.Helper()
		expectPipeResults(t, pipeStr, rows, rowsExpected)
	}

	f("uniq by (a)", [][]Field{
		{
			{"a", `2`},
			{"b", `3`},
		},
		{
			{"a", "2"},
			{"b", "3"},
		},
		{
			{"a", `2`},
			{"b", `54`},
			{"c", "d"},
		},
	}, [][]Field{
		{
			{"a", "2"},
		},
	})

	f("uniq by a hits", [][]Field{
		{
			{"a", `2`},
			{"b", `3`},
		},
		{
			{"a", "2"},
			{"b", "3"},
		},
		{
			{"a", `2`},
			{"b", `54`},
			{"c", "d"},
		},
	}, [][]Field{
		{
			{"a", "2"},
			{"hits", "3"},
		},
	})

	f("uniq b", [][]Field{
		{
			{"a", `2`},
			{"b", `3`},
		},
		{
			{"a", "2"},
			{"b", "3"},
		},
		{
			{"a", `2`},
			{"b", `54`},
			{"c", "d"},
		},
	}, [][]Field{
		{
			{"b", "3"},
		},
		{
			{"b", "54"},
		},
	})

	f("uniq by (b) hits", [][]Field{
		{
			{"a", `2`},
			{"b", `3`},
		},
		{
			{"a", "2"},
			{"b", "3"},
		},
		{
			{"a", `2`},
			{"b", `54`},
			{"c", "d"},
		},
	}, [][]Field{
		{
			{"b", "3"},
			{"hits", "2"},
		},
		{
			{"b", "54"},
			{"hits", "1"},
		},
	})

	f("uniq by (c)", [][]Field{
		{
			{"a", `2`},
			{"b", `3`},
		},
		{
			{"a", "2"},
			{"b", "3"},
		},
		{
			{"a", `2`},
			{"b", `54`},
			{"c", "d"},
		},
	}, [][]Field{
		{
			{"c", ""},
		},
		{
			{"c", "d"},
		},
	})

	f("uniq by (c) hits", [][]Field{
		{
			{"a", `2`},
			{"b", `3`},
		},
		{
			{"a", "2"},
			{"b", "3"},
		},
		{
			{"a", `2`},
			{"b", `54`},
			{"c", "d"},
		},
	}, [][]Field{
		{
			{"c", ""},
			{"hits", "2"},
		},
		{
			{"c", "d"},
			{"hits", "1"},
		},
	})

	f("uniq by (d)", [][]Field{
		{
			{"a", `2`},
			{"b", `3`},
		},
		{
			{"a", "2"},
			{"b", "3"},
		},
		{
			{"a", `2`},
			{"b", `54`},
			{"c", "d"},
		},
	}, [][]Field{
		{
			{"d", ""},
		},
	})

	f("uniq by (d) hits", [][]Field{
		{
			{"a", `2`},
			{"b", `3`},
		},
		{
			{"a", "2"},
			{"b", "3"},
		},
		{
			{"a", `2`},
			{"b", `54`},
			{"c", "d"},
		},
	}, [][]Field{
		{
			{"d", ""},
			{"hits", "3"},
		},
	})

	f("uniq by (a, b)", [][]Field{
		{
			{"a", `2`},
			{"b", `3`},
		},
		{
			{"a", "2"},
			{"b", "3"},
		},
		{
			{"a", `2`},
			{"b", `54`},
			{"c", "d"},
		},
	}, [][]Field{
		{
			{"a", "2"},
			{"b", "3"},
		},
		{
			{"a", "2"},
			{"b", "54"},
		},
	})

	f("uniq a, b hits", [][]Field{
		{
			{"a", `2`},
			{"b", `3`},
		},
		{
			{"a", "2"},
			{"b", "3"},
		},
		{
			{"a", `2`},
			{"b", `54`},
			{"c", "d"},
		},
	}, [][]Field{
		{
			{"a", "2"},
			{"b", "3"},
			{"hits", "2"},
		},
		{
			{"a", "2"},
			{"b", "54"},
			{"hits", "1"},
		},
	})
}

func TestPipeUniqUpdateNeededFields(t *testing.T) {
	f := func(s, allowFilters, denyFilters, allowFiltersExpected, denyFiltersExpected string) {
		t.Helper()
		expectPipeNeededFields(t, s, allowFilters, denyFilters, allowFiltersExpected, denyFiltersExpected)
	}

	// all the needed fields
	f("uniq by(f1,f2)", "*", "", "f1,f2", "")
	f("uniq by(f1,f2) with hits", "*", "", "f1,f2", "")

	// all the needed fields, unneeded fields do not intersect with src
	f("uniq by(s1, s2)", "*", "f1,f2", "s1,s2", "")
	f("uniq by(s1, s2)", "*", "f*", "s1,s2", "")

	// all the needed fields, unneeded fields intersect with src
	f("uniq by(s1, s2)", "*", "s1,f1,f2", "s1,s2", "")
	f("uniq by(s1, s2)", "*", "s1,s2,f1", "s1,s2", "")
	f("uniq by(s1, s2)", "*", "s*,f*", "s1,s2", "")

	// needed fields do not intersect with src
	f("uniq by (s1, s2)", "f1,f2", "", "s1,s2", "")
	f("uniq by (s1, s2)", "f*", "", "s1,s2", "")

	// needed fields intersect with src
	f("uniq by (s1, s2)", "s1,f1,f2", "", "s1,s2", "")
	f("uniq by (s1, s2)", "s*,f*", "", "s1,s2", "")
}
