[![Build Status](https://travis-ci.com/jschlichtholz/godot.svg?branch=master)](https://travis-ci.com/jschlichtholz/godot)

# godot

Build with `go build`.

Use `source <(curl -s https://raw.githubusercontent.com/jschlichtholz/godot/develop/godotw.sh) init git@github.com:jschlichtholz/.dotfiles.git` to initialize
on any system.

## Prerequisites

You need to have CA certificates installed, such that godot is able to clone the git repository.

## ToDo

* Add module for packages (supporting `apt`, `brew`, and `pacman`)
* Add module executing scripts
* Add module building a global gitignore file using gitignore.io
