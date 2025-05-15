# list possible commands
default:
  @just --list

# This is a Justfile
run:
  mkdocs serve --dirty

run-prod:
  mkdocs serve --strict

# Initialize the project
init:
    python3 -m venv venv 
    source venv/bin/activate
    pip install -r requirements.txt