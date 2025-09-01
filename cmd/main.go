package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/vharitonsky/iniflags"
	"io"
	"log"
	"net/http"
	"time"
	"wedding/pkg/logflags"
	"wedding/pkg/tglogs"
)

var (
	httpListenAddr    = flag.String("httpListenAddr", ":8080", "HTTP listen address")
	publicStaticFSDir = flag.String("publicStaticFSDir", "", "Public static files directory")
)

func main() {
	iniflags.Parse()
	logflags.LogAllFlags()

	tglogs.Init("Wedding")
	tglogs.Send("Starting application")

	mux := http.NewServeMux()

	// 1) Serve static assets from ./public under /static/*
	//    e.g. /static/app.css, /static/app.js, /static/logo.png
	fs := http.FileServer(http.Dir(*publicStaticFSDir))
	mux.Handle("/", fs)

	// 3) Handle the form submit (POST /submit)
	mux.HandleFunc("/user_account/attendance", attendanceHandler)

	// Optional: small server hardening + request logging
	srv := &http.Server{
		Addr:              *httpListenAddr,
		Handler:           logRequests(mux),
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("Listening on %s\n", *httpListenAddr)
	log.Fatal(srv.ListenAndServe())
}

type AttendanceMessage struct {
	Name string

	AttendanceMapping string
	AttendanceMessage string

	AccommodationMapping string
	AccommodationMessage string

	Comment string
}

var attendanceOptionToMessage = map[string]string{
	"1": "Так, зможу",
	"2": "Вагаюсь з відповіддю, повідомлю пізніше",
	"3": "Не зможу прийти",
}

var accommodationOptionToMessage = map[string]string{
	"1": "Маю де заночувати",
	"2": "Не знаю, мені потрібна допомога",
	"3": "Вагаюсь, повідомлю пізніше",
}

func (a AttendanceMessage) String() string {
	return fmt.Sprintf(`
*Імʼя:*
%s
*Присутність:*
%s
*Проживання:*
%s
*Коментарі:*
%s
`, a.Name, a.AttendanceMessage, a.AccommodationMessage, a.Comment)
}

func NewAttendanceMessage(name, attendanceOption, accommodationOption, comment string) (AttendanceMessage, error) {
	attd, ok := attendanceOptionToMessage[attendanceOption]
	if !ok {
		return AttendanceMessage{}, fmt.Errorf("invalid attendance option: %s", attendanceOption)
	}

	acm, ok := accommodationOptionToMessage[accommodationOption]
	if !ok {
		return AttendanceMessage{}, fmt.Errorf("invalid attendance option: %s", attendanceOption)
	}

	return AttendanceMessage{
		Name: name,

		AttendanceMapping: attendanceOption,
		AttendanceMessage: attd,

		AccommodationMapping: accommodationOption,
		AccommodationMessage: acm,

		Comment: comment,
	}, nil
}

type AttendanceRequest struct {
	Name          string `json:"name"`
	Attendance    string `json:"attendance"`
	Accommodation string `json:"accommodation"`
	Comment       string `json:"comment"`
}

// attendanceHandler reads a simple x-www-form-urlencoded form
func attendanceHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	rawReq, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error creating message: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var req AttendanceRequest

	if err = json.Unmarshal(rawReq, &req); err != nil {
		log.Printf("Error unmarshaling message: %v %s\n", err, rawReq)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Form submit: name=%q attendance=%q accommodation=%q comment=%q", req.Name, req.Attendance, req.Accommodation, req.Comment)

	msg, err := NewAttendanceMessage(req.Name, req.Attendance, req.Accommodation, req.Comment)
	if err != nil {
		log.Printf("Error creating message: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tglogs.Send(msg.String())

	// Redirect with 303 See Other (so browser does GET after POST)
	http.Redirect(w, r, "/?ok=1", http.StatusSeeOther)
}

func logRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s (%s)", r.Method, r.URL.Path, time.Since(start).Round(time.Millisecond))
	})
}
