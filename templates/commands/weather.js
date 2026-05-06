    async weather(args) {
      const city = args.join(' ').trim();
      if (!city) { printText('usage: weather <city>'); return; }
      printText(`fetching weather for ${city} ...`);
      try {
        const geoRes = await fetch(`https://geocoding-api.open-meteo.com/v1/search?name=${encodeURIComponent(city)}&count=1`);
        const geo = await geoRes.json();
        if (!geo.results || !geo.results.length) { printText(`weather: no place named "${city}".`); return; }
        const { latitude, longitude, name, country, admin1 } = geo.results[0];
        const wRes = await fetch(`https://api.open-meteo.com/v1/forecast?latitude=${latitude}&longitude=${longitude}&current=temperature_2m,relative_humidity_2m,weather_code,wind_speed_10m&temperature_unit=celsius&wind_speed_unit=kmh`);
        const w = await wRes.json();
        const c = w.current || {};
        const WMO = { 0:'clear',1:'mainly clear',2:'partly cloudy',3:'overcast',45:'fog',48:'rime fog',51:'light drizzle',53:'drizzle',55:'dense drizzle',61:'light rain',63:'rain',65:'heavy rain',71:'light snow',73:'snow',75:'heavy snow',80:'showers',81:'showers',82:'violent showers',95:'thunderstorm',96:'thunderstorm',99:'severe thunderstorm' };
        const tC = c.temperature_2m;
        const tF = Math.round((tC * 9 / 5 + 32) * 10) / 10;
        const place = [name, admin1, country].filter(Boolean).join(', ');
        const cond = WMO[c.weather_code] || `code ${c.weather_code}`;
        print([
          `<span class="dim">${escapeHtml(place)}</span>`,
          `${tC}°C / ${tF}°F · ${escapeHtml(cond)} · humidity ${c.relative_humidity_2m}% · wind ${c.wind_speed_10m} km/h`,
        ].join('<br>'));
      } catch (e) { printText('weather: failed to fetch.'); }
    },
