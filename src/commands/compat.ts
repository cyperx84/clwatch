import { loadCompat } from '../lib/data.js';
import { title, c } from '../lib/format.js';

export function runCompat(modelId: string) {
  const data = loadCompat();
  title(`compatibility for ${modelId}`);
  const matches = Object.entries(data.harnesses).filter(([, v]) => v.supported.includes(modelId));
  if (!matches.length) {
    console.log(c.warn('No harness support found.'));
    return;
  }
  for (const [h, v] of matches) {
    console.log(`- ${c.info(h)} ${v.default === modelId ? c.ok('(default)') : ''} ${c.dim(v.notes || '')}`);
  }
}
