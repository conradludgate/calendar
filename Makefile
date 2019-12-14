dev:
	echo "Open http://localhost:5997/dev to get live changes"
	go build
	filewatcher *.go 'go build && curl http://localhost:5997/restart' & while true; do ./calendar; done