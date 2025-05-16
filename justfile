# list possible commands
default:
  @just --list

# This is a Justfile
run:
  mkdocs serve --dirty || echo "Failed to run mkdocs serve, did you run `source venv/bin/activate`?"

run-prod:
  mkdocs serve --strict || echo "Failed to run mkdocs serve, did you run `source venv/bin/activate`?"

# Initialize the project
init:
    python3 -m venv venv 
    source venv/bin/activate
    pip install -r requirements.txt