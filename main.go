package main
import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"github.com/charmbracelet/huh"
)
var templates embed.FS
type Config struct {
	User           string
	Host           string
	Dir            string
	DisplayName    string
	SiteURL        string
	Email          string
	Bio            string
	Interests      string
	Sections       []string
	UtilityCmds    []string
	FlavorCmds     []string
	Themes         []string
	Share          []string
	OutputDir      string
}
func contains(s []string, v string) bool {
	for _, x := range s {
		if x == v {
			return true
		}
	}
	return false
}
func main() {
	fmt.Println("welcome to matcha setup.")
	fmt.Println()
	cfg := Config{
		User: "drcat", Host: "matcha-host", Dir: "matcha",
		DisplayName: "your name",
		SiteURL:     "https://example.com",
		Email:       "you@example.com",
		Bio:         "tell us about yourself.",
		Interests:   "tech, books, dr. cat, coffee",
		OutputDir:   "./my-site",
	}
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewNote().Title("identity").
				Description("the prompt looks like  user@host[dir]:  $"),
			huh.NewInput().Title("username").Value(&cfg.User),
			huh.NewInput().Title("hostname").
				Description("can be unicode/japanese/cyrillic/anything").Value(&cfg.Host),
			huh.NewInput().Title("directory symbol").Value(&cfg.Dir),
			huh.NewInput().Title("display name").
				Description("used in social share text").Value(&cfg.DisplayName),
			huh.NewInput().Title("site url").Value(&cfg.SiteURL),
			huh.NewInput().Title("email").Value(&cfg.Email),
		),
		huh.NewGroup(
			huh.NewNote().Title("bio").
				Description("shown by the `myself` command. you can edit later."),
			huh.NewText().Title("bio paragraph").
				CharLimit(800).Value(&cfg.Bio),
			huh.NewInput().Title("interests").
				Description("comma-separated").Value(&cfg.Interests),
		),
		huh.NewGroup(
			huh.NewMultiSelect[string]().Title("sections").
				Options(
					huh.NewOption("blog (with rss/atom feed)", "blog").Selected(true),
					huh.NewOption("projects", "projects").Selected(true),
					huh.NewOption("reading log", "reading"),
					huh.NewOption("webring", "webring"),
				).Value(&cfg.Sections),
		),
		huh.NewGroup(
			huh.NewMultiSelect[string]().Title("utility commands").
				Options(
					huh.NewOption("weather <city>", "weather").Selected(true),
					huh.NewOption("ip", "ip").Selected(true),
					huh.NewOption("define <word>", "define").Selected(true),
					huh.NewOption("qr <text>", "qr"),
					huh.NewOption("hex / rgb (color preview)", "color"),
					huh.NewOption("stock <ticker>", "stock"),
					huh.NewOption("rss <feed>", "rss").Selected(true),
				).Value(&cfg.UtilityCmds),
		),
		huh.NewGroup(
			huh.NewMultiSelect[string]().Title("flavor commands").
				Options(
					huh.NewOption("fortune (edit list later)", "fortune").Selected(true),
					huh.NewOption("cowsay <text>", "cowsay").Selected(true),
					huh.NewOption("sudo (joke)", "sudo").Selected(true),
					huh.NewOption("neofetch with ascii cat", "neofetch").Selected(true),
					huh.NewOption("chess (lichess daily puzzle)", "chess"),
					huh.NewOption("wordle", "wordle"),
				).Value(&cfg.FlavorCmds),
		),
		huh.NewGroup(
			huh.NewMultiSelect[string]().Title("themes").
				Options(
					huh.NewOption("light + dark", "basic").Selected(true),
					huh.NewOption("japanese (paper, sakura, matcha, sumi)", "japanese"),
					huh.NewOption("crt (amber, phosphor)", "crt"),
					huh.NewOption("dev classics (solarized, nord)", "dev"),
				).Value(&cfg.Themes),
		),
		huh.NewGroup(
			huh.NewMultiSelect[string]().Title("blog post share buttons").
				Options(
					huh.NewOption("copy link", "copy").Selected(true),
					huh.NewOption("twitter", "twitter").Selected(true),
					huh.NewOption("mastodon", "mastodon").Selected(true),
					huh.NewOption("linkedin", "linkedin"),
					huh.NewOption("email", "email").Selected(true),
				).Value(&cfg.Share),
		),
		huh.NewGroup(
			huh.NewInput().Title("output directory").Value(&cfg.OutputDir),
		),
	)
	if err := form.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	fmt.Println()
	fmt.Println("generating site...")

	if err := generate(cfg); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	fmt.Println()
	fmt.Printf("done. your site is in %s\n", cfg.OutputDir)
	fmt.Println("upload all files to your web host's public folder.")
	fmt.Println("see README.md inside that folder for next steps.")
}
func generate(cfg Config) error {
	out := cfg.OutputDir
	if err := os.MkdirAll(out, 0755); err != nil {
		return err
	}
	indexBytes, err := templates.ReadFile("templates/index.html.tmpl")
	if err != nil {
		return err
	}
	index := string(indexBytes)
	var styles, commands, manEntries, helpRows []string
	coreCSS, _ := templates.ReadFile("templates/styles/_core.css")
	styles = append(styles, string(coreCSS))
	coreJS, _ := templates.ReadFile("templates/commands/_core.js")
	commands = append(commands, string(coreJS))
	allCmds := append(append([]string{}, cfg.UtilityCmds...), cfg.FlavorCmds...)
	allCmds = append(allCmds, cfg.Sections...)
	for _, cmd := range allCmds {
		cssPath := fmt.Sprintf("templates/styles/%s.css", cmd)
		if data, err := templates.ReadFile(cssPath); err == nil {
			styles = append(styles, string(data))
		}
		jsPath := fmt.Sprintf("templates/commands/%s.js", cmd)
		if data, err := templates.ReadFile(jsPath); err == nil {
			commands = append(commands, string(data))
		}
		manPath := fmt.Sprintf("templates/man/%s.txt", cmd)
		if data, err := templates.ReadFile(manPath); err == nil {
			manEntries = append(manEntries, fmt.Sprintf("        %s: %q,", cmd, string(data)))
		}
		helpPath := fmt.Sprintf("templates/help/%s.txt", cmd)
		if data, err := templates.ReadFile(helpPath); err == nil {
			helpRows = append(helpRows, "        "+strings.TrimSpace(string(data))+",")
		}
	}
	themeBlock := buildThemes(cfg.Themes)
	shareBlock := buildShare(cfg.Share)
	bio := buildBio(cfg)
	replacements := map[string]string{
		"__USER__":         cfg.User,
		"__HOST__":         cfg.Host,
		"__DIR__":          cfg.Dir,
		"__DISPLAY_NAME__": cfg.DisplayName,
		"__SITE_URL__":     cfg.SiteURL,
		"__EMAIL__":        cfg.Email,
		"__SHARE_HANDLE__": fmt.Sprintf("via %s (%s)", cfg.DisplayName, cfg.SiteURL),
		"__STYLES__":       strings.Join(styles, "\n"),
		"__COMMANDS__":     strings.Join(commands, "\n"),
		"__MAN_ENTRIES__":  strings.Join(manEntries, "\n"),
		"__HELP_ROWS__":    strings.Join(helpRows, "\n"),
		"__THEMES__":       themeBlock,
		"__SHARE_BLOCK__":  shareBlock,
		"__BIO__":          bio,
		"__ATTRIBUTION__":  `generated with <a href="https://github.com/catsmells/matcha" target="_blank" rel="noopener">matcha</a> by <a href="https://drcat.fun/" target="_blank" rel="noopener">dr. cat</a>`,
	}
	for k, v := range replacements {
		index = strings.ReplaceAll(index, k, v)
	}
	if err := os.WriteFile(filepath.Join(out, "index.html"), []byte(index), 0644); err != nil {
		return err
	}
	staticFiles := []struct{ src, dst string }{
		{"templates/error-pages/404.html", "404.html"},
		{"templates/error-pages/403.html", "403.html"},
		{"templates/error-pages/410.html", "410.html"},
		{"templates/error-pages/500.html", "500.html"},
		{"templates/htaccess", ".htaccess"},
		{"templates/README.md", "README.md"},
	}
	for _, f := range staticFiles {
		data, err := templates.ReadFile(f.src)
		if err != nil {
			continue
		}
		content := string(data)
		for k, v := range replacements {
			content = strings.ReplaceAll(content, k, v)
		}
		if err := os.WriteFile(filepath.Join(out, f.dst), []byte(content), 0644); err != nil {
			return err
		}
	}
	if contains(cfg.Sections, "blog") {
		if err := copyDir("templates/blog", filepath.Join(out, "blog")); err != nil {
			return err
		}
	}
	if contains(cfg.Sections, "projects") {
		if err := copyDir("templates/projects", filepath.Join(out, "projects")); err != nil {
			return err
		}
	}
	if contains(cfg.Sections, "webring") {
		if data, err := templates.ReadFile("templates/webring.json"); err == nil {
			os.WriteFile(filepath.Join(out, "webring.json"), data, 0644)
		}
	}
	return nil
}
func copyDir(src, dst string) error {
	if err := os.MkdirAll(dst, 0755); err != nil {
		return err
	}
	entries, err := fs.ReadDir(templates, src)
	if err != nil {
		return err
	}
	for _, e := range entries {
		data, err := templates.ReadFile(filepath.Join(src, e.Name()))
		if err != nil {
			continue
		}
		if err := os.WriteFile(filepath.Join(dst, e.Name()), data, 0644); err != nil {
			return err
		}
	}
	return nil
}
func buildThemes(selected []string) string {
	themes := []string{
		`    light:    { bg: '#E2E2E2', fg: '#000000' },`,
		`    dark:     { bg: '#000000', fg: '#E2E2E2' },`,
	}
	if contains(selected, "japanese") {
		themes = append(themes,
			`    paper:    { bg: '#F4ECD8', fg: '#3B2F2F' },`,
			`    sakura:   { bg: '#FBE4E4', fg: '#5C2A2A' },`,
			`    matcha:   { bg: '#E5EBD8', fg: '#2C3A1F' },`,
			`    sumi:     { bg: '#F2EFE6', fg: '#1A1A1A' },`,
		)
	}
	if contains(selected, "crt") {
		themes = append(themes,
			`    amber:    { bg: '#1A1A1A', fg: '#FFB000' },`,
			`    phosphor: { bg: '#0A0F0A', fg: '#33FF66' },`,
		)
	}
	if contains(selected, "dev") {
		themes = append(themes,
			`    solarized:{ bg: '#FDF6E3', fg: '#586E75' },`,
			`    nord:     { bg: '#2E3440', fg: '#D8DEE9' },`,
		)
	}
	return "{\n" + strings.Join(themes, "\n") + "\n  }"
}
func buildShare(selected []string) string {
	if len(selected) == 0 {
		return "false"
	}
	return "true"
}
func buildBio(cfg Config) string {
	interestList := ""
	if cfg.Interests != "" {
		interestList = fmt.Sprintf("'<span class=\"dim\">interests ::</span>  %s',", cfg.Interests)
	}
	return fmt.Sprintf(`'<span class="dim">—# myself #——————————————————————</span>',
        '',
        %q,
        '',
        %s`, cfg.Bio, interestList)
}
