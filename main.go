package main

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	svg "github.com/ajstarks/svgo"
)

func main() {
	http.HandleFunc("/calendar", handleEmoji)

	fmt.Println("Loaded")
	http.ListenAndServe(":5997", nil)
}

func handleEmoji(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL)

	typ := r.FormValue("emoji")
	date := r.FormValue("date")

	if typ == "apple" {
		if err := makeAppleEmoji(w, date); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Date is not supported. Error %s", err.Error())
		}
		return
	}
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("Emoji type is not supported"))
}

func pos(n int, offset int) (x int, y int) {
	x, y = 63, 22
	pos := n + offset - 1
	x += 5 * (pos % 7)
	y += 3 * (pos / 7)
	return
}

func parseDate(date string) (month string, day int, days int, offset int, err error) {
	formats := []string{
		"2006-01-02",
		"01-02",
		"Jan 02",
		"Jan 2",
	}

	date = strings.Replace(date, "+", " ", -1)

	var t time.Time
	fail := true
	for _, format := range formats {
		t, err = time.Parse(format, date)
		if err == nil {
			fail = false
			break
		}
	}
	if fail {
		return
	}

	if t.Year() == 0 {
		t = t.AddDate(time.Now().Year(), 0, 0)
	}

	day = t.Day()
	month = strings.ToUpper(t.Month().String()[:3])
	year := t.Year()
	offset = (int(t.Weekday()) + 35 - day) % 7
	days = []int{31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}[t.Month()-1]
	if days == 28 && year%4 == 0 && (year%100 != 0 || year%400 == 0) {
		days++
	}
	return
}

func makeAppleEmoji(w io.Writer, date string) error {
	month, day, days, offset, err := parseDate(date)
	if err != nil {
		return err
	}

	width := 100
	height := 100
	canvas := svg.New(w)
	canvas.Start(width, height)

	topRectPath := "M 0 0 l 0 38 l 100 0 l 0 -38"
	topRectPath += "l -27 0 l -0.5 1 l 1 1 l -1.5 0.5 l 0.5 1.5 a 6 6 0 1 1 -3 3 l -1 -1 l -2 1 l 1.5 -1.5 l -1 -2 l -3 0 l 3 -1 l 0.5 -2.5"
	topRectPath += "L 24 0 l 0 4 a 6 6 0 1 1 -4 0 l 0 -4"

	canvas.Path(topRectPath, "fill:rgb(169,95,95)")
	canvas.Rect(0, 38, 100, 72, "fill:rgb(220,220,220)")

	canvas.Text(5, 34, month, "font-family:'Super Rad', sans-serif;fill:rgb(225,225,225);font-size:20px")

	for i := 1; i <= days; i++ {
		x, y := pos(i, offset)
		canvas.Text(x, y, strconv.Itoa(i), "font-family:'Super Rad', sans-serif; fill:rgb(225,225,225); font-size: 2px;text-anchor:end")
	}

	canvas.Text(50, 85, strconv.Itoa(day), "font-family:sans-serif;fill:rgb(0,0,0);font-size:50px;text-anchor:middle")

	canvas.End()

	return nil
}
