import { differenceInDays, parseISO } from 'date-fns';
import { loadModelsProviders } from '../lib/data.js';
import { title, c } from '../lib/format.js';

export function runModels(opts: { newer?: boolean; provider?: string }) {
  const providers = loadModelsProviders();
  title('models');
  for (const p of providers) {
    if (opts.provider && p.provider.toLowerCase() !== opts.provider.toLowerCase()) continue;
    const models = (p.models || []).filter((m: any) => !opts.newer || (m.released && differenceInDays(new Date(), parseISO(m.released)) <= 90));
    if (!models.length) continue;
    console.log(c.info(`\n${p.provider}`));
    for (const m of models) {
      console.log(`  - ${m.id || m.name} ${c.dim(m.released || '')}`);
    }
  }
}
