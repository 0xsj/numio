numio/
│
├── go.mod
├── go.sum
├── README.md
├── LICENSE
├── Makefile
│
├── docs/
│   ├── roadmap.md                          # This file
│   ├── architecture.md                     # System design
│   ├── grammar.md                          # Parser grammar spec
│   ├── api.md                              # Public API documentation
│   └── examples/
│       ├── basics.md
│       ├── finance.md
│       ├── dates.md
│       └── scripting.md
│
│
│── ════════════════════════════════════════════════════════════════
│   BINARIES
│── ════════════════════════════════════════════════════════════════
│
├── cmd/
│   │
│   ├── numio/                              # TUI binary
│   │   └── main.go
│   │
│   ├── numio-cli/                          # CLI binary
│   │   └── main.go
│   │
│   │── ─ ─ ─ ─ ─ ─ FUTURE ─ ─ ─ ─ ─ ─
│   │
│   ├── numio-server/                       # HTTP/JSON-RPC server
│   │   └── main.go
│   │
│   ├── numio-lsp/                          # Language Server Protocol
│   │   └── main.go
│   │
│   └── numio-wasm/                         # WASM build for browser
│       └── main.go
│
│
│── ════════════════════════════════════════════════════════════════
│   PUBLIC API
│── ════════════════════════════════════════════════════════════════
│
├── pkg/
│   └── numio/
│       ├── numio.go                        # Main public interface
│       ├── options.go                      # Functional options pattern
│       ├── result.go                       # Result types for API consumers
│       │
│       │── ─ ─ ─ ─ ─ ─ FUTURE ─ ─ ─ ─ ─ ─
│       │
│       ├── wasm.go                         # WASM-specific exports
│       └── embed.go                        # Embeddable engine helpers
│
│
│── ════════════════════════════════════════════════════════════════
│   INTERNAL PACKAGES
│── ════════════════════════════════════════════════════════════════
│
├── internal/
│   │
│   │── ══════════════════════════════════════════════════════════
│   │   FOUNDATION LAYER (zero internal deps)
│   │── ══════════════════════════════════════════════════════════
│   │
│   ├── errors/
│   │   └── errors.go                       # ✅ DONE - Custom error types
│   │
│   ├── token/
│   │   └── token.go                        # ✅ DONE - Lexer token types
│   │
│   ├── config/
│   │   ├── config.go                       # Configuration struct
│   │   ├── defaults.go                     # Default values
│   │   ├── file.go                         # File-based config (~/.numio/config.json)
│   │   │
│   │   │── ─ ─ ─ ─ ─ ─ FUTURE ─ ─ ─ ─ ─ ─
│   │   │
│   │   ├── apikeys.go                      # API key management
│   │   └── profiles.go                     # Named config profiles
│   │
│   │── ══════════════════════════════════════════════════════════
│   │   TYPE SYSTEM LAYER
│   │── ══════════════════════════════════════════════════════════
│   │
│   ├── types/
│   │   │
│   │   │── ─ ─ ─ ─ ─ ─ CORE ─ ─ ─ ─ ─ ─
│   │   │
│   │   ├── currency.go                     # Fiat currencies (USD, EUR, TRY, etc.)
│   │   │                                   #   - 30+ curated with symbols/aliases
│   │   │                                   #   - Dynamic support for any ISO 4217
│   │   │                                   #   - "dollars" → USD, "turkish lira" → TRY
│   │   │
│   │   ├── crypto.go                       # Cryptocurrencies
│   │   │                                   #   - BTC, ETH, SOL, etc.
│   │   │                                   #   - CoinGecko ID mapping
│   │   │
│   │   ├── metal.go                        # Precious metals
│   │   │                                   #   - XAU (gold), XAG (silver)
│   │   │                                   #   - XPT (platinum), XPD (palladium)
│   │   │
│   │   ├── unit.go                         # Physical units
│   │   │                                   #   - Length: km, m, mi, ft, in
│   │   │                                   #   - Weight: kg, g, lb, oz
│   │   │                                   #   - Time: h, min, s, days, weeks
│   │   │                                   #   - Temperature: C, F, K
│   │   │                                   #   - Data: TB, GB, MB, KB, B
│   │   │                                   #   - Area: sqm, sqft, acres
│   │   │                                   #   - Volume: L, mL, gal, cups
│   │   │
│   │   ├── value.go                        # Value sum type
│   │   │                                   #   - Number, Percentage, Currency
│   │   │                                   #   - WithUnit, Empty, Error
│   │   │
│   │   │── ─ ─ ─ ─ ─ ─ FUTURE ─ ─ ─ ─ ─ ─
│   │   │
│   │   ├── date.go                         # Date/DateTime type
│   │   │                                   #   - Parsing: 2024-01-15, "Jan 15, 2024"
│   │   │                                   #   - Relative: today, tomorrow, next monday
│   │   │                                   #   - Parts: year(), month(), weekday()
│   │   │
│   │   ├── duration.go                     # Duration type (distinct from time units)
│   │   │                                   #   - "3 years, 2 months, 5 days"
│   │   │                                   #   - Result of date subtraction
│   │   │
│   │   ├── boolean.go                      # Boolean type
│   │   │                                   #   - true, false
│   │   │                                   #   - Result of comparisons
│   │   │
│   │   ├── list.go                         # List/Range type
│   │   │                                   #   - [1, 2, 3, 4, 5]
│   │   │                                   #   - 1..10, 1..100 step 5
│   │   │
│   │   ├── table.go                        # Table/structured data
│   │   │                                   #   - Named columns
│   │   │                                   #   - Query support
│   │   │
│   │   ├── compound.go                     # Compound units
│   │   │                                   #   - km/h, $/hour, kg/m³
│   │   │                                   #   - Unit algebra
│   │   │
│   │   ├── stock.go                        # Stock/equity type
│   │   │                                   #   - AAPL, TSLA, etc.
│   │   │                                   #   - Price, volume, change
│   │   │
│   │   ├── commodity.go                    # Commodities
│   │   │                                   #   - Oil, natural gas
│   │   │
│   │   └── asset.go                        # Generic asset interface
│   │                                       #   - Unifies currency, crypto, metal, stock
│   │
│   │── ══════════════════════════════════════════════════════════
│   │   PARSING LAYER
│   │── ══════════════════════════════════════════════════════════
│   │
│   ├── lexer/
│   │   ├── lexer.go                        # Core tokenizer
│   │   │                                   #   - Numbers: 42, 3.14, 1,234.56
│   │   │                                   #   - Operators: +, -, *, /, ^
│   │   │                                   #   - Identifiers, keywords
│   │   │                                   #   - Currency symbols: $, €, £, ₺
│   │   │
│   │   │── ─ ─ ─ ─ ─ ─ FUTURE ─ ─ ─ ─ ─ ─
│   │   │
│   │   ├── lexer_number.go                 # Extended number formats
│   │   │                                   #   - Hex: 0xFF, 0x1A2B
│   │   │                                   #   - Binary: 0b1010
│   │   │                                   #   - Octal: 0o777
│   │   │                                   #   - Scientific: 1.5e6, 2.3E-4
│   │   │
│   │   ├── lexer_date.go                   # Date literal tokenization
│   │   │                                   #   - 2024-01-15
│   │   │                                   #   - Jan 15, 2024
│   │   │                                   #   - 15/01/2024
│   │   │
│   │   └── lexer_string.go                 # String literals
│   │                                       #   - "hello world"
│   │                                       #   - For imports, exports
│   │
│   ├── ast/
│   │   ├── ast.go                          # Core AST nodes
│   │   │                                   #   - Literals: Number, Percent, Currency
│   │   │                                   #   - BinaryOp, UnaryOp
│   │   │                                   #   - Assignment, Variable
│   │   │                                   #   - FunctionCall, Conversion
│   │   │
│   │   │── ─ ─ ─ ─ ─ ─ FUTURE ─ ─ ─ ─ ─ ─
│   │   │
│   │   ├── ast_conditional.go              # Conditional expressions
│   │   │                                   #   - if/then/else
│   │   │                                   #   - Ternary: x > 0 ? x : -x
│   │   │
│   │   ├── ast_range.go                    # Range expressions
│   │   │                                   #   - 1..10
│   │   │                                   #   - 1..100 step 5
│   │   │
│   │   ├── ast_comprehension.go            # List comprehensions
│   │   │                                   #   - [x * 2 for x in 1..10]
│   │   │                                   #   - [x for x in 1..100 where x % 2 == 0]
│   │   │
│   │   ├── ast_procedure.go                # User-defined procedures
│   │   │                                   #   - def withTax(amount): ...
│   │   │
│   │   └── ast_import.go                   # Import statements
│   │                                       #   - import "data.csv" as data
│   │
│   ├── parser/
│   │   ├── parser.go                       # Recursive descent parser
│   │   │                                   #   - Operator precedence
│   │   │                                   #   - Error recovery
│   │   │                                   #   - Fuzzy/natural language support
│   │   │
│   │   │── ─ ─ ─ ─ ─ ─ FUTURE ─ ─ ─ ─ ─ ─
│   │   │
│   │   ├── parser_date.go                  # Date expression parsing
│   │   │                                   #   - today + 30 days
│   │   │                                   #   - days between X and Y
│   │   │
│   │   ├── parser_conditional.go           # Conditional parsing
│   │   │
│   │   ├── parser_natural.go               # Enhanced natural language
│   │   │                                   #   - "60 dollars to turkish lira"
│   │   │                                   #   - "price of gold"
│   │   │                                   #   - "20 percent of 150"
│   │   │
│   │   └── parser_recovery.go              # Error recovery strategies
│   │
│   │── ══════════════════════════════════════════════════════════
│   │   EVALUATION LAYER
│   │── ══════════════════════════════════════════════════════════
│   │
│   ├── eval/
│   │   ├── context.go                      # Evaluation context
│   │   │                                   #   - Variables map
│   │   │                                   #   - Rate cache reference
│   │   │                                   #   - Special vars: _, ANS, total
│   │   │
│   │   ├── eval.go                         # Main evaluator
│   │   │                                   #   - AST walking
│   │   │                                   #   - Type dispatch
│   │   │
│   │   ├── operators.go                    # Binary/unary operators
│   │   │                                   #   - Arithmetic: +, -, *, /, ^
│   │   │                                   #   - With type coercion
│   │   │                                   #   - Unit-aware operations
│   │   │
│   │   ├── functions.go                    # Core built-in functions
│   │   │                                   #   - sum(), avg(), min(), max()
│   │   │                                   #   - abs(), sqrt(), round()
│   │   │                                   #   - floor(), ceil()
│   │   │
│   │   ├── conversion.go                   # Unit/currency conversion
│   │   │                                   #   - "in", "to" keywords
│   │   │                                   #   - Rate lookup via cache
│   │   │
│   │   │── ─ ─ ─ ─ ─ ─ FUTURE ─ ─ ─ ─ ─ ─
│   │   │
│   │   ├── functions_math.go               # Extended math functions
│   │   │                                   #   - log(), ln(), log10()
│   │   │                                   #   - sin(), cos(), tan()
│   │   │                                   #   - asin(), acos(), atan()
│   │   │                                   #   - pow(), exp()
│   │   │                                   #   - factorial(), nPr(), nCr()
│   │   │
│   │   ├── functions_stats.go              # Statistics functions
│   │   │                                   #   - median(), mode()
│   │   │                                   #   - stddev(), variance()
│   │   │                                   #   - percentile()
│   │   │                                   #   - correlation()
│   │   │
│   │   ├── functions_finance.go            # Finance functions
│   │   │                                   #   - loan(principal, rate, term)
│   │   │                                   #   - compound(principal, rate, time)
│   │   │                                   #   - npv(rate, cashflows...)
│   │   │                                   #   - irr(cashflows...)
│   │   │                                   #   - depreciation()
│   │   │
│   │   ├── functions_date.go               # Date functions
│   │   │                                   #   - days_between(d1, d2)
│   │   │                                   #   - workdays_between(d1, d2)
│   │   │                                   #   - age(birthdate)
│   │   │                                   #   - add_days(), add_months()
│   │   │                                   #   - year(), month(), day()
│   │   │                                   #   - weekday(), quarter(), week()
│   │   │
│   │   ├── functions_string.go             # String functions
│   │   │                                   #   - base64(), base64decode()
│   │   │                                   #   - md5(), sha256()
│   │   │                                   #   - uuid()
│   │   │                                   #   - urlencode(), urldecode()
│   │   │
│   │   ├── functions_list.go               # List functions
│   │   │                                   #   - map(), filter(), reduce()
│   │   │                                   #   - first(), last(), nth()
│   │   │                                   #   - reverse(), sort()
│   │   │                                   #   - unique(), flatten()
│   │   │
│   │   ├── functions_bitwise.go            # Bitwise operations
│   │   │                                   #   - band(), bor(), bxor(), bnot()
│   │   │                                   #   - lshift(), rshift()
│   │   │
│   │   ├── comparison.go                   # Comparison operators
│   │   │                                   #   - ==, !=, <, >, <=, >=
│   │   │                                   #   - and, or, not
│   │   │
│   │   ├── coercion.go                     # Type coercion rules
│   │   │                                   #   - Number + Currency → Currency
│   │   │                                   #   - Percentage of X → X type
│   │   │
│   │   └── conditional.go                  # Conditional evaluation
│   │                                       #   - if/then/else
│   │                                       #   - Ternary operator
│   │
│   │── ══════════════════════════════════════════════════════════
│   │   DATA LAYER
│   │── ══════════════════════════════════════════════════════════
│   │
│   ├── cache/
│   │   ├── cache.go                        # Rate cache interface
│   │   │                                   #   - Get/Set rates
│   │   │                                   #   - TTL management
│   │   │
│   │   ├── memory.go                       # In-memory cache
│   │   │                                   #   - Fast access
│   │   │                                   #   - LRU eviction
│   │   │
│   │   ├── file.go                         # File-based persistence
│   │   │                                   #   - ~/.numio/cache/rates.json
│   │   │                                   #   - Load on startup
│   │   │
│   │   ├── bfs.go                          # BFS path-finding
│   │   │                                   #   - Find conversion path
│   │   │                                   #   - USD → TRY via EUR
│   │   │
│   │   ├── fallback.go                     # Hardcoded fallback rates
│   │   │                                   #   - Offline mode
│   │   │                                   #   - Last resort
│   │   │
│   │   │── ─ ─ ─ ─ ─ ─ FUTURE ─ ─ ─ ─ ─ ─
│   │   │
│   │   └── sqlite.go                       # SQLite persistence
│   │                                       #   - Historical rates
│   │                                       #   - Query past data
│   │
│   ├── fetch/
│   │   ├── fetch.go                        # HTTP client wrapper
│   │   │                                   #   - Rate limiting
│   │   │                                   #   - Retry logic
│   │   │                                   #   - Timeout handling
│   │   │
│   │   ├── provider.go                     # Provider interface
│   │   │                                   #   - Pluggable data sources
│   │   │
│   │   ├── fiat.go                         # Fiat currency rates
│   │   │                                   #   - open.er-api.com (primary)
│   │   │                                   #   - frankfurter.app (fallback)
│   │   │
│   │   ├── crypto.go                       # Crypto prices
│   │   │                                   #   - CoinGecko (primary)
│   │   │                                   #   - CoinCap (fallback)
│   │   │
│   │   ├── metals.go                       # Precious metal prices
│   │   │                                   #   - Gold, silver, platinum
│   │   │
│   │   ├── registry.go                     # Provider registry
│   │   │                                   #   - Register/lookup providers
│   │   │                                   #   - Fallback chain
│   │   │
│   │   │── ─ ─ ─ ─ ─ ─ FUTURE ─ ─ ─ ─ ─ ─
│   │   │
│   │   ├── stocks.go                       # Stock prices
│   │   │                                   #   - Yahoo Finance
│   │   │                                   #   - Alpha Vantage
│   │   │
│   │   ├── commodities.go                  # Commodity prices
│   │   │                                   #   - Oil, natural gas
│   │   │
│   │   └── scheduler.go                    # Background refresh
│   │                                       #   - Periodic rate updates
│   │                                       #   - Configurable interval
│   │
│   │── ══════════════════════════════════════════════════════════
│   │   ENGINE LAYER
│   │── ══════════════════════════════════════════════════════════
│   │
│   ├── engine/
│   │   ├── engine.go                       # Core engine
│   │   │                                   #   - Eval(), EvalMultiple()
│   │   │                                   #   - Variable management
│   │   │                                   #   - Rate cache integration
│   │   │
│   │   ├── line.go                         # Line result tracking
│   │   │                                   #   - Input, Value, consumed flag
│   │   │                                   #   - Continuation logic
│   │   │
│   │   ├── continuation.go                 # Continuation handling
│   │   │                                   #   - "+ 10" continues from previous
│   │   │                                   #   - "in EUR" converts previous
│   │   │
│   │   ├── totals.go                       # Grouped totals
│   │   │                                   #   - Sum by currency
│   │   │                                   #   - Sum by unit type
│   │   │
│   │   │── ─ ─ ─ ─ ─ ─ FUTURE ─ ─ ─ ─ ─ ─
│   │   │
│   │   ├── session.go                      # Named sessions
│   │   │                                   #   - Multiple independent contexts
│   │   │                                   #   - Switch between workspaces
│   │   │
│   │   ├── history.go                      # Undo/redo
│   │   │                                   #   - State snapshots
│   │   │                                   #   - Variable history tracking
│   │   │
│   │   ├── macro.go                        # Macro support
│   │   │                                   #   - Record/playback
│   │   │                                   #   - Named macros
│   │   │                                   #   - Persistence
│   │   │
│   │   ├── procedure.go                    # User-defined procedures
│   │   │                                   #   - def withTax(x): x * 1.21
│   │   │                                   #   - Parameters, body
│   │   │
│   │   ├── import.go                       # Data import
│   │   │                                   #   - CSV files
│   │   │                                   #   - JSON files
│   │   │                                   #   - HTTP URLs
│   │   │
│   │   ├── export.go                       # Data export
│   │   │                                   #   - JSON, CSV, Markdown
│   │   │                                   #   - Session state dump
│   │   │
│   │   └── watch.go                        # Live data watching
│   │                                       #   - Auto-refresh variables
│   │                                       #   - "watch btc every 30s"
│   │
│   │── ══════════════════════════════════════════════════════════
│   │   VISUALIZATION LAYER
│   │── ══════════════════════════════════════════════════════════
│   │
│   ├── graph/
│   │   │── ─ ─ ─ ─ ─ ─ FUTURE ─ ─ ─ ─ ─ ─
│   │   │
│   │   ├── graph.go                        # Graph interface
│   │   │
│   │   ├── sparkline.go                    # Inline spark charts
│   │   │                                   #   - spark(1, 4, 2, 8) → ▁▃▂▇
│   │   │
│   │   ├── bar.go                          # Horizontal bar charts
│   │   │                                   #   - ASCII art bars
│   │   │
│   │   ├── histogram.go                    # Histograms
│   │   │                                   #   - Frequency distribution
│   │   │
│   │   ├── line.go                         # Line charts
│   │   │                                   #   - Time series
│   │   │                                   #   - Variable history
│   │   │
│   │   ├── scatter.go                      # Scatter plots
│   │   │
│   │   ├── pie.go                          # Pie charts (ASCII)
│   │   │
│   │   ├── table.go                        # Table formatting
│   │   │                                   #   - Aligned columns
│   │   │                                   #   - Borders
│   │   │
│   │   └── export/
│   │       ├── svg.go                      # SVG export
│   │       ├── png.go                      # PNG export (via SVG)
│   │       └── html.go                     # Interactive HTML
│   │
│   │── ══════════════════════════════════════════════════════════
│   │   INTERFACE LAYER
│   │── ══════════════════════════════════════════════════════════
│   │
│   ├── highlight/
│   │   ├── highlight.go                    # Syntax highlighting
│   │   │                                   #   - Token → Color mapping
│   │   │                                   #   - Numbers, operators, units
│   │   │
│   │   │── ─ ─ ─ ─ ─ ─ FUTURE ─ ─ ─ ─ ─ ─
│   │   │
│   │   └── themes.go                       # Color themes
│   │                                       #   - Dark, light, solarized
│   │                                       #   - Custom themes
│   │
│   ├── format/
│   │   ├── format.go                       # Output formatting
│   │   │                                   #   - Value → String
│   │   │
│   │   │── ─ ─ ─ ─ ─ ─ FUTURE ─ ─ ─ ─ ─ ─
│   │   │
│   │   ├── precision.go                    # Precision control
│   │   │                                   #   - Decimal places
│   │   │                                   #   - Significant figures
│   │   │
│   │   ├── locale.go                       # Locale-aware formatting
│   │   │                                   #   - 1,234.56 vs 1.234,56
│   │   │                                   #   - Currency symbols
│   │   │
│   │   ├── notation.go                     # Number notation
│   │   │                                   #   - Scientific: 1.5e6
│   │   │                                   #   - Engineering: 1.5M
│   │   │                                   #   - Fraction: 1/3
│   │   │
│   │   └── accounting.go                   # Accounting format
│   │                                       #   - ($500) for negative
│   │                                       #   - Aligned decimals
│   │
│   ├── repl/
│   │   ├── repl.go                         # REPL loop
│   │   │                                   #   - Read, eval, print
│   │   │                                   #   - Command handling
│   │   │
│   │   │── ─ ─ ─ ─ ─ ─ FUTURE ─ ─ ─ ─ ─ ─
│   │   │
│   │   ├── commands.go                     # REPL commands
│   │   │                                   #   - :help, :clear, :quit
│   │   │                                   #   - :save, :load
│   │   │                                   #   - :vars, :history
│   │   │
│   │   ├── completion.go                   # Tab completion
│   │   │                                   #   - Variables, functions
│   │   │                                   #   - Currency codes, units
│   │   │
│   │   └── history.go                      # Input history
│   │                                       #   - Arrow key navigation
│   │                                       #   - Persistent history file
│   │
│   ├── server/
│   │   │── ─ ─ ─ ─ ─ ─ FUTURE ─ ─ ─ ─ ─ ─
│   │   │
│   │   ├── server.go                       # HTTP server
│   │   │
│   │   ├── jsonrpc.go                      # JSON-RPC 2.0 handler
│   │   │                                   #   - eval, eval_lines
│   │   │                                   #   - clear, get_variables
│   │   │                                   #   - reload_rates
│   │   │
│   │   ├── rest.go                         # REST API
│   │   │                                   #   - POST /eval
│   │   │                                   #   - GET /variables
│   │   │                                   #   - GET /rates
│   │   │
│   │   ├── websocket.go                    # WebSocket support
│   │   │                                   #   - Live updates
│   │   │                                   #   - Streaming results
│   │   │
│   │   └── middleware.go                   # HTTP middleware
│   │                                       #   - CORS, logging, auth
│   │
│   ├── lsp/
│   │   │── ─ ─ ─ ─ ─ ─ FUTURE ─ ─ ─ ─ ─ ─
│   │   │
│   │   ├── lsp.go                          # LSP main handler
│   │   │
│   │   ├── completion.go                   # Autocomplete
│   │   │                                   #   - Variables in scope
│   │   │                                   #   - Functions, units, currencies
│   │   │
│   │   ├── hover.go                        # Hover information
│   │   │                                   #   - Variable values
│   │   │                                   #   - Currency rates
│   │   │                                   #   - Function signatures
│   │   │
│   │   ├── diagnostics.go                  # Error diagnostics
│   │   │                                   #   - Parse errors
│   │   │                                   #   - Undefined variables
│   │   │                                   #   - Type mismatches
│   │   │
│   │   ├── definition.go                   # Go to definition
│   │   │                                   #   - Variable assignments
│   │   │
│   │   ├── references.go                   # Find references
│   │   │                                   #   - Variable usages
│   │   │
│   │   └── formatting.go                   # Document formatting
│   │
│   ├── tui/
│   │   │── ─ ─ ─ ─ ─ ─ FUTURE ─ ─ ─ ─ ─ ─
│   │   │
│   │   ├── tui.go                          # TUI main loop
│   │   │
│   │   ├── editor.go                       # Vim-style editor
│   │   │                                   #   - Normal/Insert modes
│   │   │                                   #   - Motions: h, j, k, l
│   │   │                                   #   - Commands: dd, yy, p
│   │   │
│   │   ├── viewport.go                     # Scrolling viewport
│   │   │                                   #   - Page up/down
│   │   │                                   #   - Mouse wheel
│   │   │
│   │   ├── statusbar.go                    # Status bar
│   │   │                                   #   - Mode indicator
│   │   │                                   #   - File name
│   │   │                                   #   - Cursor position
│   │   │
│   │   ├── results.go                      # Results panel
│   │   │                                   #   - Right-aligned values
│   │   │                                   #   - Syntax coloring
│   │   │
│   │   ├── footer.go                       # Footer / totals
│   │   │                                   #   - Grouped totals
│   │   │
│   │   ├── popup.go                        # Popup dialogs
│   │   │                                   #   - Help, quit confirm
│   │   │
│   │   ├── keybindings.go                  # Key mapping
│   │   │
│   │   └── mouse.go                        # Mouse handling
│   │
│   └── test/
│       │── ─ ─ ─ ─ ─ ─ FUTURE ─ ─ ─ ─ ─ ─
│       │
│       ├── assert.go                       # Assertion engine
│       │                                   #   - assert 2 + 2 == 4
│       │
│       └── runner.go                       # Test runner
│                                           #   - Run test blocks
│                                           #   - Report results
│
│
│── ════════════════════════════════════════════════════════════════
│   TESTING
│── ════════════════════════════════════════════════════════════════
│
└── test/
    │
    ├── testdata/
    │   ├── basic.numio                     # Basic arithmetic tests
    │   ├── percentages.numio               # Percentage operations
    │   ├── variables.numio                 # Variable tests
    │   ├── currencies.numio                # Currency tests
    │   ├── units.numio                     # Unit conversion tests
    │   ├── continuation.numio              # Continuation tests
    │   ├── functions.numio                 # Function tests
    │   │
    │   │── ─ ─ ─ ─ ─ ─ FUTURE ─ ─ ─ ─ ─ ─
    │   │
    │   ├── dates.numio                     # Date arithmetic tests
    │   ├── finance.numio                   # Finance function tests
    │   ├── ranges.numio                    # Range/list tests
    │   ├── conditionals.numio              # Conditional tests
    │   └── edge_cases.numio                # Edge cases
    │
    ├── integration/
    │   │── ─ ─ ─ ─ ─ ─ FUTURE ─ ─ ─ ─ ─ ─
    │   │
    │   ├── cli_test.go                     # CLI integration tests
    │   ├── server_test.go                  # Server integration tests
    │   └── lsp_test.go                     # LSP integration tests
    │
    ├── benchmarks/
    │   │── ─ ─ ─ ─ ─ ─ FUTURE ─ ─ ─ ─ ─ ─
    │   │
    │   ├── lexer_bench_test.go
    │   ├── parser_bench_test.go
    │   └── eval_bench_test.go
    │
    └── fuzz/
        │── ─ ─ ─ ─ ─ ─ FUTURE ─ ─ ─ ─ ─ ─
        │
        ├── lexer_fuzz_test.go
        └── parser_fuzz_test.go