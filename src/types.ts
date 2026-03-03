export interface Release {
  version: string;
  name: string;
  body: string;
  date: string;
  url: string;
  prerelease?: boolean;
}

export interface ToolReleaseFile {
  tool: { id: string; name: string; category?: string };
  releases: Release[];
  lastUpdated?: string;
}

export interface CompatData {
  models: { id: string; provider: string; name: string }[];
  harnesses: Record<string, { supported: string[]; default?: string; notes?: string }>;
}

export interface Deprecation {
  harness: string;
  item: string;
  deprecated: string;
  removal?: string | null;
  severity: 'critical' | 'warning' | 'info';
  migration: string;
  affects: string;
}
