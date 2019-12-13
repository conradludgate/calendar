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
	http.ListenAndServe(":8080", nil)
}

func handleEmoji(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL)

	typ := r.FormValue("emoji")
	day, err := strconv.Atoi(r.FormValue("day"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Day was not a valid number"))
		return
	}
	month, err := strconv.Atoi(r.FormValue("month"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Month was not a valid number"))
		return
	}
	year, err := strconv.Atoi(r.FormValue("year"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Year was not a valid number"))
		return
	}
	if year == 0 {
		year = time.Now().Year()
	}

	if typ == "apple" {
		makeAppleEmoji(w, day, time.Month(month), year)
		return
	}
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("Type is not supported"))
}

func pos(n int, offset int) (x int, y int) {
	x, y = 63, 22
	pos := n + offset - 1
	x += 5 * (pos % 7)
	y += 3 * (pos / 7)
	return
}

func parseDate(m time.Month, day int, year int) (month string, days int, offset int) {
	month = strings.ToUpper(m.String()[:3])
	t := time.Date(year, m, day, 0, 0, 0, 0, time.UTC)
	offset = (int(t.Weekday()) + 35 - day) % 7
	days = []int{31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}[m-1]
	if days == 28 && year%4 == 0 && (year%100 != 0 || year%400 == 0) {
		days++
	}
	return
}

func makeAppleEmoji(w io.Writer, day int, mon time.Month, year int) {
	width := 100
	height := 100
	canvas := svg.New(w)
	canvas.Start(width, height)

	topRectPath := "M 0 0 l 0 38 l 100 0 l 0 -38"
	topRectPath += "l -27 0 l -0.5 1 l 1 1 l -1.5 0.5 l 0.5 1.5 a 6 6 0 1 1 -3 3 l -1 -1 l -2 1 l 1.5 -1.5 l -1 -2 l -3 0 l 3 -1 l 0.5 -2.5"
	topRectPath += "L 24 0 l 0 4 a 6 6 0 1 1 -4 0 l 0 -4"

	canvas.Path(topRectPath, "fill:rgb(169,95,95)")
	canvas.Rect(0, 38, 100, 72, "fill:rgb(220,220,220)")

	month, days, offset := parseDate(mon, day, year)

	canvas.Text(5, 34, month, "font-family:'Super Rad', sans-serif;fill:rgb(225,225,225);font-size:20px")

	for i := 1; i <= days; i++ {
		x, y := pos(i, offset)
		canvas.Text(x, y, strconv.Itoa(i), "font-family:'Super Rad', sans-serif; fill:rgb(225,225,225); font-size: 2px;text-anchor:end")
	}

	canvas.Text(50, 85, strconv.Itoa(day), "font-family:sans-serif;fill:rgb(0,0,0);font-size:50px;text-anchor:middle")

	canvas.End()
}
