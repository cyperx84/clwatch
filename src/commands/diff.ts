import { parseISO, subDays } from 'date-fns';
import { loadAllReleases } from '../lib/data.js';
import { c, title } from '../lib/format.js';

function sinceToDays(s: string) {
  const m = s.match(/^(\d+)([dwm])$/i);
  if (!m) return 7;
  const n = Number(m[1]);
  return m[2].toLowerCase() === 'w' ? n * 7 : m[2].toLowerCase() === 'm' ? n * 30 : n;
}

export function runDiff(harness?: string, since = '7d') {
  const days = sinceToDays(since);
  const cutoff = subDays(new Date(), days);
  const all = loadAllReleases();
  const entries = Object.entries(all).filter(([id]) => !harness || id === harness);

  title(`changes in last ${since}${harness ? ` for ${harness}` : ''}`);
  for (const [id, file] of entries) {
    const recent = file.releases.filter((r) => parseISO(r.date) >= cutoff).slice(0, 3);
    if (!recent.length) continue;
    console.log(c.info(`${id}`));
    for (const r of recent) {
      console.log(`  ${c.ok('+')} ${r.version} ${c.dim(`(${r.date})`)}`);
      const first = r.body?.split('\n').find((l) => l.trim().length > 0)?.trim() ?? r.name;
      console.log(`    ${first.slice(0, 120)}`);
    }
  }
}
