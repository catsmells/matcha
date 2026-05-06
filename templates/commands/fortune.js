    fortune() {
      const f = [
        '"the unexamined life is not worth living." — socrates',
        '"talk is cheap. show me the code." — linus torvalds',
        '"there is no spoon." — the matrix',
        '"we are stuck with technology when what we really want is just stuff that works." — douglas adams',
        'today is yesterday\'s tomorrow.',
        'edit this list in commands/fortune.js — make it your own.',
      ];
      printText(f[Math.floor(Math.random() * f.length)]);
    },
