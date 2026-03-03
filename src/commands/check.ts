import { readFileSync } from 'node:fs';
import { detectHarness } from '../lib/config-detect.js';
import { loadDeprecations, loadRecommendations } from './shared.js';
import { c, title } from '../lib/format.js';

export function runCheck(configPath: string) {
  const harness = detectHarness(configPath);
  title(`config check: ${configPath}`);
  if (!harness) {
    console.log(c.warn('Could not detect harness from filename.'));
    return;
  }
  console.log(c.info(`Detected harness: ${harness}`));
  const content = readFileSync(configPath, 'utf-8');

  const deprecations = loadDeprecations().filter((d) => d.harness === harness);
  for (const d of deprecations) {
    const needle = d.item.toLowerCase().split(' ')[0];
    if (content.toLowerCase().includes(needle)) {
      console.log(`${c.warn('⚠')} ${d.item} — ${d.migration}`);
    }
  }

  const recs = loadRecommendations().filter((r) => r.harness === harness);
  if (recs.length) {
    console.log(c.h('\nRecommended actions:'));
    recs.slice(0, 3).forEach((r) => console.log(`  - ${r.action}: ${r.how}`));
  }
}
