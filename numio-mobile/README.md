# numio-mobile

A natural language calculator for iOS and Android — the mobile companion to [numio](https://github.com/0xsj/numio).

## Features

- Natural language math: `rent = 45000TL/month`, `$100 to lira`, `15% of 200`
- 170+ fiat currencies, 30+ cryptocurrencies, precious metals with live rates
- Unit conversions: length, weight, temperature, data, area, volume, speed
- Variables, line continuations, multi-target comparisons
- Two modes: quick single-line calculator and full multi-line notepad
- Syntax highlighting with multiple themes (Dracula, Monokai, Gruvbox, Light)
- Offline-first with cached exchange rates

## Installation

```bash
npm install
```

## Quick Start

```bash
npx expo start
```

Scan the QR code with Expo Go (Android) or Camera (iOS), or press `i` for iOS simulator / `a` for Android emulator.

## Architecture

The app contains a full TypeScript port of the numio engine (`lib/engine/`), with zero React Native dependencies. The engine handles:

```
NLP preprocess -> Lexer (tokenize) -> Parser (AST) -> Evaluator (Value)
```

See [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md) for the full system design.

## Project Structure

See [docs/TREE.md](docs/TREE.md) for the complete file tree with build phases.

## Requirements

- Node.js 18+
- Expo SDK 54
- iOS 15+ / Android 6+
- Expo Go app (for development) or EAS Build (for production)

## License

Private
