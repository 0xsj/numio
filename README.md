# numio

A natural language calculator for the terminal. Write math like you think.
```
$100 to turkish lira              →  4,312.90₺
15% tip on $85                    →  97.75
today + 30 days                   →  2026-02-12
mortgage(250000, 6.5, 30)         →  1,580.17
spark(1, 4, 2, 8, 5, 3, 6)        →  ▁▄▂█▅▃▆
```

## Features

**Natural Language Input**
- Write expressions the way you think: `$60 to euros`, `30 days ago`, `15% of 200`
- Fuzzy autocorrect for typos: `squrt(16)` automatically evaluates as `sqrt(16)`

**Live Currency & Crypto Conversion**
- 170+ fiat currencies with live rates
- Cryptocurrencies: BTC, ETH, SOL, and more
- Precious metals: gold, silver, platinum, palladium
- Automatic fallback between multiple rate providers

**Date & Time Math**
- Relative dates: `today`, `tomorrow`, `yesterday`
- Date arithmetic: `today + 30 days`, `next monday`
- Time zones: `3pm PST`, `what time is it in Tokyo`
- Date functions: `daysbetween()`, `age()`, `workdays()`

**Financial Calculations**
- Loan & mortgage payments
- Compound interest, NPV, IRR
- Depreciation (straight-line, MACRS, DDB)
- ROI, CAGR, break-even analysis

**Statistics**
- Central tendency: mean, median, mode, geometric mean
- Dispersion: variance, standard deviation, IQR
- Correlation, regression, z-scores
- Percentiles and quartiles

**Inline Visualizations**
- Sparklines: `spark(1, 4, 2, 8)` → `▁▃▂█`
- Gauges: `gauge(75)` → `████████░░`
- Histograms, trends, dot plots

**Developer Experience**
- Vim-style keybindings in TUI
- Syntax highlighting
- Variables and expressions
- Line continuation

## Installation
```bash
go install github.com/0xsj/numio/cmd/numio-tui@latest
go install github.com/0xsj/numio/cmd/numio-cli@latest
```

Or build from source:
```bash
git clone https://github.com/0xsj/numio.git
cd numio
make build
```

## Usage

### TUI Mode (Interactive)
```bash
numio-tui
```

### CLI Mode (One-shot)
```bash
numio-cli -e "sqrt(144)"
numio-cli -e "\$100 to eur"
```

## Examples

### Basic Math
```
2 + 2                             →  4
15% of 200                        →  30
sqrt(144)                         →  12
2^10                              →  1024
```

### Currency Conversion
```
$100 to EUR                       →  €85.53
100 usd to turkish lira           →  4,312.90₺
1 btc to usd                      →  94,521.00
1 oz gold to usd                  →  2,891.45
```

### Date & Time
```
today                             →  2026-01-13
today + 90 days                   →  2026-04-13
30 days ago                       →  2025-12-14
next friday                       →  2026-01-17
daysbetween(date(2025,1,1), today()) →  377
age(date(1990, 5, 15))            →  35
what time is it in tokyo          →  2026-01-13 20:30:00
```

### Financial
```
# Monthly mortgage payment
mortgage(250000, 6.5, 30)         →  1,580.17

# Compound interest
compound(10000, 7, 10)            →  19,671.51

# Tip calculator
15% tip on $85                    →  97.75
splittip(120, 20, 4)              →  36.00

# Investment
roi(1500, 1000)                   →  50
cagr(10000, 25000, 5)             →  20.11
```

### Statistics
```
avg(85, 90, 78, 92, 88)           →  86.6
median(1, 2, 3, 4, 100)           →  3
stddev(2, 4, 4, 4, 5, 5, 7, 9)    →  2
percentile(90, 1, 2, 3, 4, 5)     →  4.6
```

### Visualizations
```
spark(1, 4, 2, 8, 5, 3, 6)        →  ▁▄▂█▅▃▆
gauge(75)                         →  ████████░░
histogram(1,1,2,2,2,3,3,4,5)      →  ██ ███ ██ █ █
trend(10, 12, 11, 15, 14, 18)     →  ↗ +80.0%
stars(4.5, 5)                     →  ★★★★☆
```

### Variables
```
price = 99.99
tax = 8.25%
total = price + (tax of price)    →  108.24

subtotal = $150
shipping = $12.99
subtotal + shipping in eur        →  €139.47
```

## Function Reference

<details>
<summary><strong>Math Functions</strong></summary>

| Function | Description | Example |
|----------|-------------|---------|
| `abs(x)` | Absolute value | `abs(-5)` → `5` |
| `sqrt(x)` | Square root | `sqrt(16)` → `4` |
| `cbrt(x)` | Cube root | `cbrt(27)` → `3` |
| `pow(x, y)` | Power | `pow(2, 8)` → `256` |
| `log(x)` | Natural logarithm | `log(e)` → `1` |
| `log10(x)` | Base-10 logarithm | `log10(100)` → `2` |
| `sin(x)`, `cos(x)`, `tan(x)` | Trigonometric | `sin(0)` → `0` |
| `round(x, n)` | Round to n decimals | `round(3.14159, 2)` → `3.14` |
| `floor(x)` | Round down | `floor(3.9)` → `3` |
| `ceil(x)` | Round up | `ceil(3.1)` → `4` |
| `factorial(n)` | Factorial | `factorial(5)` → `120` |
| `gcd(a, b)` | Greatest common divisor | `gcd(12, 18)` → `6` |
| `lcm(a, b)` | Least common multiple | `lcm(4, 6)` → `12` |

</details>

<details>
<summary><strong>Date Functions</strong></summary>

| Function | Description | Example |
|----------|-------------|---------|
| `today()` | Current date | `today()` → `2026-01-13` |
| `now()` | Current datetime | `now()` → `2026-01-13 14:30:00` |
| `date(y, m, d)` | Create date | `date(2025, 12, 25)` |
| `adddays(date, n)` | Add days | `adddays(today(), 30)` |
| `addmonths(date, n)` | Add months | `addmonths(today(), 3)` |
| `daysbetween(d1, d2)` | Days between dates | `daysbetween(d1, d2)` |
| `workdaysbetween(d1, d2)` | Business days | `workdaysbetween(d1, d2)` |
| `age(birthdate)` | Age in years | `age(date(1990, 5, 15))` |
| `year(date)` | Extract year | `year(today())` → `2026` |
| `month(date)` | Extract month | `month(today())` → `1` |
| `weekday(date)` | Day of week (0-6) | `weekday(today())` |
| `isweekend(date)` | Check if weekend | `isweekend(today())` |
| `startofmonth(date)` | First of month | `startofmonth(today())` |
| `endofmonth(date)` | Last of month | `endofmonth(today())` |

</details>

<details>
<summary><strong>Finance Functions</strong></summary>

| Function | Description | Example |
|----------|-------------|---------|
| `loan(p, r, y)` | Monthly payment | `loan(250000, 6.5, 30)` |
| `compound(p, r, t)` | Compound interest | `compound(10000, 7, 10)` |
| `simple(p, r, t)` | Simple interest | `simple(1000, 5, 2)` |
| `pv(fv, r, n)` | Present value | `pv(10000, 5, 10)` |
| `fv(pv, r, n)` | Future value | `fv(1000, 5, 10)` |
| `npv(r, cf...)` | Net present value | `npv(10, -1000, 300, 400, 500)` |
| `irr(cf...)` | Internal rate of return | `irr(-1000, 300, 400, 500)` |
| `roi(gain, cost)` | Return on investment | `roi(1500, 1000)` |
| `cagr(start, end, years)` | Compound annual growth | `cagr(10000, 25000, 5)` |
| `sln(cost, salvage, life)` | Straight-line depreciation | `sln(10000, 1000, 5)` |
| `tip(bill, pct)` | Calculate tip | `tip(85, 15)` |
| `splittip(bill, pct, n)` | Split bill with tip | `splittip(100, 20, 4)` |

</details>

<details>
<summary><strong>Statistics Functions</strong></summary>

| Function | Description | Example |
|----------|-------------|---------|
| `avg(...)` | Mean | `avg(1, 2, 3, 4, 5)` |
| `median(...)` | Median | `median(1, 2, 3, 4, 5)` |
| `mode(...)` | Mode | `mode(1, 2, 2, 3)` |
| `min(...)` | Minimum | `min(3, 1, 4, 1, 5)` |
| `max(...)` | Maximum | `max(3, 1, 4, 1, 5)` |
| `sum(...)` | Sum | `sum(1, 2, 3, 4, 5)` |
| `count(...)` | Count | `count(1, 2, 3)` |
| `variance(...)` | Population variance | `variance(2, 4, 4, 4, 5)` |
| `stddev(...)` | Standard deviation | `stddev(2, 4, 4, 4, 5)` |
| `percentile(p, ...)` | Percentile | `percentile(90, 1, 2, 3, 4, 5)` |
| `quartile(q, ...)` | Quartile (1-3) | `quartile(2, 1, 2, 3, 4, 5)` |
| `iqr(...)` | Interquartile range | `iqr(1, 2, 3, 4, 5)` |
| `correlation(x..., y...)` | Correlation coefficient | `corr(1,2,3, 2,4,6)` |
| `zscore(val, ...)` | Z-score | `zscore(7, 2, 4, 4, 4, 5)` |

</details>

<details>
<summary><strong>Visualization Functions</strong></summary>

| Function | Description | Example |
|----------|-------------|---------|
| `spark(...)` | Sparkline | `spark(1,4,2,8)` → `▁▄▂█` |
| `gauge(val, max)` | Progress gauge | `gauge(75)` → `████████░░` |
| `bar(val, max)` | Horizontal bar | `bar(7, 10)` |
| `histogram(...)` | Histogram | `histogram(1,1,2,2,3)` |
| `trend(...)` | Trend indicator | `trend(10, 15, 12, 18)` |
| `stars(val, max)` | Star rating | `stars(4.5, 5)` → `★★★★☆` |
| `battery(pct)` | Battery indicator | `battery(75)` |

</details>

## TUI Keybindings

| Key | Action |
|-----|--------|
| `i` | Enter insert mode |
| `Esc` | Enter normal mode |
| `j` / `k` | Move down / up |
| `gg` / `G` | Go to top / bottom |
| `dd` | Delete line |
| `yy` | Yank line |
| `p` | Paste |
| `o` / `O` | New line below / above |
| `Ctrl+C` | Quit |

## Architecture
```
┌─────────────────────────────────────────────────────────────┐
│                         numio                               │
├─────────────────────────────────────────────────────────────┤
│  TUI (Bubble Tea)  │  CLI  │  [Future: Server, LSP, WASM]  │
├─────────────────────────────────────────────────────────────┤
│                       Engine                                │
│  ┌─────────┐  ┌────────┐  ┌──────┐  ┌───────────────────┐  │
│  │   NLP   │→ │ Lexer  │→ │Parser│→ │     Evaluator     │  │
│  │ Patterns│  │Tokenize│  │  AST │  │ Functions/Context │  │
│  └─────────┘  └────────┘  └──────┘  └───────────────────┘  │
├─────────────────────────────────────────────────────────────┤
│  Types        │  Cache      │  Fetch       │  Fuzzy        │
│  Currency     │  Memory     │  Fiat Rates  │  Levenshtein  │
│  Crypto       │  File       │  Crypto      │  Jaro-Winkler │
│  Metals       │  BFS Path   │  Metals      │  Damerau-Lev  │
│  Units        │             │              │               │
└─────────────────────────────────────────────────────────────┘
```

## Roadmap

### Implemented
- [x] Natural language parsing
- [x] Live currency conversion (170+ currencies)
- [x] Cryptocurrency prices
- [x] Precious metals (gold, silver, platinum, palladium)
- [x] Unit conversion
- [x] Date/time arithmetic
- [x] Financial functions (loan, compound, NPV, IRR, depreciation)
- [x] Statistics functions
- [x] Inline visualizations (sparklines, gauges, histograms)
- [x] Fuzzy autocorrect for typos
- [x] TUI with vim keybindings
- [x] Syntax highlighting
- [x] Multiple rate provider fallbacks

### Planned
- [ ] Stock prices
- [ ] Commodities (oil, natural gas)
- [ ] Boolean/conditional expressions
- [ ] List operations and ranges
- [ ] User-defined functions
- [ ] Data import (CSV, JSON)
- [ ] HTTP server / JSON-RPC API
- [ ] Language Server Protocol (LSP)
- [ ] WASM build for browser
- [ ] Historical rate data
- [ ] Tab completion

## Contributing

Contributions are welcome! Please feel free to submit issues and pull requests.
```bash
# Run tests
make test

# Run with hot reload
make dev

# Build all binaries
make build
```

## License

MIT License - see [LICENSE](LICENSE) for details.

## Acknowledgments

Inspired by [Numi](https://numi.app/) and [numr](https://github.com/timvisee/numr).

---

<p align="center">
  <sub>Built with Go and ☕</sub>
</p>