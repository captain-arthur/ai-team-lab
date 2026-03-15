# How to use your CLI task tracker (sample)

This is the **user-facing doc** the Writer would produce so you can get started quickly. (Sample content.)

---

## What you have

A recommendation to use **todo.txt** with a CLI (e.g. **todo.txt-cli** or **topydo**): one file, plain text, minimal setup, terminal-only.

---

## Install

Example (macOS with Homebrew):

```bash
brew install todo-txt-cli
```

Check the official repo for your OS and for the exact package name.

---

## Where your tasks live

By default, the tool uses a file like `~/todo.txt`. You can move this file into a Dropbox folder and point the tool at it via its config so all devices see the same list.

---

## Basic commands

```bash
t add "Your task here"
t ls
t do 1
```

(Exact commands may vary by CLI; see the tool’s help.)

---

## Next steps

1. Install the CLI for your OS.
2. Run `t add "First task"` and `t ls` to confirm.
3. Put `todo.txt` in Dropbox if you want sync.
4. Add aliases or a tiny script if you want shortcuts (e.g. "today" or "work").

For the full reasoning and alternatives, see the rest of `projects/_sample/`.
