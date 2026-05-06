# matcha

<img width="653" height="429" alt="image" src="https://github.com/user-attachments/assets/1c38344e-32f8-4075-bb14-84d6ff3e180e" />

A wizard that generates a personal site shaped like a command-line shell. Pure HTML / CSS / JS / jQuery output. Zero build step. Deploy by uploading.

```
you@host[~]:  $ help
blog [slug|n]    list blog posts, or read one
contact          show contact information
fortune          print a random fortune
help             show this list
myself           about me
theme [name]     cycle theme; or `theme <name>`, `theme list`
weather <city>   current weather for a city
...
```

## Why

Personal sites are getting templated, vibe-coded, hostable-on-Vercel/Netlify, generic. This is the opposite. Every command is intentional, the output is one folder of static files you understand and own, and the result has personality. Plus, everyone wanted a website like my [personal one](https://drcat.fun/)!

## Install

Download the binary for your OS from the [releases page](https://github.com/catsmells/matcha/releases), or build from source:

```sh
git clone https://github.com/catsmells/matcha
cd matcha
go build -o matcha
```

## Use

```sh
./matcha
```

The wizard asks you a few questions and writes a finished site to `./my-site/` (or wherever you point it). Upload the contents of that folder to your web host's public directory.

## What You Get

- `index.html` - The shell, single file, ~30KB
- `blog/` - Drop new posts here, edit `posts.json`
- `projects/` - Same scaffold for projects
- `404.html`, `403.html`, etc. - Themed error pages
- `.htaccess` - Wires error pages on Apache/LiteSpeed hosts
- `webring.json` - List of personal sites you link to

Every part is editable plain text. Re-run the wizard if you want to regenerate from scratch with different options, or just edit the output.

## Features Available in the Wizard

**Sections:** blog (with Atom feed builder), projects, reading log, webring.

**Utility Commands:** weather, ip, define, qr, hex/rgb color preview, stock quotes, rss reader.

**Flavor Commands:** fortune, cowsay, sudo joke, neofetch, lichess daily chess puzzle, wordle.

**Themes:** light/dark, japanese (paper/sakura/matcha/sumi), CRT (amber/phosphor), dev classics (solarized/nord).

**Blog Social Share Buttons:** copy link, twitter, mastodon, linkedin, email.

## Built Using

- [huh](https://github.com/charmbracelet/huh) for the TUI
- jQuery in the generated output (yes really, it works fine and is small)
- `go:embed` for the templates

## Inspiration

The original implementation lives at [drcat.fun](https://drcat.fun). This generator extracts the reusable parts so others can have their own.

## License

MIT. Do whatever you want with it, but contact me if you'd like to use it commercially.

## Contributing

Issues and PRs welcome. The codebase is small enough to read in an afternoon. See `docs/developing.md` for how the templating works.

## Cool Bits

- I regularly browse the web development threads on 4chan and Lainchan. If I see a cool site using matcha that attributes my site, I'll mention it on the main repository here!
- Theme bundles I find cool or interesting will be linked on this main repository. I love seeing cool color combinations!
