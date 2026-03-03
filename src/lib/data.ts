import { readFileSync, readdirSync, existsSync } from 'node:fs';
import { join } from 'node:path';
import type { CompatData, Deprecation, ToolReleaseFile } from '../types.js';

const DATA_DIR = join(process.cwd(), 'data');

function safeRead<T>(path: string, fallback: T): T {
  try {
    return JSON.parse(readFileSync(path, 'utf-8')) as T;
  } catch {
    return fallback;
  }
}

export function loadAllReleases(): Record<string, ToolReleaseFile> {
  const dir = join(DATA_DIR, 'releases');
  if (!existsSync(dir)) return {};
  const out: Record<string, ToolReleaseFile> = {};
  for (const f of readdirSync(dir).filter((x) => x.endsWith('.json'))) {
    out[f.replace('.json', '')] = safeRead<ToolReleaseFile>(join(dir, f), { tool: { id: f.replace('.json', ''), name: f }, releases: [] });
  }
  return out;
}

export function loadCompat(): CompatData {
  return safeRead<CompatData>(join(DATA_DIR, 'compatibility.json'), { models: [], harnesses: {} });
}

export function loadDeprecations(): Deprecation[] {
  return safeRead<Deprecation[]>(join(DATA_DIR, 'deprecations.json'), []);
}

export function loadModelsProviders(): any[] {
  const dir = join(DATA_DIR, 'models');
  if (!existsSync(dir)) return [];
  return readdirSync(dir).filter((x) => x.endsWith('.json')).map((f) => safeRead<any>(join(dir, f), { provider: f, models: [] }));
}
