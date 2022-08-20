package camper

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
)

type Track struct {
	trackno  int
	title    string
	duration time.Duration
}

type Release struct {
	title    string
	released time.Time
	tracks   []Track
	genres   []string
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
	fmt.Println(splitTimes)

	reverse(splitTimes)

	fmt.Println(splitTimes)

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

	c.OnHTML("h2.trackTitle", func(e *colly.HTMLElement) {
		release.title = strings.TrimSpace(e.Text)
	})

	c.OnHTML("tr.track_row_view", func(e *colly.HTMLElement) {
		tracknum, _ := strconv.Atoi(strings.Split(e.ChildText("div.track_number"), ".")[0])
		title := e.ChildText(".track-title")
		duration := bandcampTimeToDuration(e.ChildText("span.time"))
		track := Track{trackno: tracknum, title: title, duration: duration}

		fmt.Printf("duration: %v", duration.Seconds())

		release.tracks = append(release.tracks, track)
	})

	c.OnHTML("div.tralbum-credits", func(e *colly.HTMLElement) {
		// grab the first line of the entry - this will always be "release(s|d) Month Day, Year"
		released := strings.Split(strings.TrimSpace(e.Text), "\n")[0]

		// grab everything after the first word, join into a new string
		released = strings.Join(strings.Fields(released)[1:], " ")

		// parse as a datetime
		releaseDate, _ := time.Parse("January 2, 2006", released)
		release.released = releaseDate
	})

	c.OnHTML("div.tralbumData.tralbum-tags", func(e *colly.HTMLElement) {
		tags := e.ChildTexts("a.tag")
		release.genres = tags
	})

	c.Visit(url)

	return release
}
