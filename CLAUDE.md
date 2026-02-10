# numio

> A natural language calculator written in Go — a more powerful alternative to apps like Numi.

## Project Overview

numio is an intelligent, natural-language-aware calculator that supports currencies, units, cryptocurrencies, precious metals, financial functions, and more. It exposes multiple interfaces: CLI, TUI (terminal), and desktop (WIP). The core philosophy is to let users write math the way they think — `rent = 45000TL/month`, `$100 to lira`, `15% of 200`, `1 btc in usd` — while maintaining a precise, extensible evaluation engine underneath.

## Architecture
```
numio/
├── cmd/
│   ├── numio-cli/              # CLI REPL binary
│   ├── numio-tui/              # TUI binary (Bubble Tea)
│   └── numio-desktop/          # Desktop binary (Fyne, WIP)
├── pkg/                        # Public API surface
│   ├── engine/                 # Engine: orchestrates parse → eval, rate refresh
│   ├── types/                  # Value sum type, Currency, Crypto, Metal, Unit, Period
│   ├── cache/                  # Rate cache with BFS pathfinding for indirect conversions
│   └── errors/                 # Shared error types
├── internal/                   # Private implementation
│   ├── token/                  # Token definitions (operators, keywords, literals)
│   ├── lexer/                  # Tokenizer (multi-word identifiers, currency symbols)
│   ├── ast/                    # AST node definitions
│   ├── parser/                 # Recursive descent parser with operator precedence
│   ├── eval/                   # Evaluator, context (variables, line history, continuations)
│   │   ├── operators.go        # Arithmetic, percentage, rate arithmetic
│   │   ├── conversions.go      # Currency/unit/crypto/metal conversion
│   │   ├── continuations.go    # Line continuation logic (previous result carries forward)
│   │   ├── functions.go        # Built-in function registry + dispatch
│   │   └── context.go          # Variable store, line results, rate cache adapter
│   ├── nlp/                    # Natural language preprocessing
│   ├── explain/                # Step-by-step trace builder + formatter
│   ├── fetch/                  # Live rate fetching (multi-provider with fallback)
│   │   ├── registry.go         # Provider orchestration
│   │   ├── fiat.go             # Frankfurter, ExchangeRate-API
│   │   ├── crypto.go           # CoinGecko, CoinCap
│   │   └── metals.go           # Metals.live, GoldAPI, MetalpriceAPI
│   ├── highlight/              # Syntax highlighting (themes: Dracula, Monokai, Gruvbox, etc.)
│   ├── fuzzy/                  # Fuzzy matching for autocorrect
│   ├── tui/                    # TUI app (Bubble Tea + Lip Gloss)
│   │   └── keymap/             # Vim-style keybindings, actions, motions
│   ├── rpc/                    # JSON-RPC 2.0 protocol (for desktop UI communication)
│   └── editor/                 # Headless editor model (buffer, cursor, selection, history)
└── docs/
    └── roadmap.md
```

## Key Design Patterns

- **Registry / Provider** — Fetch layer uses provider interface with priority-ordered fallback chains
- **BFS Pathfinding** — Cache resolves indirect conversions (e.g. USD→EUR→TRY) via graph traversal
- **Sum Type (Value)** — Single `Value` type carries Number, Currency, Unit, Crypto, Metal, Percentage, Rate, Date, String, Error
- **Adapter** — `RateCacheAdapter` decouples evaluator from concrete cache implementation
- **Strategy** — Highlight themes, fetch providers, and key bindings are all swappable
- **Separation of Concerns** — Eval is split into operators, conversions, continuations, functions, and context
- **Leaf-upward development** — Always implement from the lowest-dependency module first

## What's Implemented

### Core Language
- Arithmetic: `+`, `-`, `*`, `/`, `^`, `!`, `%`, parentheses
- `x` / `X` as natural multiplication (`175 x 2`)
- Variables: `x = 42`, `rent = 45000TL/month`
- Percentages: `15% of 200`, `200 + 10%`, `200 - 25%`
- Line continuations: previous result carries forward implicitly
- Multi-target comparison: `$100 in EUR, GBP, JPY`

### Type System
- Currencies: 170+ fiat currencies with aliases (`lira` → TRY, `dollars` → USD)
- Crypto: BTC, ETH, SOL, etc. via CoinGecko/CoinCap
- Precious metals: XAU, XAG, XPT, XPD with live spot pricing
- Units: length, weight, time, temperature, data, area, volume, speed, frequency
- Rates & periods: `45000TL/month`, `rent * year`, period conversion (month↔year↔week↔day↔hour)

### Functions
- Math: `sqrt`, `abs`, `ceil`, `floor`, `round`, `log`, `ln`, `sin`, `cos`, `tan`, `pow`, `factorial`
- Statistics: `sum`, `avg`, `min`, `max`, `median`, `stddev`, `variance`
- Finance: `loan`, `compound`, `npv`, `irr`, `pmt`, `fv`, `pv`
- Visualization: `spark`, `sparktrend`, `sparkstats`, `bar`, `progress`
- Calculus: `derivative`, `integral`, `limit`, `root`, `taylor`

### Interfaces
- **CLI** — Simple REPL
- **TUI** — Full vim-modal editor with syntax highlighting, undo/redo, yank buffer, help overlay, explain mode (`Ctrl+E`), rate refresh (`Ctrl+R`), multiple color themes
- **Desktop (Fyne)** — WIP, basic evaluation works, vim keybindings partially ported

### Infrastructure
- Live rate fetching with multi-provider fallback (free APIs, optional API keys)
- Rate caching with BFS for indirect conversion paths
- Explain mode: step-by-step calculation trace
- Syntax highlighting with pluggable themes
- Hot reload dev workflow (Air)
- GitHub Actions CI/CD + Dependabot

## Known Issues & Bugs

- `1 gram gold in lira` parses `1 gram` as weight, `gold` as separate identifier — never connects them
- `gold price in usd` doesn't resolve metal spot price inline
- Turkish gold denominations (çeyrek, yarım, tam) not supported
- Some time/date expressions return `0` (`istanbul time to kst`, `utc time now`)
- `day` parsed as frequency (`1/day`) instead of date keyword
- Desktop: Fyne lacks native macOS window transparency/vibrancy
- Desktop: `Ctrl+E` (explain) doesn't trigger reliably in Fyne
- Desktop: cursor sizing mismatch with text size

## Planned Work

### High Priority
- **Session persistence** — save/restore `.numio` notebook files (variables, functions, history)
- **User-defined functions** — `def withTax(amount): amount * 1.08` with structural typing
- **What-if analysis** — `what if price = $80..$120 step $10` sensitivity tables
- **Fix gold/metal inline expressions** — connect `1 gram gold` parsing, add çeyrek/yarım/tam Turkish gold denominations

### Medium Priority
- **Conditional expressions** — `if/then/else`, ternary, comparison operators
- **List operations & ranges** — `[1,2,3]`, `1..10`, `sum(1..100)`
- **Stock prices** — new fetch provider for equities
- **Date/time improvements** — fix timezone conversions, add `days until christmas`, `today + 30 days`

### Desktop UI
- Resolve Fyne modifier key detection for `Ctrl+` shortcuts
- Fix cursor rendering proportional to font size
- Richer syntax highlighting in desktop (currency/unit colors)
- Explore native macOS vibrancy (may require CGo or platform-specific code)
- Alternative: revisit Tauri or Swift/AppKit if Fyne limitations block progress

### Long-term Roadmap
- **HTTP / JSON-RPC server** — expose engine as a service
- **LSP support** — editor integration (VS Code, Neovim)
- **WASM build** — browser-based numio
- **Boolean type** — true/false with logical operators
- **Tables & compound types** — structured data
- **Import/export** — CSV, JSON data ingestion
- **Plugin system** — user-extensible providers and functions

## Development Workflow

1. Always generate a module tree / file plan before starting work
2. Implement leaf modules first (types, errors, constants) → work upward
3. One file at a time; confirm before proceeding to the next
4. Offer tests for any testable unit of work
5. Complete file implementations preferred over partial patches
6. Maintain backward compatibility when adding features

## Build & Run
```bash
make build          # Build all binaries
make tui            # Run TUI with hot reload (Air)
make cli            # Run CLI with hot reload
make desktop        # Run Fyne desktop app
make test           # Run all tests
```
