package camper

import (
	"testing"
	"time"

	"golang.org/x/exp/slices"
)

func TestPoolKids(t *testing.T) {
	url := "https://poolkidsband.bandcamp.com/album/pool-kids"
	release := ReleaseFromURL(url)
	wantTitle := "Pool Kids"
	if release.Title != wantTitle {
		t.Fatalf("got title %s, wanted %s", release.Title, wantTitle)
	}

	wantArtist := "Pool Kids"
	if release.Artist != wantArtist {
		t.Fatalf("got artist %s, wanted %s", release.Artist, wantArtist)
	}

	wantReleased := time.Date(2022, 07, 22, 0, 0, 0, 0, time.UTC)
	if release.Released != wantReleased {
		t.Fatalf("got release date %v, wanted %v", release.Released, wantReleased)
	}

	if len(release.Tracks) != 12 {
		t.Fatalf("expected 12 tracks")
	}

	if !slices.Contains(release.Genres, "Tallahassee") {
		t.Fatalf("Genres does not appear to be fully populated - %v", release.Genres)
	}

	// check the second track, verify it's all good
	if release.Tracks[1].Title != "That's Physics, Baby" {
		t.Fatalf("Second track doesn't exist or is in incorrect position")
	}

	if release.Tracks[1].Trackno != 2 {
		t.Fatalf("Bad track number for track #2")
	}

	wantDuration, _ := time.ParseDuration("3m54s")
	if release.Tracks[1].Duration != wantDuration {
		t.Fatalf("Bad track duration for track #2")
	}
}

func TestLongTrack(t *testing.T) {
	url := "https://vickychow.bandcamp.com/album/tristan-perich-surface-image"
	release := ReleaseFromURL(url)

	wantDuration, _ := time.ParseDuration("1h3m1s")
	gotDuration := release.Tracks[0].Duration

	if gotDuration != wantDuration {
		t.Fatalf("Bad track duration, wanted %v, got %v", wantDuration, gotDuration)
	}
}

func TestShortTrack(t *testing.T) {
	url := "https://pioulard.bandcamp.com/album/the-beno-t-pioulard-listening-matter-2"
	release := ReleaseFromURL(url)

	wantDuration, _ := time.ParseDuration("44s")
	gotDuration := release.Tracks[3].Duration

	if gotDuration != wantDuration {
		t.Fatalf("Bad track duration, wanted %v, got %v", wantDuration, gotDuration)
	}
}
