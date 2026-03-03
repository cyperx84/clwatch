import chalk from 'chalk';

export const c = {
  h: (s: string) => chalk.green(s),
  dim: (s: string) => chalk.gray(s),
  ok: (s: string) => chalk.green(s),
  warn: (s: string) => chalk.yellow(s),
  err: (s: string) => chalk.red(s),
  info: (s: string) => chalk.cyan(s),
};

export function title(t: string) {
  console.log(c.h(`\n> ${t}`));
}
