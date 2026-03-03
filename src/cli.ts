import { Command } from 'commander';
import { runDiff } from './commands/diff.js';
import { runCheck } from './commands/check.js';
import { runModels } from './commands/models.js';
import { runCompat } from './commands/compat.js';
import { runDeprecations } from './commands/deprecations.js';
import { runTui } from './commands/tui.js';

const program = new Command();
program.name('cliwatch').description('Track harness changelog updates and config impact').version('0.1.0');

program.command('diff').argument('[harness]').option('--since <duration>', 'e.g. 7d, 2w, 1m', '7d').action((h, o) => runDiff(h, o.since));
program.command('check').argument('<configPath>').action((p) => runCheck(p));
program.command('models').option('--new', 'last 90 days').option('--provider <name>').action((o) => runModels({ newer: o.new, provider: o.provider }));
program.command('compat').argument('<modelId>').action((m) => runCompat(m));
program.command('deprecations').option('--harness <name>').action((o) => runDeprecations(o.harness));
program.command('tui').action(() => runTui());

program.parse();
