package mergesort

import "testing"

var tests = []struct {
	initial, want []byte
}{
	{initial: []byte{1, 2}, want: []byte{1, 2}},
	{initial: []byte{2, 1}, want: []byte{1, 2}},
	{initial: []byte{3, 2, 1}, want: []byte{1, 2, 3}},
	{initial: []byte{7, 9, 5, 6, 4, 8, 2, 3, 1}, want: []byte{1, 2, 3, 4, 5, 6, 7, 8, 9}},
	{
		initial: []byte{7, 9, 5, 6, 4, 8, 2, 3, 1, 3, 5, 6, 8, 9, 7, 4, 5, 6, 2, 3, 5, 9, 7, 5, 6, 4, 2, 3, 1, 56, 7, 8, 99, 5, 3, 1, 8, 48, 4, 6, 6, 21, 3, 6, 4, 8, 9, 4, 6, 7, 8},
		want:    []byte{1, 1, 1, 2, 2, 2, 3, 3, 3, 3, 3, 3, 4, 4, 4, 4, 4, 4, 5, 5, 5, 5, 5, 5, 6, 6, 6, 6, 6, 6, 6, 6, 7, 7, 7, 7, 7, 8, 8, 8, 8, 8, 8, 9, 9, 9, 9, 21, 48, 56, 99}},

	{initial: []byte{'a', 'b'}, want: []byte{'a', 'b'}},
	{initial: []byte{'m', 'z', 'a'}, want: []byte{'a', 'm', 'z'}},
	{initial: []byte("maxence"), want: []byte("aceemnx")},
	{initial: []byte("hello, mister gopher!"), want: []byte("  !,eeeghhillmooprrst")},
}

func TestSort(t *testing.T) {
	for _, test := range tests {
		got := Sort(test.initial)
		for i, e := range got {
			if e != test.want[i] {
				t.Fatalf("got %v, want %v\n    error was : %v != %v", got, test.want, e, test.want[i])
			}
		}
	}
}
