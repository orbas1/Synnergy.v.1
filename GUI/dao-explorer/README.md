# DAO Explorer GUI

This reference interface demonstrates basic DAO management through the Synnergy CLI.

## Usage

```bash
# Create a DAO
node -e "import('./src/main.ts').then(m => m.createDAO('mydao','alice').then(id => console.log(id)))"

# Join a DAO
node -e "import('./src/main.ts').then(m => m.joinDAO('daoID','bob'))"

# List DAOs
node -e "import('./src/main.ts').then(m => m.listDAOs().then(console.log))"
```

The module spawns the local `synnergy` binary and therefore inherits all CLI configuration
such as network settings, gas table values, and authentication.
