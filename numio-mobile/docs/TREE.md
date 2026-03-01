# Project Structure

## Directory Layout

```
numio-mobile/
├── app/                                # Expo Router (file-based routing)
│   ├── _layout.tsx                     # Root layout (providers, theme, fonts)
│   ├── modal.tsx                       # Generic modal route
│   ├── (tabs)/
│   │   ├── _layout.tsx                 # Tab navigator (Quick, Notepad, Settings)
│   │   ├── index.tsx                   # Quick mode (single-line calculator)
│   │   ├── notepad.tsx                 # Notepad mode (multi-line editor)
│   │   └── settings.tsx                # Preferences screen
│   └── explain/
│       └── [line].tsx                  # Explain modal (step-by-step trace)
│
├── lib/
│   └── engine/                         # Pure TS engine (zero RN deps)
│       ├── index.ts                    # Public API barrel
│       ├── engine.ts                   # Engine class (eval, explain, variables, rates)
│       ├── errors/
│       │   └── errors.ts               # NumioError, ParseError, EvalError
│       ├── types/
│       │   ├── value.ts                # Value discriminated union + constructors + formatting
│       │   ├── currency.ts             # Currency interface + registry (40+ fiat)
│       │   ├── unit.ts                 # Unit interface + registry (100+ units, 10 categories)
│       │   ├── crypto.ts               # Crypto interface + registry (30+ coins)
│       │   ├── metal.ts                # Metal interface + registry (10 metals)
│       │   └── period.ts               # Period enum + conversion math
│       ├── token/
│       │   └── token.ts                # TokenType enum, Token interface, keyword/symbol maps
│       ├── lexer/
│       │   └── lexer.ts                # Tokenizer (UTF-8 symbols, multi-word identifiers)
│       ├── ast/
│       │   └── ast.ts                  # AST node types (discriminated unions), Walk, helpers
│       ├── parser/
│       │   └── parser.ts               # Recursive descent parser, precedence climbing
│       ├── eval/
│       │   ├── eval.ts                 # Evaluator (AST dispatch, tracing hooks)
│       │   ├── context.ts              # Context (variables, previous, lines, rate adapter)
│       │   ├── operators.ts            # Binary/unary arithmetic, percentage math, rate ops
│       │   ├── conversions.ts          # Currency/unit/crypto/metal conversion dispatch
│       │   ├── continuations.ts        # Line continuation logic (previous result carry)
│       │   └── functions/
│       │       ├── registry.ts         # Function lookup + dispatch
│       │       ├── math.ts             # sqrt, abs, ceil, floor, round, log, ln, trig, pow
│       │       ├── stats.ts            # sum, avg, min, max, median, stddev, variance
│       │       └── finance.ts          # loan, compound, npv, irr, pmt, fv, pv
│       ├── nlp/
│       │   └── processor.ts            # Pattern-based NLP preprocessing
│       ├── cache/
│       │   └── rate-cache.ts           # RateCache with BFS pathfinding
│       ├── fetch/
│       │   ├── registry.ts             # Provider orchestration + fallback chains
│       │   ├── fiat.ts                 # Frankfurter, ExchangeRate-API
│       │   ├── crypto.ts               # CoinGecko, CoinCap
│       │   └── metals.ts               # Metals.live, GoldAPI
│       ├── explain/
│       │   ├── trace.ts                # Trace, Step data structures
│       │   ├── builder.ts              # Trace builder (records eval steps)
│       │   └── formatter.ts            # Format traces to human-readable strings
│       ├── highlight/
│       │   └── highlighter.ts          # Token -> Span[] (UI-agnostic, no ANSI)
│       └── fuzzy/
│           └── fuzzy.ts                # Levenshtein/Jaro-Winkler for autocorrect
│
├── components/
│   ├── quick/
│   │   ├── quick-input.tsx             # Single-line TextInput with live preview
│   │   ├── quick-result.tsx            # Large result display
│   │   └── quick-history.tsx           # Scrollable committed results
│   ├── notepad/
│   │   ├── notepad-editor.tsx          # Multi-line editor container
│   │   ├── notepad-line.tsx            # Single line: input + right-aligned result
│   │   ├── line-result.tsx             # Inline result display (tap for explain)
│   │   └── variables-panel.tsx         # Bottom sheet showing defined variables
│   ├── shared/
│   │   ├── styled-line.tsx             # Syntax-highlighted line (Text overlay)
│   │   ├── styled-span.tsx             # Single colored text span
│   │   ├── explain-view.tsx            # Step-by-step explanation display
│   │   └── rate-status-badge.tsx       # Rate freshness indicator
│   ├── keyboard/
│   │   └── keyboard-accessory.tsx      # Operator row above keyboard (%, $, +, etc.)
│   └── ui/                             # Base UI primitives (from Expo template)
│       ├── collapsible.tsx
│       ├── icon-symbol.tsx
│       └── icon-symbol.ios.tsx
│
├── hooks/
│   ├── use-engine.ts                   # Engine singleton + initialization
│   ├── use-eval.ts                     # Debounced eval/preview
│   ├── use-rates.ts                    # Rate refresh lifecycle (foreground, TTL)
│   ├── use-session.ts                  # Auto-save/restore notepad state
│   ├── use-theme.ts                    # Theme switching + system detection
│   ├── use-haptics.ts                  # Haptic feedback (eval commit, errors)
│   ├── use-color-scheme.ts             # System color scheme (Expo default)
│   └── use-theme-color.ts              # Themed color lookup (Expo default)
│
├── store/
│   ├── engine-store.ts                 # Zustand: engine instance, quick/notepad state
│   └── settings-store.ts              # Zustand: theme, precision, NLP toggle
│
├── theme/
│   ├── colors.ts                       # Color palettes per theme
│   ├── themes.ts                       # Dracula, Monokai, Gruvbox, Light definitions
│   └── tokens.ts                       # Token class -> color mapping per theme
│
├── services/
│   ├── storage.ts                      # MMKV wrapper (settings, rates, sessions)
│   └── session.ts                      # Session serialize/deserialize
│
├── constants/
│   └── theme.ts                        # Expo template theme constants
│
├── assets/
│   ├── fonts/                          # Custom fonts (if any)
│   └── images/                         # App icon, splash, tab icons
│
├── docs/
│   ├── ARCHITECTURE.md                 # System design + diagrams
│   └── TREE.md                         # This file
│
├── __tests__/
│   ├── engine/
│   │   ├── lexer.test.ts               # Tokenization tests
│   │   ├── parser.test.ts              # AST generation tests
│   │   ├── eval.test.ts                # Evaluation tests
│   │   ├── types.test.ts               # Registry lookup + Value formatting
│   │   └── integration.test.ts         # End-to-end engine tests
│   └── fixtures/
│       └── eval-cases.json             # Shared test vectors (Go + TS parity)
│
├── CLAUDE.md                           # Development guide
├── README.md                           # Project overview
├── app.json                            # Expo configuration
├── package.json                        # Dependencies + scripts
└── tsconfig.json                       # TypeScript configuration
```

## Build Phases

Implementation order follows leaf-first: types with zero deps first, then layers that depend on them.

### Phase 1: Engine Foundation

The core pipeline — types through evaluation. No network, no UI.

```
Step  File                              Depends On
─────────────────────────────────────────────────────────────
1.1   lib/engine/errors/errors.ts        (none)
1.2   lib/engine/types/period.ts         (none)
1.3   lib/engine/types/currency.ts       (none)
1.4   lib/engine/types/unit.ts           (none)
1.5   lib/engine/types/crypto.ts         (none)
1.6   lib/engine/types/metal.ts          (none)
1.7   lib/engine/types/value.ts          period, currency, unit, crypto, metal
1.8   lib/engine/token/token.ts          (none)
1.9   lib/engine/lexer/lexer.ts          token, types
1.10  lib/engine/ast/ast.ts              types
1.11  lib/engine/parser/parser.ts        lexer, ast, token, types
1.12  lib/engine/eval/context.ts         types
1.13  lib/engine/eval/operators.ts       types, ast
1.14  lib/engine/eval/conversions.ts     types, context
1.15  lib/engine/eval/continuations.ts   types, context
1.16  lib/engine/eval/functions/math.ts  types
1.17  lib/engine/eval/functions/stats.ts types
1.18  lib/engine/eval/functions/finance.ts types
1.19  lib/engine/eval/functions/registry.ts math, stats, finance
1.20  lib/engine/eval/eval.ts            ast, context, operators, conversions,
                                          continuations, functions
1.21  lib/engine/engine.ts               eval, parser, types
1.22  lib/engine/index.ts                engine (barrel export)
```

### Phase 2: NLP + Rates + Explain + Highlighting

Supporting subsystems that enhance the core.

```
Step  File                              Depends On
─────────────────────────────────────────────────────────────
2.1   lib/engine/nlp/processor.ts        types
2.2   lib/engine/cache/rate-cache.ts     types
2.3   lib/engine/fetch/registry.ts       cache
2.4   lib/engine/fetch/fiat.ts           registry
2.5   lib/engine/fetch/crypto.ts         registry
2.6   lib/engine/fetch/metals.ts         registry
2.7   lib/engine/explain/trace.ts        types, ast
2.8   lib/engine/explain/builder.ts      trace
2.9   lib/engine/explain/formatter.ts    trace
2.10  lib/engine/highlight/highlighter.ts token, types
2.11  lib/engine/fuzzy/fuzzy.ts          (none)
```

### Phase 3: App Shell + Theme

Navigation structure, theme system, state management.

```
Step  File                              Depends On
─────────────────────────────────────────────────────────────
3.1   theme/colors.ts                    (none)
3.2   theme/themes.ts                    colors
3.3   theme/tokens.ts                    themes
3.4   services/storage.ts                (none — MMKV)
3.5   store/settings-store.ts            storage, themes
3.6   store/engine-store.ts              engine, storage
3.7   hooks/use-engine.ts                engine-store
3.8   hooks/use-theme.ts                 settings-store, themes
3.9   hooks/use-haptics.ts               (none — expo-haptics)
3.10  app/_layout.tsx                    use-theme, fonts, splash
3.11  app/(tabs)/_layout.tsx             use-theme
```

### Phase 4: Quick Mode

Single-line calculator screen and components.

```
Step  File                              Depends On
─────────────────────────────────────────────────────────────
4.1   hooks/use-eval.ts                  use-engine
4.2   components/keyboard/keyboard-accessory.tsx  (none)
4.3   components/quick/quick-result.tsx  use-theme
4.4   components/quick/quick-input.tsx   use-eval, keyboard-accessory
4.5   components/quick/quick-history.tsx use-theme, use-haptics
4.6   app/(tabs)/index.tsx               quick-input, quick-result, quick-history
```

### Phase 5: Notepad Mode

Multi-line editor with syntax highlighting and inline results.

```
Step  File                              Depends On
─────────────────────────────────────────────────────────────
5.1   components/shared/styled-span.tsx  theme/tokens
5.2   components/shared/styled-line.tsx  styled-span, highlighter
5.3   components/notepad/line-result.tsx use-theme
5.4   components/notepad/notepad-line.tsx styled-line, line-result, use-eval
5.5   components/notepad/notepad-editor.tsx notepad-line, use-engine
5.6   components/notepad/variables-panel.tsx use-engine
5.7   components/shared/explain-view.tsx explain/formatter
5.8   app/explain/[line].tsx             explain-view
5.9   app/(tabs)/notepad.tsx             notepad-editor, variables-panel
```

### Phase 6: Rates + Persistence + Settings

Network layer, caching, session management, settings screen.

```
Step  File                              Depends On
─────────────────────────────────────────────────────────────
6.1   hooks/use-rates.ts                 use-engine, fetch, storage
6.2   components/shared/rate-status-badge.tsx use-rates
6.3   services/session.ts                storage, engine
6.4   hooks/use-session.ts               session, use-engine
6.5   app/(tabs)/settings.tsx            settings-store, use-theme
```

### Phase 7: Tests

Test suite with shared vectors.

```
Step  File                              Depends On
─────────────────────────────────────────────────────────────
7.1   __tests__/fixtures/eval-cases.json (none — data)
7.2   __tests__/engine/types.test.ts     types
7.3   __tests__/engine/lexer.test.ts     lexer
7.4   __tests__/engine/parser.test.ts    parser
7.5   __tests__/engine/eval.test.ts      eval
7.6   __tests__/engine/integration.test.ts engine
```

## File Descriptions

### lib/engine/

| File | Purpose |
|---|---|
| `errors/errors.ts` | Error types with position and context info |
| `types/period.ts` | Time period enum (second through year) with conversion math |
| `types/currency.ts` | Currency struct, registry with 40+ fiat currencies, lookup by code/symbol/alias |
| `types/unit.ts` | Unit struct, 10 categories, 100+ units, linear + temperature conversion |
| `types/crypto.ts` | Crypto struct, 30+ coins, CoinGecko ID mapping |
| `types/metal.ts` | Metal struct, 10 metals (precious + industrial) |
| `types/value.ts` | Discriminated union (11 kinds), constructors, formatting, type compatibility |
| `token/token.ts` | Token type enum, Token interface, keyword and currency symbol maps |
| `lexer/lexer.ts` | Character-level tokenizer with UTF-8, numbers, multi-word identifiers |
| `ast/ast.ts` | All AST nodes as discriminated unions, BinaryOp/UnaryOp, Walk visitor |
| `parser/parser.ts` | Recursive descent with precedence climbing, rate suffix detection |
| `eval/context.ts` | Variable store, previous result, line history, rate cache adapter interface |
| `eval/operators.ts` | Arithmetic dispatch, percentage add/sub, rate-period multiplication |
| `eval/conversions.ts` | Unit/currency/crypto/metal conversion via registries and rate cache |
| `eval/continuations.ts` | Previous-result carry-forward for operator and conversion continuations |
| `eval/functions/registry.ts` | Function name lookup and dispatch to category handlers |
| `eval/functions/math.ts` | 15+ math functions (sqrt, trig, log, pow, factorial) |
| `eval/functions/stats.ts` | 7 statistics functions (sum, avg, min, max, median, stddev, variance) |
| `eval/functions/finance.ts` | 7 finance functions (loan, compound, npv, irr, pmt, fv, pv) |
| `eval/eval.ts` | Main evaluator: AST node dispatch, binary/unary ops, special forms |
| `engine.ts` | Public Engine class: eval, evalPreview, explain, variables, rates |
| `index.ts` | Barrel export of Engine, Value, types, and utilities |

### app/

| File | Purpose |
|---|---|
| `_layout.tsx` | Root layout: theme provider, font loading, splash screen |
| `(tabs)/_layout.tsx` | Bottom tab navigator with Quick, Notepad, Settings tabs |
| `(tabs)/index.tsx` | Quick mode screen |
| `(tabs)/notepad.tsx` | Notepad mode screen |
| `(tabs)/settings.tsx` | Settings screen (theme, precision, NLP toggle) |
| `explain/[line].tsx` | Modal showing step-by-step explanation for a result |

### components/

| File | Purpose |
|---|---|
| `quick/quick-input.tsx` | TextInput with debounced live preview |
| `quick/quick-result.tsx` | Large formatted result display |
| `quick/quick-history.tsx` | Scrollable list of committed results (swipe to delete) |
| `notepad/notepad-editor.tsx` | Multi-line editor container managing line array |
| `notepad/notepad-line.tsx` | Single line with transparent input over styled text |
| `notepad/line-result.tsx` | Right-aligned result for a notepad line |
| `notepad/variables-panel.tsx` | Bottom sheet listing defined variables |
| `shared/styled-line.tsx` | Renders highlighted spans for a line |
| `shared/styled-span.tsx` | Single colored Text element |
| `shared/explain-view.tsx` | Renders explain trace steps |
| `shared/rate-status-badge.tsx` | Shows rate cache age and refresh status |
| `keyboard/keyboard-accessory.tsx` | Row of operator buttons above keyboard |
