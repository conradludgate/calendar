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
		"2006 Jan 2",
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

	canvas := svg.New(w)
	canvas.Startpercent(100, 100, "viewBox='0 0 100 100'", "preserveAspectRatio='none'")

	topRectPath := Path{}.MoveTo(0, 0).Ver(38).Hor(100).Ver(-38)
	topRectPath = topRectPath.Hor(-27).Line(-0.5, 1).Line(1, 1).Line(-1.5, 0.5).Line(0.5, 1.5)
	topRectPath = topRectPath.Arc(6, 6, -3, 3, -45, true, true)
	topRectPath = topRectPath.Line(-1, -1).Line(-2, 1).Line(1.5, -1.5).Line(-1, -2).Line(-3, 0).Line(3, -1).Line(0.5, -2.5)
	topRectPath = topRectPath.HorTo(24).Ver(4).Arc(6, 6, -4, 0, 0, true, true).Ver(-4)

	canvas.Path(topRectPath.String(), "fill:rgb(169,95,95)")
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
