"""
This hook was generated based on the following prompt:
"Create a hook that generate warning in the mkdocs logs when a file in the docs/static (recursively) is not used in markdown"

The hook scans all markdown files and static files in the docs directory, then checks if each static file
is referenced in any markdown file. If a static file is not referenced, it generates a warning in the
MkDocs logs using the format: "WARNING - [unused file]: path/to/file"
"""

import os
import re
from pathlib import Path
from mkdocs.plugins import event_priority
from mkdocs.structure.files import get_files
from mkdocs.utils import meta
from mkdocs.utils import log

@event_priority(-100)  # Run early in the process
def on_files(files, config):
    """Check for unused static files and generate warnings."""
    # Get all markdown files
    markdown_files = []
    for file in files:
        if file.src_path.endswith('.md'):
            markdown_files.append(file.src_path)

    # Get all static files
    static_dir = Path('docs/static')
    static_files = []
    if static_dir.exists():
        for root, _, filenames in os.walk(static_dir):
            for filename in filenames:
                static_files.append(os.path.relpath(os.path.join(root, filename), 'docs'))

    # Read all markdown content
    markdown_content = []
    for md_file in markdown_files:
        with open(os.path.join('docs', md_file), 'r', encoding='utf-8') as f:
            markdown_content.append(f.read())

    # Check each static file
    unused_files = []
    for static_file in static_files:
        # Convert static file path to a pattern that might appear in markdown
        # This handles both direct references and relative paths
        static_patterns = [
            f'/{static_file}',  # Direct reference
            f'../{static_file}',  # Relative path
            f'../../{static_file}',  # Two levels up
            f'../../../{static_file}',  # Three levels up
            f'../../../../{static_file}',  # Four levels up
        ]
        
        # Check if any pattern is found in any markdown file
        is_used = any(
            any(pattern in content for pattern in static_patterns)
            for content in markdown_content
        )
        
        if not is_used:
            unused_files.append(static_file)

    # Generate warnings for unused files
    for file in sorted(unused_files):
        log.warning(f"[unused file]: {file}")

    return files 