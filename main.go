package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	svg "github.com/ajstarks/svgo"
)

func main() {
	http.HandleFunc("/calendar", handleEmoji)

	fmt.Println("Loaded")
	log.Fatalln(http.ListenAndServe(":5997", nil))
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
	pos := n + offset - 1
	x = (pos % 7) * 20
	y = (pos / 7) * 12
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

	canvas.Startpercent(100, 100, "viewBox='0 0 200 200'", "preserveAspectRatio='none'")

	canvas.Gid("background")

	topRectPath := Path{}.MoveTo(0, 0).Ver(76).Hor(200).Ver(-76)
	topRectPath = topRectPath.Hor(-54).Line(-1, 2).Line(2, 2).Line(-3, 1).Line(1, 3)
	topRectPath = topRectPath.Arc(12, 12, -6, 6, -45, true, true)
	topRectPath = topRectPath.Line(-2, -2).Line(-4, 2).Line(3, -3).Line(-2, -4).Line(-6, 0).Line(6, -2).Line(1, -5)
	topRectPath = topRectPath.HorTo(48).Ver(8).Arc(12, 12, -8, 0, 0, true, true).Ver(-8)

	canvas.Path(topRectPath.String(), "fill:rgb(169,95,95)")
	canvas.Rect(0, 76, 200, 144, "fill:rgb(220,220,220)")

	canvas.Gend()

	canvas.Text(10, 68, month, "font-family:'Super Rad', sans-serif;fill:rgb(225,225,225);font-size:40px")

	canvas.Group("id='month_view'", "transform='translate(126, 42) scale(0.5)'")
	for i := 1; i <= days; i++ {
		x, y := pos(i, offset)
		canvas.Text(x, y, strconv.Itoa(i), "font-family:'Super Rad', sans-serif; fill:rgb(225,225,225); font-size: 8px;text-anchor:end")
	}
	canvas.Gend()

	canvas.Text(100, 170, strconv.Itoa(day), "font-family:sans-serif;fill:rgb(0,0,0);font-size:100px;text-anchor:middle")

	canvas.End()

	return nil
}
