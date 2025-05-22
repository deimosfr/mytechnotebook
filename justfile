# List all possible commands
default:
  @just --list

# Run the quick reload server for development
run:
  source venv/bin/activate && mkdocs serve --dirty

# Run the production server
run-prod:
  source venv/bin/activate && mkdocs serve --strict

# Initialize the project
init:
  python3 -m venv venv 
  source venv/bin/activate && pip install -r requirements.txt