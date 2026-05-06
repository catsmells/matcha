    echo(args) { printText(args.join(' ')); },
    clear() { $term.empty(); $currentInput = null; },
    whoami() { printText(USER); },
    pwd()    { printText('/home/' + USER + '/' + DIR); },
    date()   { printText(new Date().toString()); },
    uptime() {
      const s = Math.floor((Date.now() - startTime) / 1000);
      const h = Math.floor(s / 3600), m = Math.floor((s % 3600) / 60), sec = s % 60;
      printText(`up ${h}h ${m}m ${sec}s`);
    },
    history() {
      print(history.map((h, i) =>
        `<span class="dim">${String(i + 1).padStart(4)}</span>  ${escapeHtml(h)}`
      ).join('<br>') || '<span class="dim">(no history yet)</span>');
    },
    site(args) {
      if (!args[0]) { printText('usage: site <url>[.tld]'); return; }
      let target = args[0];
      if (!/\./.test(target)) target += '.com';
      if (!/^https?:\/\//i.test(target)) target = 'https://' + target;
      printText(`opening ${target} ...`);
      window.open(target, '_blank', 'noopener');
    },
    theme(args) {
      const names = Object.keys(THEMES);
      if (!args.length) {
        const idx = names.indexOf(currentTheme);
        applyTheme(names[(idx + 1) % names.length]);
        printText(`theme: ${currentTheme}`);
        return;
      }
      if (args[0] === 'list') {
        print(names.map(n =>
          n === currentTheme
            ? `<a data-cmd="theme ${n}">${escapeHtml(n)}</a> <span class="dim">(current)</span>`
            : `<a data-cmd="theme ${n}">${escapeHtml(n)}</a>`
        ).join('   '));
        return;
      }
      if (!THEMES[args[0]]) { printText(`theme: unknown theme "${args[0]}".`); return; }
      applyTheme(args[0]);
      printText(`theme: ${currentTheme}`);
    },
    exit() {
      printText('closing...');
      setTimeout(() => {
        window.close();
        setTimeout(() => printText('(your browser refused to close the tab.)'), 200);
      }, 120);
    },
