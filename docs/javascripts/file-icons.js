/**
 * MkDocs Material File Icons Script
 * 
 * AI Prompt used to generate this script:
 * "Create a JavaScript script for MkDocs Material that adds file icons before tab labels 
 * that contain configuration file extensions. For example, when a tab is labeled 'config.yaml', 
 * add the Material Design file icon (twemoji format) before the text. 
 * Only target configuration files like .yaml, .yml, .json, etc. 
 * Make the script efficient and compatible with MkDocs Material's instant loading feature."
 * 
 * This script efficiently adds file icons to tab labels that have file extensions
 * matching the configExtensions list. It works with MkDocs Material's instant loading
 * and periodically checks for new tabs.
 */

// Configuration files extensions to target
const configExtensions = ['.yaml', '.yml', '.json', '.conf', '.config', '.ini', '.toml', '.sh', '.service'];

// Function to add file icons to tab labels with file extensions
function addFileIcons() {
  // Find all tab labels
  const tabLabels = document.querySelectorAll('label[for^="__tabbed_"]');
  
  tabLabels.forEach(label => {
    // Skip if this label already has an icon
    if (label.querySelector('.twemoji')) return;
    
    const labelText = label.textContent.trim();
    
    // Check if the text ends with one of our target extensions
    for (const ext of configExtensions) {
      if (labelText.endsWith(ext)) {
        // Create the file icon HTML
        const iconHTML = `<span class="twemoji"><svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24"><path d="M13 9V3.5L18.5 9M6 2c-1.11 0-2 .89-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8l-6-6z"></path></svg></span> `;
        
        // Add the icon to the label
        if (label.querySelector('a')) {
          // If there's a link, insert inside it
          const link = label.querySelector('a');
          link.innerHTML = iconHTML + link.innerHTML;
        } else {
          // Otherwise insert directly in the label
          label.innerHTML = iconHTML + label.innerHTML;
        }
        
        break; // Exit the loop once we've added an icon
      }
    }
  });
}

// Run when the page loads
document.addEventListener('DOMContentLoaded', function() {
  // Run immediately
  addFileIcons();
  
  // Run again after a short delay
  setTimeout(addFileIcons, 500);
});

// Handle MkDocs Material instant loading
document.addEventListener('DOMContentSwitch', function() {
  // Run with a slight delay
  setTimeout(addFileIcons, 100);
});

// Run periodically but not too frequently
setInterval(addFileIcons, 2000); 