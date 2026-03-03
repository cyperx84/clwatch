import { basename } from 'node:path';

export function detectHarness(configPath: string): string | null {
  const f = basename(configPath).toLowerCase();
  if (f === 'claude.md' || f.includes('.claude') || f === 'settings.json') return 'claude-code';
  if (f === '.cursorrules') return 'cursor';
  if (f.includes('aider')) return 'aider';
  if (f.includes('cline')) return 'cline';
  if (f.includes('continue')) return 'continue';
  if (f.includes('codex')) return 'codex-cli';
  if (f.includes('openclaw') || f.includes('gateway')) return 'openclaw';
  if (f.includes('gemini')) return 'gemini-cli';
  return null;
}
