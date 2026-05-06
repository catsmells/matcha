    cowsay(args) {
      const msg = args.join(' ') || 'moo.';
      const top = ' ' + '_'.repeat(msg.length + 2);
      const bot = ' ' + '-'.repeat(msg.length + 2);
      const out = [
        top, `< ${msg} >`, bot,
        '        \\   ^__^',
        '         \\  (oo)\\_______',
        '            (__)\\       )\\/\\',
        '                ||----w |',
        '                ||     ||',
      ].join('\n');
      print(`<span class="ascii">${escapeHtml(out)}</span>`);
    },
