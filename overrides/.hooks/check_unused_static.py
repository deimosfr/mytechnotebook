#!/usr/bin/env python3
"""
CLI script to check for unused static files in the docs directory.
This script reuses the logic from the static_file_checker.py hook but can be run independently.

Usage:
    python check_unused_static.py [--verbose]

Options:
    --verbose    Show more detailed information about the search process
"""

import os
import re
import sys
import argparse
from pathlib import Path
from typing import List, Set

def find_markdown_files(docs_dir: str) -> List[str]:
    """Find all markdown files in the docs directory."""
    markdown_files = []
    for root, _, filenames in os.walk(docs_dir):
        for filename in filenames:
            if filename.endswith('.md'):
                markdown_files.append(os.path.join(root, filename))
    return markdown_files

def find_static_files(docs_dir: str) -> List[str]:
    """Find all files in the docs/static directory."""
    static_dir = os.path.join(docs_dir, 'static')
    static_files = []
    if os.path.exists(static_dir):
        for root, _, filenames in os.walk(static_dir):
            for filename in filenames:
                static_files.append(os.path.relpath(os.path.join(root, filename), docs_dir))
    return static_files

def check_unused_files(docs_dir: str, verbose: bool = False) -> List[str]:
    """Check for unused static files in the documentation."""
    # Get all markdown files
    markdown_files = find_markdown_files(docs_dir)
    if verbose:
        print(f"Found {len(markdown_files)} markdown files")

    # Get all static files
    static_files = find_static_files(docs_dir)
    if verbose:
        print(f"Found {len(static_files)} static files")

    # Read all markdown content
    markdown_content = []
    for md_file in markdown_files:
        with open(md_file, 'r', encoding='utf-8') as f:
            markdown_content.append(f.read())

    # Check each static file
    unused_files = []
    for static_file in static_files:
        # Convert static file path to patterns that might appear in markdown
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

    return sorted(unused_files)

def main():
    parser = argparse.ArgumentParser(description='Check for unused static files in the docs directory')
    parser.add_argument('--verbose', '-v', action='store_true', help='Show more detailed information')
    args = parser.parse_args()

    # Get the docs directory (assuming we're running from the project root)
    docs_dir = 'docs'
    if not os.path.exists(docs_dir):
        print(f"Error: {docs_dir} directory not found. Make sure you're running this script from the project root.")
        sys.exit(1)

    unused_files = check_unused_files(docs_dir, args.verbose)
    
    if unused_files:
        print("\nUnused static files found:")
        for file in unused_files:
            print(f"WARNING - [unused file]: {file}")
        sys.exit(1)  # Exit with error code if unused files are found
    else:
        print("No unused static files found.")
        sys.exit(0)

if __name__ == '__main__':
    main() 