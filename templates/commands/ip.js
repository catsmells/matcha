    async ip() {
      try {
        const res = await fetch('https://api.ipify.org?format=json');
        if (!res.ok) throw new Error('HTTP ' + res.status);
        const data = await res.json();
        printText(data.ip || '(unknown)');
      } catch (e) { printText('ip: failed to fetch.'); }
    },
