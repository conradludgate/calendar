package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
)

var u string

func init() {
	fmt.Println("Init dev tools")

	u = strconv.FormatInt(int64(time.Since(time.Unix(0, 0))), 16)
	http.HandleFunc("/uuid", handleUUID)
	http.HandleFunc("/dev", handleDev)
	http.HandleFunc("/restart", handleRestart)
}

func handleRestart(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Restart signal received")
	w.Write([]byte("Stopping..."))
	os.Exit(1)
}

func handleUUID(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(u))
}

func handleDev(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`
<!DOCTYPE html>
<head>
</head>
<body>
	<input id="date" type="date" onchange="updateImage()" value="2019-12-19" /><br/>
	<object width="700" height="700" id="image"></object>
	<script>
		let uuid = "";

		function updateImage() {
			let date = document.querySelector("#date").value;
			document.querySelector("#image").data="/calendar?emoji=apple&date=" + date + "&uuid=" + uuid;
		}

		setInterval(() => fetch("/uuid").then(resp => resp.text()).then(new_uuid => {
			if (new_uuid !== uuid) {
				uuid = new_uuid;
				updateImage();
			}
		}), 100);
	</script>
</body>
	`))
}
