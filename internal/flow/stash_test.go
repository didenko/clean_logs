package flow

import (
	"testing"
)

var rnd_str = []string{
	"h^a*Hp99vyC}E$3aV8h6",
	"@8LV,/c$4vwds47jf2Nh",
	"paHb/z=@mzm22M#T8r66",
	"",
	"&g6Yt8?aGds4YA=23N)P",
	"Ct33n7P[M6FeYMbi&8].",
	"U@D)N9x3d{G8M8bEr,8e",
	"M=k4NTe2(vL8T8+3A?Zf",
	"9j$7R*WpLP4gQ3qK&7)B",
	"",
	"h^a*Hp99vyC}E$3aV8h6@8LV,/c$4vwds47jf2NhpaHb/z=@mzm22M#T8r66&g6Yt8?aGds4YA=23N)PCt33n7P[M6FeYMbi&8].U@D)N9x3d{G8M8bEr,8e,M=k4NTe2(vL8T8+3A?Zf9j$7R*WpLP4gQ3qK&7)B",
}

func TestStash(t *testing.T) {

	var (
		in   = make(chan string)
		err  = make(chan error)
		out  = make(chan string)
		done = make(chan string)
	)

	go func() {
		for e := range err {
			t.Error(e)
		}
	}()

	go Stash(in, done, out, err)

	for _, l := range rnd_str {
		in <- l
	}

	close(done)

	close(in)

	for i, s := range rnd_str {
		replay := <-out
		if s != replay {
			t.Error("Line ", i, " differs in the replay")
		}
	}

	if _, ok := <-out; ok {
		t.Error("Replay channel produced extraneous values")
	}
}

func TestFork(t *testing.T) {
	var (
		in   = make(chan string)
		out1 = make(chan string)
		out2 = make(chan string)
	)

	go Fork(in, out1, out2)

	go func() {
		defer close(in)
		for _, l := range rnd_str {
			in <- l
		}
	}()

	for i, s := range rnd_str {
		r1 := <-out1
		if r1 != s {
			t.Error("Line ", i, "of the first output differs in the replay")
		}
		r2 := <-out2
		if r2 != s {
			t.Error("Line ", i, "of the second output differs in the replay")
		}
	}

	if _, ok := <-out1; ok {
		t.Error("First replay channel produced extraneous values")
	}

	if _, ok := <-out2; ok {
		t.Error("Second replay channel produced extraneous values")
	}
}
