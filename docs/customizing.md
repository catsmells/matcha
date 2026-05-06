# Customizing Your Generated Site

Your site is plain HTML, CSS, and JavaScript. Open `index.html` in any
text editor and you can change anything. This guide covers the most
common edits.

## Updating Your Bio

Search for `myself()` in `index.html`. The bio is a JavaScript array of
strings, joined with line breaks. Edit the strings to whatever you want:

```javascript
myself() {
  print([
    '<span class="dim">—# myself #——————————————————————</span>',
    '',
    'your new bio paragraph goes here.',
    '',
    '<span class="dim">interests ::</span>  whatever you want to list',
  ].join('<br>'));
},
```

You can use any HTML inside the strings - `<a>` for links, `<strong>`
for emphasis, `<ul>` for lists, etc. The styles in the page already
handle them.

## Adding a Blog Post

1. Create a new file `blog/<slug>.html` with your post body. Use plain
   HTML — `<p>`, `<a>`, `<code>`, `<pre>`, `<ul>`. No `<html>` or
   `<body>` wrapper needed.

2. Open `blog/posts.json` and add an entry to the array:

   ```json
   [
     { "slug": "my-new-post", "title": "what i think about X", "date": "2025-12-01" },
     { "slug": "older-post", "title": "an older one", "date": "2025-11-01" }
   ]
   ```

   Posts are automatically sorted newest-first by date.

3. If you have the `subscribe` command, run `subscribe build` in your
   shell to regenerate `feed.xml`, then upload the downloaded file to
   your site root.

Projects work identically — the directory is `projects/` instead of
`blog/`.

## Adding a New Theme

Search for `THEMES` in `index.html`. Add a line:

```javascript
const THEMES = {
  light:    { bg: '#E2E2E2', fg: '#000000' },
  dark:     { bg: '#000000', fg: '#E2E2E2' },
  oceanic:  { bg: '#1B2B34', fg: '#C0C5CE' },
};
```

Two colors, that's it - background and foreground. The shell uses CSS
variables so the new theme works everywhere automatically. Test with
`theme oceanic`.

## Adding a Custom Command

Search for `const commands = {` in `index.html`. Add a new entry:

```javascript
const commands = {
  // ...existing commands...
  hello() {
    printText('hello, friend.');
  },
  shout(args) {
    printText(args.join(' ').toUpperCase() + '!!!');
  },
  async joke() {
    try {
      const res = await fetch('https://icanhazdadjoke.com/', {
        headers: { 'Accept': 'application/json' }
      });
      const data = await res.json();
      printText(data.joke);
    } catch (e) {
      printText('joke: failed to fetch.');
    }
  },
};
```

Helper functions you can use inside command bodies:

- `print(html)` - print a line of HTML
- `printText(text)` - print a line of plain text (escaped)
- `blank()` - print a blank line
- `escapeHtml(s)` - make user input safe to embed in HTML

To make your command appear in `help`, search for the `rows` array in
the `help()` command and add a row:

```javascript
['hello',           'greet you back'],
['shout <text>',    'shout text in caps'],
```

To add a manual page, search for `pages` in `man()` and add an entry:

```javascript
hello: 'HELLO(1)\n\nNAME\n    hello — say hi\n\nSYNOPSIS\n    hello',
```

## Changing Your Prompt

Search for `const USER`, `const HOST`, `const DIR` near the top of the
script. Edit any of them. The prompt format is fixed at
`USER@HOST[DIR]:  $` but the values can be anything including unicode.

## Removing the "matcha" Attribution

You're allowed to (as much as I'd love the traffic to my site). Search for `generated with` near the bottom of the
script and delete that `print(...)` line.

## Custom Error Pages

Edit `404.html`, `403.html`, `410.html`, and `500.html`. They're
standalone files with their own inline styles - change the message,
add ASCII art, whatever fits.

The `.htaccess` file wires them to the right HTTP status codes on
Apache and LiteSpeed hosts. If you're on nginx, you'll need different
configuration in your server block:

```nginx
error_page 403 /403.html;
error_page 404 /404.html;
error_page 410 /410.html;
error_page 500 /500.html;
```

## More Extensive Changes

If you want to do anything bigger - different layouts, multiple pages,
build-tool integration - at that point you've outgrown the matcha
output and should treat it as a starting point you've forked. The
generated site is yours; do what you want with it (you're probably
more creative than me anyhow)!
