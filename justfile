# List all possible commands
default:
  @just --list

# Run the quick reload server for development
run:
  mkdocs serve --dirty || echo "Failed to run mkdocs serve, did you run `source venv/bin/activate`?"

# Run the production server
run-prod:
  mkdocs serve --strict || echo "Failed to run mkdocs serve, did you run `source venv/bin/activate`?"

# Initialize the project
init:
    python3 -m venv venv 
    source venv/bin/activate
    pip install -r requirements.txt