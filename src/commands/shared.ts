import { readFileSync } from 'node:fs';
import { join } from 'node:path';
import type { Deprecation } from '../types.js';

export interface Recommendation {
  harness: string;
  version: string;
  action: string;
  impact: 'high'|'medium'|'low';
  description: string;
  how: string;
  tags?: string[];
}

export function loadRecommendations(): Recommendation[] {
  try { return JSON.parse(readFileSync(join(process.cwd(), 'data/recommendations.json'),'utf-8')); }
  catch { return []; }
}

export function loadDeprecations(): Deprecation[] {
  try { return JSON.parse(readFileSync(join(process.cwd(), 'data/deprecations.json'),'utf-8')); }
  catch { return []; }
}
