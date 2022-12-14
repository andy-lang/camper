package camper

import (
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
)

type Track struct {
	Trackno  int
	Title    string
	Duration time.Duration
}

type Release struct {
	Title    string
	Artist   string
	Released time.Time
	Tracks   []Track
	Genres   []string
}

func reverse(s []string) {
	for i := len(s)/2 - 1; i >= 0; i-- {
		opp := len(s) - 1 - i
		s[i], s[opp] = s[opp], s[i]
	}
}

func bandcampTimeToDuration(t string) time.Duration {
	// split on colons (heheheh), convert each to their constituent parts
	splitTimes := strings.Split(t, ":")

	reverse(splitTimes)

	secs, _ := strconv.Atoi(splitTimes[0])
	mins, _ := strconv.Atoi(splitTimes[1])
	hours := 0
	if len(splitTimes) == 3 {
		hours, _ = strconv.Atoi(splitTimes[2])
	}
	return (time.Duration(hours)*time.Hour + time.Duration(mins)*time.Minute) + (time.Duration(secs) * time.Second)
}

func ReleaseFromURL(url string) Release {
	c := colly.NewCollector()
	c.URLFilters = []*regexp.Regexp{regexp.MustCompile(`^https?://[a-z]+\.bandcamp.com/`)}

	release := Release{}

	// populate Track Title, Artist
	c.OnHTML("div#name-section", func(e *colly.HTMLElement) {
		release.Title = strings.TrimSpace(e.ChildText("h2.trackTitle"))
		release.Artist = strings.TrimSpace(e.ChildText("a[href]"))
	})

	// Populate any track entries
	c.OnHTML("tr.track_row_view", func(e *colly.HTMLElement) {
		tracknum, _ := strconv.Atoi(strings.Split(e.ChildText("div.track_number"), ".")[0])
		title := e.ChildText(".track-title")

        // duration is not set for some tracks, e.g. if it's an upcoming album and only certain songs are previewed. Make a best effort guess.
        duration := time.Duration(0)
        if len(e.ChildText("span.time")) > 0 {
            duration = bandcampTimeToDuration(e.ChildText("span.time"))
        }

		track := Track{Trackno: tracknum, Title: title, Duration: duration}

		release.Tracks = append(release.Tracks, track)
	})

	/// Populate release date
	c.OnHTML("div.tralbum-credits", func(e *colly.HTMLElement) {
		// grab the first line of the entry - this will always be "release(s|d) Month Day, Year"
		released := strings.Split(strings.TrimSpace(e.Text), "\n")[0]

		// grab everything after the first word, join into a new string
		released = strings.Join(strings.Fields(released)[1:], " ")

		// parse as a datetime
		releaseDate, _ := time.Parse("January 2, 2006", released)
		release.Released = releaseDate
	})

	// Populate genres
	c.OnHTML("div.tralbumData.tralbum-tags", func(e *colly.HTMLElement) {
		tags := e.ChildTexts("a.tag")
		release.Genres = tags
	})

	c.OnScraped(func(r *colly.Response) {
		sort.SliceStable(release.Tracks, func(i, j int) bool {
			return release.Tracks[i].Trackno < release.Tracks[j].Trackno
		})
	})

	c.Visit(url)

	return release
}
