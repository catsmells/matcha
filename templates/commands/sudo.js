    sudo(args) {
      if (!args.length) { printText('usage: sudo <command>'); return; }
      printText(`sorry, ${USER} is not in the sudoers file. this incident will be reported.`);
    },
