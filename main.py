#!/usr/bin/env python3
import os
import re
import yaml

"""
This is the hook for the variables, macros and filters.
"""
def define_env(env):
    """
    Prompt for Claude:
    Create/update a function named main_index that generates a grid card layout of folders with their subfolders.
    The function should:
    - Generate HTML with a grid of cards, each representing a folder
    - Read folder metadata (title and icon) from the index.md file in each folder
    - For icons, directly convert any '/' in the icon name to '-' in the output (e.g., simple/openbsd → simple-openbsd)
    - Respect spaces at the beginning of lines for proper markdown formatting
    - Read metadata from index.md files using format:
      ```
      ---
      title: Folder Title
      icon: simple/openbsd
      ---
      ```
    - For each folder card, show all its subfolders as links separated by bullets (•)
    - For each subfolder link, use the title from its own index.md with 3 hash marks (###)
    - Skip folders that don't have an index.md file
    - The title should be a link to the index.md of the folder
    - Handle errors gracefully (file not found, parsing errors)
    - Follow the exact format as shown in the example below:
    
    ```
    <div class="grid cards" markdown style="grid-template-columns: repeat(auto-fill, minmax(min(100%, calc(50% - 1rem)), 1fr));">

    -   ### :simple-openbsd: BSD

        ---

        [Filesystems](BSD/Filesystems/index.md) • [Kernel](BSD/Kernel/index.md)
    </div>
    ```
    """
    @env.macro
    def main_index(folder_path):
        base_path = "./docs/" + folder_path
        
        # HTML opening for the grid
        html = '<div class="grid cards" markdown style="grid-template-columns: repeat(auto-fill, minmax(min(100%, calc(50% - 1rem)), 1fr));">\n\n'
        
        # Get immediate subfolders
        try:
            subfolders = [f for f in os.listdir(base_path) if os.path.isdir(os.path.join(base_path, f)) and not f.startswith('.')]
            subfolders.sort()
            
            for folder in subfolders:
                folder_index_path = os.path.join(base_path, folder, 'index.md')
                
                # Skip if index.md doesn't exist
                if not os.path.exists(folder_index_path):
                    continue
                
                # Read metadata from index.md
                title = folder
                icon = None
                
                try:
                    with open(folder_index_path, 'r') as f:
                        content = f.read()
                        # Extract YAML front matter
                        if content.startswith('---'):
                            yaml_match = re.search(r'^---\s+(.*?)\s+---', content, re.DOTALL)
                            if yaml_match:
                                yaml_text = yaml_match.group(1)
                                metadata = yaml.safe_load(yaml_text)
                                
                                # Get title and icon
                                title = metadata.get('title', folder)
                                icon = metadata.get('icon', '')
                
                except Exception as e:
                    print(f"Error reading metadata from {folder_index_path}: {e}")
                
                # Get this folder's subfolders for the links
                current_path = os.path.join(base_path, folder)
                subfolder_links = []
                
                try:
                    child_folders = [f for f in os.listdir(current_path) if os.path.isdir(os.path.join(current_path, f)) and not f.startswith('.')]
                    child_folders.sort()
                    
                    for child in child_folders:
                        child_index_path = os.path.join(current_path, child, 'index.md')
                        child_title = child
                        
                        # Try to get title from child's index.md
                        if os.path.exists(child_index_path):
                            try:
                                with open(child_index_path, 'r') as f:
                                    content = f.read()
                                    yaml_match = re.search(r'^---\s+(.*?)\s+---', content, re.DOTALL)
                                    if yaml_match:
                                        yaml_text = yaml_match.group(1)
                                        metadata = yaml.safe_load(yaml_text)
                                        child_title = metadata.get('title', child)
                            except:
                                pass
                                
                        # Create a link for each subfolder
                        link_path = f"{folder_path}{folder}/{child}/index.md"
                        subfolder_links.append(f"[{child_title}]({link_path})")
                
                except Exception as e:
                    print(f"Error listing subfolders of {current_path}: {e}")
                
                # Build the card for this folder
                # Direct string manipulation for icon formatting to ensure proper rendering
                icon_formatted = icon.replace('/', '-') if icon else ""
                folder_link = f"{folder_path}{folder}/index.md"
                html += f"-   ### :{icon_formatted}: [{title}]({folder_link})\n\n"
                html += "    ---\n\n"
                
                if subfolder_links:
                    html += "    " + " • ".join(subfolder_links) + "\n\n"
                else:
                    html += "    *No subfolders*\n\n"
        
        except Exception as e:
            return f"Error processing directory {base_path}: {e}"
        
        # Close the HTML div
        html += "</div>"
        return html

    """
    Prompt for Claude:
    Create/update a function named sub_index that generates a grid card layout of folders with their subfolders.
    The function should:
    - Generate HTML with a grid of cards, each representing a folder
    - Read folder metadata (title and icon) from the index.md file in each folder
    - For icons, directly convert any '/' in the icon name to '-' in the output (e.g., simple/openbsd → simple-openbsd)
    - Respect spaces at the beginning of lines for proper markdown formatting
    - Read metadata from index.md files using format:
      ```
      ---
      title: Folder Title
      icon: simple/openbsd
      ---
      ```
    - For each folder card, show all its subfolders as links with 4 hash marks (####) like #### :simple-freebsd: [FreeBSD](./Misc/FreeBSD/index.md)
    - For each subfolder link, use the title from its own index.md with 3 hash marks (###)
    - Skip folders that don't have an index.md file
    - The title should be a link to the index.md of the folder
    - For each subfolder, create links with bullet points pointing to every pages, but not index.md
    - Handle errors gracefully (file not found, parsing errors)
    - If there are markdown files (not index.md) at the root of the folder path, display them in a separate grid BEFORE the subfolder grid
    - For root-level files, use a wider grid (90%) with the current folder's title and icon
    - For subfolders, use a 50% width grid layout
    - Follow the exact format as shown in these examples:
    
    For folders with their subfolders:
    ```
    <div class="grid cards" markdown style="grid-template-columns: repeat(auto-fill, minmax(min(100%, calc(50% - 1rem)), 1fr));">

    -   ### :material-harddisk: [Filesystems](./Filesystems/index.md)

        ---

        * [MFS: Using a RAM Filesystem](./Filesystems/mfs_using_a_ram_filesystem.md)
        * [Mounting a Filesystem in Multiple Places Simultaneously](./Filesystems/mounting_a_filesystem_in_multiple_places_simultaneously.md)

    -   ### :material-harddisk: [Filesystems](./Filesystems/index.md)

        ---

        * [MFS: Using a RAM Filesystem](./Filesystems/mfs_using_a_ram_filesystem.md)

        #### :simple-freebsd: [FreeBSD](./Misc/FreeBSD/index.md)

        * [Mounting a Filesystem in Multiple Places Simultaneously](./Filesystems/mounting_a_filesystem_in_multiple_places_simultaneously.md)
    </div>
    ```

    For root-level markdown files, add them before the subfolder grid:
    ```
    <div class="grid cards" markdown style="grid-template-columns: repeat(auto-fill, minmax(min(100%, calc(90% - 1rem)), 1fr));">

    -   ### :material-harddisk: [Misc](./index.md)

        ---

        * [Migration: Migrating Linux accounts to BSD](./migration_migrating_linux_accounts_to_bsd.md)
    </div>
    
    <div class="grid cards" markdown style="grid-template-columns: repeat(auto-fill, minmax(min(100%, calc(50% - 1rem)), 1fr));">

    -   ### :simple-freebsd: [FreeBSD](./FreeBSD/index.md)

        ---

        * [Creating an Apple Time Machine Network on FreeBSD](./FreeBSD/creating_an_apple_time_machine_network_on_freebsd.md)
    </div>
    ```
    """
    @env.macro
    def sub_index(folder_path):
        base_path = "./docs/" + folder_path
        
        # Check if there are any root-level markdown files (excluding index.md)
        root_files = []
        try:
            root_files = [f for f in os.listdir(base_path) 
                         if os.path.isfile(os.path.join(base_path, f)) 
                         and f.endswith('.md') and f != 'index.md']
            root_files.sort()
        except Exception as e:
            print(f"Error listing root files in {base_path}: {e}")
        
        html = ""
        
        # If root files exist, create a separate card section for them
        if root_files:
            # Get folder title and icon from index.md
            folder_title = os.path.basename(folder_path)
            folder_icon = ""
            folder_index_path = os.path.join(base_path, 'index.md')
            
            if os.path.exists(folder_index_path):
                try:
                    with open(folder_index_path, 'r') as f:
                        content = f.read()
                        if content.startswith('---'):
                            yaml_match = re.search(r'^---\s+(.*?)\s+---', content, re.DOTALL)
                            if yaml_match:
                                yaml_text = yaml_match.group(1)
                                metadata = yaml.safe_load(yaml_text)
                                folder_title = metadata.get('title', folder_title)
                                folder_icon = metadata.get('icon', '')
                except Exception as e:
                    print(f"Error reading metadata from {folder_index_path}: {e}")
            
            # Format icon (replace '/' with '-')
            icon_formatted = folder_icon.replace('/', '-') if folder_icon else "material-harddisk"
            
            # Create HTML for root files
            html += '<div class="grid cards" markdown style="grid-template-columns: repeat(auto-fill, minmax(min(100%, calc(90% - 1rem)), 1fr));">\n\n'
            html += f"-   ### :{icon_formatted}: [{folder_title}](./index.md)\n\n"
            html += "    ---\n\n"
            
            # Add links for root files
            root_links = []
            for file in root_files:
                file_title = os.path.splitext(file)[0].replace('_', ' ').title()
                file_path = os.path.join(base_path, file)
                
                try:
                    with open(file_path, 'r') as f:
                        content = f.read()
                        yaml_match = re.search(r'^---\s+(.*?)\s+---', content, re.DOTALL)
                        if yaml_match:
                            yaml_text = yaml_match.group(1)
                            metadata = yaml.safe_load(yaml_text)
                            if 'title' in metadata:
                                file_title = metadata['title']
                except Exception as e:
                    print(f"Error reading metadata from {file_path}: {e}")
                
                file_link = f"./{file}"
                root_links.append(f"    * [{file_title}]({file_link})")
            
            if root_links:
                html += "\n".join(root_links) + "\n\n"
            
            html += "</div>\n"
        
        # Start the regular grid for subfolders
        html += '<div class="grid cards" markdown style="grid-template-columns: repeat(auto-fill, minmax(min(100%, calc(50% - 1rem)), 1fr));">\n\n'
        
        try:
            # Get immediate subfolders
            subfolders = [f for f in os.listdir(base_path) if os.path.isdir(os.path.join(base_path, f)) and not f.startswith('.')]
            subfolders.sort()
            
            for folder in subfolders:
                folder_index_path = os.path.join(base_path, folder, 'index.md')
                folder_path_full = os.path.join(base_path, folder)
                
                # Skip if index.md doesn't exist
                if not os.path.exists(folder_index_path):
                    continue
                
                # Read metadata from index.md
                title = folder
                icon = None
                
                try:
                    with open(folder_index_path, 'r') as f:
                        content = f.read()
                        # Extract YAML front matter
                        if content.startswith('---'):
                            yaml_match = re.search(r'^---\s+(.*?)\s+---', content, re.DOTALL)
                            if yaml_match:
                                yaml_text = yaml_match.group(1)
                                metadata = yaml.safe_load(yaml_text)
                                
                                # Get title and icon
                                title = metadata.get('title', folder)
                                icon = metadata.get('icon', '')
                
                except Exception as e:
                    print(f"Error reading metadata from {folder_index_path}: {e}")
                
                # Format icon (replace '/' with '-')
                icon_formatted = icon.replace('/', '-') if icon else ""
                
                # Create folder card header with icon and title as link with 3 hash marks
                folder_link = f"./{folder}/index.md"
                html += f"-   ### :{icon_formatted}: [{title}]({folder_link})\n\n"
                html += "    ---\n\n"
                
                # Process files in current folder (excluding index.md)
                page_links = []
                
                try:
                    # Get all .md files in the folder except index.md
                    files = [f for f in os.listdir(folder_path_full) 
                           if os.path.isfile(os.path.join(folder_path_full, f)) 
                           and f.endswith('.md') and f != 'index.md']
                    files.sort()
                    
                    for file in files:
                        # Extract title from file content or use filename
                        file_title = os.path.splitext(file)[0].replace('_', ' ').title()
                        file_path = os.path.join(folder_path_full, file)
                        
                        try:
                            with open(file_path, 'r') as f:
                                content = f.read()
                                # Extract YAML front matter for title
                                yaml_match = re.search(r'^---\s+(.*?)\s+---', content, re.DOTALL)
                                if yaml_match:
                                    yaml_text = yaml_match.group(1)
                                    metadata = yaml.safe_load(yaml_text)
                                    if 'title' in metadata:
                                        file_title = metadata['title']
                        except Exception as e:
                            print(f"Error reading metadata from {file_path}: {e}")
                        
                        # Create file link with proper formatting
                        file_link = f"./{folder}/{file}"
                        page_links.append(f"    * [{file_title}]({file_link})")
                
                except Exception as e:
                    print(f"Error listing files in {folder_path_full}: {e}")
                
                # Always add file links before subfolder links
                if page_links:
                    html += "\n".join(page_links) + "\n\n"
                
                # Process subfolders
                try:
                    child_folders = [f for f in os.listdir(folder_path_full) 
                                   if os.path.isdir(os.path.join(folder_path_full, f)) 
                                   and not f.startswith('.')]
                    child_folders.sort()
                    
                    for child in child_folders:
                        child_index_path = os.path.join(folder_path_full, child, 'index.md')
                        
                        # Skip if child doesn't have index.md
                        if not os.path.exists(child_index_path):
                            continue
                        
                        # Get child folder metadata
                        child_title = child
                        child_icon = ""
                        
                        try:
                            with open(child_index_path, 'r') as f:
                                content = f.read()
                                yaml_match = re.search(r'^---\s+(.*?)\s+---', content, re.DOTALL)
                                if yaml_match:
                                    yaml_text = yaml_match.group(1)
                                    metadata = yaml.safe_load(yaml_text)
                                    child_title = metadata.get('title', child)
                                    child_icon = metadata.get('icon', '')
                        except Exception as e:
                            print(f"Error reading metadata from {child_index_path}: {e}")
                        
                        # Format child icon and create header for child folder with 4 hash marks
                        child_icon_formatted = child_icon.replace('/', '-') if child_icon else ""
                        child_link = f"./{folder}/{child}/index.md"
                        
                        # Add the subfolder header with the correct icon format
                        if child_icon_formatted:
                            html += f"    #### :{child_icon_formatted}: [{child_title}]({child_link})\n\n"
                        else:
                            html += f"    #### [{child_title}]({child_link})\n\n"
                        
                        # Process files in child folder (excluding index.md)
                        child_page_links = []
                        child_folder_path = os.path.join(folder_path_full, child)
                        
                        try:
                            child_files = [f for f in os.listdir(child_folder_path) 
                                        if os.path.isfile(os.path.join(child_folder_path, f)) 
                                        and f.endswith('.md') and f != 'index.md']
                            child_files.sort()
                            
                            for file in child_files:
                                # Extract title from file content or use filename
                                file_title = os.path.splitext(file)[0].replace('_', ' ').title()
                                file_path = os.path.join(child_folder_path, file)
                                
                                try:
                                    with open(file_path, 'r') as f:
                                        content = f.read()
                                        yaml_match = re.search(r'^---\s+(.*?)\s+---', content, re.DOTALL)
                                        if yaml_match:
                                            yaml_text = yaml_match.group(1)
                                            metadata = yaml.safe_load(yaml_text)
                                            if 'title' in metadata:
                                                file_title = metadata['title']
                                except Exception as e:
                                    print(f"Error reading metadata from {file_path}: {e}")
                                
                                # Create file link with proper formatting for subfolder files
                                file_link = f"./{folder}/{child}/{file}"
                                child_page_links.append(f"    * [{file_title}]({file_link})")
                        
                        except Exception as e:
                            print(f"Error listing files in {child_folder_path}: {e}")
                        
                        # Add child file links
                        if child_page_links:
                            html += "\n".join(child_page_links) + "\n\n"
                
                except Exception as e:
                    print(f"Error processing subfolders of {folder_path_full}: {e}")
        
        except Exception as e:
            return f"Error processing directory {base_path}: {e}"
        
        # Close the HTML div
        html += "</div>"
        return html
