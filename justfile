# list possible commands
default:
  @just --list

# This is a Justfile
run:
  mkdocs serve

# Initialize the project
init:
    python3 -m venv venv 
    source venv/bin/activate.fish
    pip install -r requirements.txt