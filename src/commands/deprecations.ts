import { differenceInDays, parseISO } from 'date-fns';
import { loadDeprecations } from './shared.js';
import { title, c } from '../lib/format.js';

export function runDeprecations(harness?: string) {
  const items = loadDeprecations().filter((d) => !harness || d.harness === harness);
  title(`deprecations${harness ? ` (${harness})` : ''}`);
  for (const d of items) {
    const sev = d.severity === 'critical' ? c.err('critical') : d.severity === 'warning' ? c.warn('warning') : c.info('info');
    const days = d.removal ? differenceInDays(parseISO(d.removal), new Date()) : null;
    console.log(`- ${c.info(d.harness)}: ${d.item} [${sev}]`);
    console.log(`  deprecated ${d.deprecated}${d.removal ? `, removal ${d.removal} (${days}d)` : ''}`);
    console.log(`  -> ${d.migration}`);
  }
}
