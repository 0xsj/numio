# Project Development Guide

> Instructions for AI-assisted development. This document defines how we build software together.

---

## 1. Core Principles

### 1.1 Interface-First Development

- Define contracts before implementations
- Interfaces belong in the consumer package
- Accept interfaces, return concrete types
- Small interfaces (1-3 methods) are preferred
- Every major component should be swappable

### 1.2 Dependency Injection

- Dependencies are explicit, passed via constructors
- No global state, no hidden dependencies
- No magic — if it's used, it's passed in

### 1.3 Leaf-First Build Order

- Start with code that has zero internal dependencies
- Build upward: utilities → domain → application → infrastructure → entry point
- Each file should compile and test before moving to the next

### 1.4 Test Where It Matters

- Business logic: always test
- Complex algorithms: always test
- I/O boundaries: integration test
- Glue code: rarely needs tests
- Offer to write tests when a unit of work is testable

---

## 2. Project Initialization

When starting a new project, create these documents in order:

### Step 1: Discuss the Use Case

Before any code or documentation:

- Clarify the problem being solved
- Identify who uses it and how
- Define success criteria
- Establish what's explicitly out of scope

### Step 2: Create Documentation Structure

```
{{project}}/
├── README.md              # What, why, how to use
├── docs/
│   ├── ARCHITECTURE.md    # System design, component diagrams
│   ├── TREE.md            # File structure with descriptions
│   └── decisions/         # Architecture Decision Records (optional)
├── notes/                 # Learning notes (Obsidian-compatible)
│   ├── language/          # Language-specific concepts
│   ├── patterns/          # Design patterns used
│   ├── domain/            # Domain-specific knowledge
│   └── techniques/        # Implementation techniques
└── CLAUDE.md              # This file (or project-specific variant)
```

### Step 3: Draft Core Documents

1. **README.md** — Overview, installation, usage examples
2. **docs/ARCHITECTURE.md** — ASCII diagrams, components, data flow
3. **docs/TREE.md** — Complete file tree with build phases

Only then begin implementation.

---

## 3. Documentation Standards

### 3.1 README.md Structure

```markdown
# Project Name

One-line description.

## Features

- Bullet points of capabilities

## Installation

[install commands]

## Quick Start

[simplest possible usage]

## Usage

[detailed usage with examples]

## Configuration

[configuration options]

## Requirements

[dependencies, prerequisites]

## License
```

### 3.2 ARCHITECTURE.md Structure

```markdown
# Architecture

## Overview

[One paragraph describing the system]

## System Diagram

[ASCII diagram showing major components and their relationships]

## Components

[For each component:]

- Responsibility (single sentence)
- Key interfaces it defines or implements
- Dependencies

## Data Flow

[How data moves through the system for key operations]

## Key Design Decisions

[Why things are the way they are]
```

### 3.3 TREE.md Structure

```markdown
# Project Structure

## Directory Layout

[Complete ASCII tree of all files/folders]

## Build Phases

[Ordered phases showing which files to build first]

## File Descriptions

[Brief description of each file's purpose]
```

### 3.4 ASCII Diagrams

Use ASCII diagrams for architecture visualization:

```
Component Diagram:
┌─────────────┐     ┌─────────────┐
│  Component  │────▶│  Component  │
└─────────────┘     └─────────────┘

Layered Architecture:
┌─────────────────────────────────┐
│          Entry Point            │
├─────────────────────────────────┤
│         Application             │
├─────────────────────────────────┤
│           Domain                │
├─────────────────────────────────┤
│        Infrastructure           │
└─────────────────────────────────┘

Data Flow:
Input ──▶ Validate ──▶ Process ──▶ Store ──▶ Output

Pipeline:
┌───────┐   ┌───────┐   ┌───────┐
│ Stage │──▶│ Stage │──▶│ Stage │
│   1   │   │   2   │   │   3   │
└───────┘   └───────┘   └───────┘
```

---

## 4. Development Workflow

### 4.1 Per-File Process

1. **Confirm scope** — Single responsibility, clear purpose
2. **Identify interfaces** — What contracts does it need or provide?
3. **Implement** — Keep it minimal, make it work
4. **Test** — If it's testable and non-trivial, write tests
5. **Verify** — No compile errors, tests pass
6. **Document** — Update notes if something new was learned
7. **Confirm** — Get approval before moving to next file

### 4.2 One File at a Time

- Never output multiple files in a single response
- Always confirm before writing the next file
- Exception: multiple files only if explicitly requested

### 4.3 When Stuck or Unclear

- Ask for relevant existing files or context
- Don't assume implementation details
- Clarify requirements before proceeding

---

## 5. Note-Taking System

### 5.1 Purpose

Capture learnings for future reference. Notes are:

- Obsidian-compatible (Markdown with `[[links]]`)
- Organized by category
- Written as the project progresses, not all upfront

### 5.2 Categories

```
notes/
├── language/      # Language syntax, idioms, standard library
├── patterns/      # Design patterns, architectural patterns
├── domain/        # Project-specific domain knowledge
└── techniques/    # Implementation techniques, algorithms
```

### 5.3 When to Create Notes

Create a note when:

- A new concept is used for the first time
- A non-obvious technique solves a problem
- A gotcha or edge case is discovered
- A pattern is applied

### 5.4 Note Format

```markdown
# {{Concept Name}}

## What

[One paragraph explanation]

## Why

[When to use this, what problem it solves]

## Example

[Code or pseudocode demonstrating the concept]

## Gotchas

[Common mistakes, edge cases]

## Related

[[other-note]], [[another-note]]
```

---

## 6. Code Principles (Language-Agnostic)

### 6.1 Structure

- Separate concerns: domain logic vs I/O vs configuration
- Dependencies point inward (infrastructure depends on domain, not vice versa)
- Entry point does wiring only

### 6.2 Interfaces

- Define behavior contracts, not data structures
- Keep them small and focused
- Place in consumer, not provider

### 6.3 Error Handling

- Errors carry context about what failed
- Fail fast on invalid configuration
- Never silently swallow errors

### 6.4 Naming

- Names reveal intent
- Consistent naming conventions throughout
- Avoid abbreviations except for well-known acronyms

### 6.5 Functions

- Do one thing
- Few parameters (consider grouping into a config/options object)
- Return early to avoid deep nesting

---

## 7. Testing Guidelines

### 7.1 What to Test

| Priority | What               | Why                          |
| -------- | ------------------ | ---------------------------- |
| High     | Business logic     | Core value, easy to test     |
| High     | Algorithms         | Complex, error-prone         |
| Medium   | Integration points | Catches interface mismatches |
| Low      | Glue code          | Simple, unlikely to break    |

### 7.2 Test Characteristics

- Fast: unit tests run in milliseconds
- Isolated: no shared state between tests
- Repeatable: same result every time
- Self-validating: pass or fail, no manual inspection

### 7.3 Test Naming

Tests should read like specifications:

- `test_valid_input_returns_expected_output`
- `test_empty_input_raises_error`
- `test_duplicate_entry_is_rejected`

### 7.4 Offering Tests

After implementing a testable unit, offer:

> "This is testable. Want me to write tests for it?"

---

## 8. Patterns Quick Reference

### Creational

| Pattern        | Use When                                      |
| -------------- | --------------------------------------------- |
| Factory        | Multiple implementations of an interface      |
| Builder        | Complex object construction with many options |
| Options/Config | Optional parameters without constructor bloat |

### Structural

| Pattern   | Use When                                       |
| --------- | ---------------------------------------------- |
| Adapter   | Wrapping external code to match your interface |
| Decorator | Adding behavior without modifying original     |
| Facade    | Simplifying a complex subsystem                |

### Behavioral

| Pattern                 | Use When                                    |
| ----------------------- | ------------------------------------------- |
| Strategy                | Swappable algorithms                        |
| Observer                | Event-driven notifications                  |
| Chain of Responsibility | Pipeline processing                         |
| Command                 | Encapsulating operations for undo/queue/log |

### Concurrency

| Pattern        | Use When                              |
| -------------- | ------------------------------------- |
| Worker Pool    | Bounded parallelism                   |
| Pipeline       | Staged concurrent processing          |
| Fan-Out/Fan-In | Parallel work with aggregated results |

---

## 9. Checklist Reference

### Before Starting Implementation

- [ ] Use case discussed and understood
- [ ] README.md drafted
- [ ] ARCHITECTURE.md drafted with ASCII diagrams
- [ ] TREE.md drafted with build phases
- [ ] First file identified

### Per File

- [ ] Scope confirmed (single responsibility)
- [ ] Interfaces identified
- [ ] Implementation complete
- [ ] Tests written (if applicable)
- [ ] Compiles without errors
- [ ] Notes added (if new concept learned)
- [ ] Confirmed ready for next file

### Before Moving On

- [ ] All tests pass
- [ ] Documentation is current
- [ ] No TODOs left unexplained

---

## 10. Communication Style

### Asking for Context

> "I need to see [specific file/directory] to understand [what]."

### Proposing Implementation

> "Here's the plan for [component]:
>
> 1. [Step]
> 2. [Step]
>    Does this approach work?"

### Confirming Before Proceeding

> "Ready to implement [file]. Confirm?"

### Offering Tests

> "This [component] is testable. Want me to write tests?"

### Creating Notes

> "Learned something about [topic]. Adding to notes/[category]/[name].md"

---

## 11. Project Status Tracking

Maintain a `STATUS.md` file at the project root. This file is gitignored and tracks progress organically across sessions.

- Update `STATUS.md` at the end of each session with what was accomplished
- Include: what was built, decisions made, what's next
- Keep it concise — a running log, not a report
- This file is for development continuity, not documentation

---

## 12. Adaptation

This guide is a starting point. As the project evolves:

- Add project-specific sections to CLAUDE.md
- Create additional note categories as needed
- Record major decisions in docs/decisions/
- Update TREE.md as structure changes
