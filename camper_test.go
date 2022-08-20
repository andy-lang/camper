package camper

import (
	"testing"
	"time"

	"golang.org/x/exp/slices"
)

func TestPoolKids(t *testing.T) {
	url := "https://poolkidsband.bandcamp.com/album/pool-kids"
	release := releaseFromURL(url)
	wantTitle := "Pool Kids"
	if release.title != wantTitle {
		t.Fatalf("got %s, wanted %s", release.title, wantTitle)
	}

	wantReleased := time.Date(2022, 07, 22, 0, 0, 0, 0, time.UTC)
	if release.released != wantReleased {
		t.Fatalf("got release date %v, wanted %v", release.released, wantReleased)
	}

	if len(release.tracks) != 12 {
		t.Fatalf("expected 12 tracks")
	}

	if !slices.Contains(release.genres, "Tallahassee") {
		t.Fatalf("Genres does not appear to be fully populated - %v", release.genres)
	}

	// check the second track, verify it's all good
	if release.tracks[1].title != "That's Physics, Baby" {
		t.Fatalf("Second track doesn't exist or is in incorrect position")
	}

	if release.tracks[1].trackno != 2 {
		t.Fatalf("Bad track number for track #2")
	}

	wantDuration, _ := time.ParseDuration("3m54s")
	if release.tracks[1].duration != wantDuration {
		t.Fatalf("Bad track duration for track #2")
	}
}

func TestLongTrack(t *testing.T) {
	url := "https://vickychow.bandcamp.com/album/tristan-perich-surface-image"
	release := releaseFromURL(url)

	wantDuration, _ := time.ParseDuration("1h3m1s")
	gotDuration := release.tracks[0].duration

	if gotDuration != wantDuration {
		t.Fatalf("Bad track duration, wanted %v, got %v", wantDuration, gotDuration)
	}
}

// Track duration that's longer than an hour
// Track duration that's less than a minute
