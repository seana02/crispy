# Crispy

A personal finance app built around double-entry bookkeeping.

## Features
- [x] Track transactions
- [x] Tag Transactions
- [ ] Charts and Graphs
- [ ] Budgets
- [ ] Subscriptions

## Building
Requires go 1.23.0, Node 15
Built using [Wails](https://wails.io/)

For development, `wails dev` starts in development mode with hot-reload (or `wails dev -tags webkit2_41` if on a Linux machine without `webkit2gtk-4.0`)

To build, `wails build` (with `-tags webkit2_41` if needed)
