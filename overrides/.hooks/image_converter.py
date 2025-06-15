"""
MkDocs Image Converter Hook

This hook automatically converts images in the docs/static/images directory to AVIF format during the MkDocs build process.
It supports both static images (JPG, PNG) and animated images (GIF).

Features:
- Converts images to AVIF format for better compression and modern browser support
- Preserves animation for GIF files
- Automatically deletes original files after successful conversion
- Handles both static and animated images appropriately
- Provides detailed logging of conversion process and results
- Skips existing AVIF files to speed up subsequent builds
- Resizes large images to reasonable dimensions (max 1920x1080)
- Maintains transparency for PNG files by converting to white background
- Cleans up failed conversions automatically

Requirements:
- pillow-avif-plugin>=1.3.1
- Pillow (PIL)

Usage:
1. Place this file in your MkDocs project's overrides/.hooks/ directory
2. Ensure pillow-avif-plugin is installed in your virtual environment
3. Run mkdocs build or mkdocs serve as normal
4. Images will be automatically converted during the build process

The hook will:
1. Scan the docs/static/images directory for supported image files
2. Convert each image to AVIF format
3. Verify the conversion was successful
4. Delete the original file if conversion succeeds
5. Log the results, including file sizes and compression ratios

Supported formats:
- Input: JPG, JPEG, PNG, GIF (including animated)
- Output: AVIF (static and animated)

Note: The original files will be permanently deleted after successful conversion.
Make sure to keep backups of your original images if needed.
"""

import os
import logging
from pathlib import Path
from PIL import Image
import pillow_avif  # This registers the AVIF format with Pillow

log = logging.getLogger("mkdocs.hooks.image_converter")

def process_image(img, target_size=None):
    """Process image for optimal AVIF conversion."""
    # For animated images, we don't process individual frames
    if getattr(img, "is_animated", False):
        return img

    # Convert to RGB if necessary
    if img.mode in ('RGBA', 'LA'):
        # Create white background
        background = Image.new('RGB', img.size, (255, 255, 255))
        # Paste image using alpha channel as mask
        background.paste(img, mask=img.split()[-1])
        img = background
    elif img.mode != 'RGB':
        img = img.convert('RGB')
    
    # Resize if target size is specified
    if target_size and img.size[0] > target_size[0] or img.size[1] > target_size[1]:
        img.thumbnail(target_size, Image.Resampling.LANCZOS)
    
    return img

def on_pre_build(config):
    """Convert images to AVIF format before building the site."""
    source_dir = Path("docs/static/images")
    skip_existing = True
    max_size = (1920, 1080)  # Maximum dimensions for large images

    if not source_dir.exists():
        log.warning(f"Source directory {source_dir} does not exist")
        return

    for file in source_dir.glob('**/*'):
        if file.is_file() and file.suffix.lower() in ['.jpg', '.jpeg', '.png', '.gif']:
            avif_path = file.with_suffix('.avif')
            
            # Skip if AVIF already exists and skip_existing is True
            if skip_existing and avif_path.exists():
                log.debug(f"Skipping {file} as AVIF already exists")
                continue

            try:
                # Open and process image
                with Image.open(file) as img:
                    # Get original image info
                    original_size = os.path.getsize(file)
                    original_format = img.format
                    original_mode = img.mode
                    is_animated = getattr(img, "is_animated", False)
                    frame_count = getattr(img, "n_frames", 1)
                    
                    if is_animated:
                        log.info(f"Processing animated {original_format} with {frame_count} frames")
                        # For animated images, we save directly without processing
                        img.save(
                            avif_path,
                            format='AVIF',
                            save_all=True  # Required for animated images
                        )
                    else:
                        # Process and save static image
                        processed_img = process_image(img, max_size)
                        processed_img.save(
                            avif_path,
                            format='AVIF'
                        )
                    
                    # Get new file size
                    new_size = os.path.getsize(avif_path)
                    compression_ratio = (1 - (new_size / original_size)) * 100
                    
                    # Verify the AVIF file was created successfully
                    if avif_path.exists() and new_size > 0:
                        try:
                            # Try to delete the original file
                            if file.exists():
                                file.unlink()
                                log.info(f"Successfully deleted original file: {file}")
                            else:
                                log.warning(f"Original file {file} no longer exists")
                            
                            log.info(
                                f"Converted {file} to {avif_path} "
                                f"(Original: {original_format} {original_mode}, "
                                f"{'animated, ' if is_animated else ''}"
                                f"{original_size/1024:.1f}KB, "
                                f"New: AVIF, {new_size/1024:.1f}KB, "
                                f"Compression: {compression_ratio:.1f}%)"
                            )
                        except PermissionError as pe:
                            log.error(f"Permission denied when trying to delete {file}: {str(pe)}")
                        except Exception as e:
                            log.error(f"Failed to delete original file {file}: {str(e)}")
                    else:
                        log.error(f"AVIF conversion failed for {file} - output file is empty or missing")
                        if avif_path.exists():
                            try:
                                avif_path.unlink()  # Clean up failed conversion
                                log.info(f"Cleaned up failed conversion: {avif_path}")
                            except Exception as e:
                                log.error(f"Failed to clean up failed conversion {avif_path}: {str(e)}")
                    
            except Exception as e:
                log.error(f"Failed to convert {file} to AVIF: {str(e)}")
                # Log more details about the error
                if 'img' in locals():
                    log.debug(f"Image details: size={img.size}, mode={img.mode}, format={img.format}, "
                             f"animated={getattr(img, 'is_animated', False)}, "
                             f"frames={getattr(img, 'n_frames', 1)}")
                if avif_path.exists():
                    try:
                        avif_path.unlink()  # Clean up failed conversion
                        log.info(f"Cleaned up failed conversion: {avif_path}")
                    except Exception as e:
                        log.error(f"Failed to clean up failed conversion {avif_path}: {str(e)}") 