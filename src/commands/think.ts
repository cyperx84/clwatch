import { execSync } from 'node:child_process';
import { parseISO, subDays } from 'date-fns';
import { loadAllReleases } from '../lib/data.js';
import { c, title } from '../lib/format.js';

interface LatticeModel {
  model_name: string;
  model_slug: string;
  category: string;
}

interface LatticeThinkResult {
  problem: string;
  models: LatticeModel[];
  summary?: string;
}

function latticeAvailable(): boolean {
  try {
    execSync('which lattice', { stdio: 'ignore' });
    return true;
  } catch {
    return false;
  }
}

function buildDiffSummary(harness: string): string {
  const days = 14;
  const cutoff = subDays(new Date(), days);
  const all = loadAllReleases();
  const entries = Object.entries(all).filter(([id]) => !harness || id === harness);

  const parts: string[] = [];
  for (const [id, file] of entries) {
    const recent = file.releases.filter((r) => parseISO(r.date) >= cutoff).slice(0, 3);
    if (!recent.length) continue;
    for (const r of recent) {
      const first = r.body?.split('\n').find((l) => l.trim().length > 0)?.trim() ?? r.name;
      parts.push(`${id} ${r.version}: ${first.slice(0, 120)}`);
    }
  }
  return parts.join('\n') || `Recent changes for ${harness || 'all harnesses'}`;
}

export function runThink(harness: string) {
  if (!latticeAvailable()) {
    console.log(c.warn('Install lattice for mental model analysis: brew install cyperx84/tap/lattice'));
    return;
  }

  const summary = buildDiffSummary(harness);

  title(`mental models for ${harness || 'recent changes'}`);

  try {
    const output = execSync(`lattice think "${summary.replace(/"/g, '\\"')}" --json`, {
      encoding: 'utf-8',
      timeout: 30000,
    });

    const result: LatticeThinkResult = JSON.parse(output);

    if (result.models?.length) {
      console.log('');
      for (const m of result.models) {
        console.log(`  ${c.info(m.model_name)} ${c.dim(`(${m.category})`)}`);
      }
    }

    if (result.summary) {
      console.log(`\n  ${c.dim('Synthesis:')}`);
      for (const line of result.summary.split('\n')) {
        console.log(`  ${line}`);
      }
    }
  } catch (e) {
    console.log(c.err(`Lattice failed: ${e instanceof Error ? e.message : String(e)}`));
  }
}

export function getLatticeSection(summary: string): string | null {
  if (!latticeAvailable()) return null;

  try {
    const output = execSync(`lattice think "${summary.replace(/"/g, '\\"')}" --json`, {
      encoding: 'utf-8',
      timeout: 30000,
    });

    const result: LatticeThinkResult = JSON.parse(output);
    if (!result.models?.length) return null;

    const lines: string[] = ['', '## 🧠 Mental Models', ''];
    for (const m of result.models) {
      lines.push(`  ${m.model_name} (${m.category})`);
    }
    if (result.summary) {
      lines.push('');
      lines.push(`  ${result.summary.split('\n')[0]}`);
    }
    return lines.join('\n');
  } catch {
    return null;
  }
}
