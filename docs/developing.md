# Developing matcha

How matcha is built and how to contribute.

## Architecture

matcha is a Go CLI that generates static sites. There's no server
component, no runtime. Once a user runs the wizard, the output is
plain HTML/CSS/JS with no dependency on matcha itself.

The Pieces:

- **`main.go`** - the wizard. Builds a `Config` from user answers via
  charmbracelet/huh, then runs `generate(cfg)` to produce the output
  directory.

- **`templates/`** - embedded into the binary via `go:embed`. Contains
  the engine template (`index.html.tmpl`), partial command files, CSS
  files, error pages, and starter content.

- **The Output** - a folder with `index.html`, error pages, an
  `.htaccess`, optional `blog/` and `projects/` directories, and an
  optional `webring.json`.

## How Generation Works

`generate()` does string substitution and conditional file
concatenation. There's no template engine - just `strings.ReplaceAll`
on `__PLACEHOLDER__` markers. I avoided Go's `text/template` because
the output is HTML/CSS/JS that has its own `{{` and `}}` syntax in
places, and escaping those would have been miserable.

The Flow:

1. Read `templates/index.html.tmpl` into memory.
2. Always concat `templates/styles/_core.css` and
   `templates/commands/_core.js`.
3. For each command the user selects, look for matching files in
   `templates/styles/<name>.css`, `templates/commands/<name>.js`,
   `templates/help/<name>.txt`, and `templates/man/<name>.txt`. Append
   if they exist; skip silently if they don't.
4. Build the THEMES JS object from `cfg.Themes` selections.
5. Run `strings.ReplaceAll` for each `__PLACEHOLDER__` against the
   accumulated content.
6. Write the result to `<outDir>/index.html`.
7. Copy error pages, `.htaccess`, and starter content directories,
   running the same replacements on them.

Placeholders in `index.html.tmpl`:

| Marker | Source |
|---|---|
| `__USER__`, `__HOST__`, `__DIR__` | wizard input |
| `__DISPLAY_NAME__`, `__SITE_URL__`, `__EMAIL__` | wizard input |
| `__SHARE_HANDLE__` | derived: "via {name} ({url})" |
| `__STYLES__` | concatenated CSS files |
| `__COMMANDS__` | concatenated JS command files |
| `__HELP_ROWS__` | help-row tuples assembled from `help/*.txt` |
| `__MAN_ENTRIES__` | manpage entries assembled from `man/*.txt` |
| `__THEMES__` | JS object literal built by `buildThemes()` |
| `__SHARE_BLOCK__` | `true` or `false` for share button rendering |
| `__BIO__` | JS-array string of bio lines built from cfg |
| `__ATTRIBUTION__` | hardcoded matcha + dr. cat link |

## Adding a New Optional Command

Most contributions will be new commands. The pattern is:

1. Drop the implementation in `templates/commands/<name>.js`. It
   should be a single command function (or a few related ones) that
   slots into the `const commands = { ... }` object. The output engine
   provides `print()`, `printText()`, `blank()`, `escapeHtml()`, plus
   access to `USER`, `HOST`, `DIR`, `THEMES`, and similar globals
   declared near the top of the template.

   ```javascript
       async joke() {
         try {
           const res = await fetch('https://icanhazdadjoke.com/', {
             headers: { 'Accept': 'application/json' }
           });
           const data = await res.json();
           printText(data.joke);
         } catch (e) { printText('joke: failed to fetch.'); }
       },
   ```

   Note the leading 4-space indent and trailing comma - your file is
   inserted directly into the JS object literal, so it has to be valid
   in that context.

2. Add a help row in `templates/help/<name>.txt`. Just one line, a
   JS array literal of [name, description]:

   ```
   ['joke', 'tell a dad joke']
   ```

3. Optional: add a `man` page in `templates/man/<name>.txt`:

   ```
   JOKE(1)\n\nNAME\n    joke - random dad joke\n\nSYNOPSIS\n    joke
   ```

   Use literal `\n` in the file - the Go side reads it as a string and
   `\n` becomes a newline at output time. (Yes, double-escaping. I
   could use a multi-line text format instead. PRs are welcome if you
   want to clean that up.)

4. Add CSS if needed in `templates/styles/<name>.css`. Most
   commands don't need any. The wordle and chess commands do.

5. Wire it into the wizard form in `main.go`. Find the right
   `huh.NewMultiSelect` group (utility, flavor, or sections) and add
   an option:

   ```go
   huh.NewOption("joke (random dad joke)", "joke"),
   ```

   The option's value (`"joke"` here) must match the filename stem
   you used in `templates/`.

The wizard will pick it up automatically because `generate()` 
iterates the user's selections and looks for matching files.

## Adding a New Theme Bundle

Edit `buildThemes()` in `main.go`. The function takes the user's
selected bundle keys (like `"japanese"` or `"crt"`) and returns a JS
object literal as a string.

To add a new bundle:

```go
if contains(selected, "ocean") {
    themes = append(themes,
        `    oceanic:  { bg: '#1B2B34', fg: '#C0C5CE' },`,
        `    abyss:    { bg: '#0B0E14', fg: '#7FDBFF' },`,
    )
}
```

Then add it to the wizard form's themes group. I'll link
cool-looking bundles on the main Github repo for matcha.

## Testing Locally

There's no test harness yet. The fastest way to iterate:

```sh
go run . 
# answer prompts
cd your-site
python3 -m http.server 8000
# open http://localhost:8000
```

`python3 -m http.server` is important - opening `index.html` directly
with `file://` will break the `fetch()` calls for blog and projects
content because of CORS.

## Cross-Compiling for Releases

```sh
GOOS=darwin  GOARCH=arm64  go build -o matcha-darwin-arm64
GOOS=darwin  GOARCH=amd64  go build -o matcha-darwin-amd64
GOOS=linux   GOARCH=amd64  go build -o matcha-linux-amd64
GOOS=windows GOARCH=amd64  go build -o matcha-windows-amd64.exe
```

Or just push a tag matching `v*` and let `.github/workflows/release.yml`
handle it. Be as lazy as possible about it.

## Things I Deliberately Didn't Build

For context if you're wondering why something seems missing:

- No template engine. Plain string replacement is enough. Anything
  more would fight the JS/HTML output.
- No config file. The wizard is the config. Saving config to a
  file just to re-read it adds a step. Re-running the wizard takes
  thirty seconds. Efficiency!
- No live preview. The output is plain HTML; the user opens it in
  a browser. Building a dev-server mode would 10x the codebase.
- No backend, no comments, no analytics. Personal sites don't
  need them. If you want them, fork off already.

## Style

- Plain Go, stdlib where possible, one external dep (huh).
- Errors return up. Don't `log.Fatal` from inside helpers.
- matcha's source has no comments deliberately - the code is short
  enough to read straight through. The generated output is also
  comment-free for the same reason. PRs adding comments will probably
  be redirected toward making the code clearer instead.

## Questions

Open an issue on https://github.com/catsmells/matcha or drop a note on
my contact page (https://drcat.fun/).
