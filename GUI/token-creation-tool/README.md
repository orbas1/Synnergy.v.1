# Token Creation Tool GUI

This lightweight interface demonstrates how new token contracts can be
created through the Synnergy CLI. It spawns the `synnergy` binary with
the desired token module and forwards standard flags.

## Usage

```bash
# Create a SYN500 token
node -e "import('./src/main.ts').then(m => m.createToken('syn500', {name:'Gold',symbol:'GLD',owner:'alice',decimals:2,supply:1000}))"

# List registered SYN500 tokens
node -e "import('./src/main.ts').then(m => m.listTokens('syn500')).then(console.log)"
```

The script assumes the `synnergy` CLI is on the PATH and configured with
the appropriate network and gas settings.
