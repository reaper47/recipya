---
title: Introduction
weight: 1
---

This chapter serves as the developer's guide to building Recipya. 

## Fetch the Code

Run the following command to get the code.

```bash
git clone --recurse-submodules https://github.com/reaper47/recipya.git
```

## Dependencies

The following software is required to build the project.

| Software                                  | Version  |
|-------------------------------------------|----------|
| [Go](https://go.dev/dl)                   | 1.22+    |
| [Node.js](https://nodejs.org/en/download) | 20.10.0+ |
| [Task](https://taskfile.dev/)             | latest   |
| [Templ](https://templ.guide/)             | latest   | 
| [Hugo](https://gohugo.io/installation/)   | latest   |

## Recommended CLI Programs

The following lists CLI programs you should install to help you develop the project.

- The [Goose](https://github.com/pressly/goose?tab=readme-ov-file#install) database migration tool

## Technology Stack

| Frontend                                 | Backend                                     |
|------------------------------------------|---------------------------------------------|
| [daisyUI](https://daisyui.com/)          | [Go](https://go.dev/)                       |
| [htmx](https://htmx.org/)                | [SQLite](https://www.sqlite.org/index.html) |
| [_hyperscript](https://hyperscript.org/) |                                             |
| [templ](https://templ.guide/)            |                                             |