---
title: 'OctoPrint: add Prusa Connect'
description: 'Connect OctoPrint to Prusa Connect'
categories: ['3D Print', 'OctoPrint', 'Prusa']
tags: ['3D Print', 'OctoPrint', 'Prusa']
date: '2025-05-04T21:29:45+02:00'
---

## Introduction

Octoprint is a powerful open-source 3D printer management tool that allows you to control and monitor your 3D printer remotely. Prusa Connect is a cloud-based service that enables you to manage your Prusa printers from anywhere. By connecting OctoPrint to Prusa Connect, you can enhance your 3D printing experience by accessing additional features and functionalities.

## Connecting OctoPrint to Prusa Connect

Follow those steps to connect your OctoPrint instance to Prusa Connect:

- Go to the Cameras section at https://connect.prusa3d.com
- Add a new camera.
- Click the QR code link
- Click "Start Camera"
- Open your browser's inspector window and look for the "/snapshot" request.
- Copy the "Fingerprint" and "Token" headers into the file below.
- Save prusaconnect_upload_cam.sh from below to `/usr/local/bin/prusaconnect_upload_cam.sh`:

```bash
#!/bin/bash

# Set default values for environment variables
: "${HTTP_URL:=https://webcam.connect.prusa3d.com/c/snapshot}"
: "${DELAY_SECONDS:=10}"
: "${LONG_DELAY_SECONDS:=60}"

FINGERPRINT="xxx"
TOKEN="xxx"
SNAPSHOTURL="http://127.0.0.1:8080/?action=snapshot"

while true; do
    # grab from octopi
    curl -s "$SNAPSHOTURL" --output /tmp/output.jpg

    # If no error, upload it.
    if [ $? -eq 0 ]; then
        # POST the image to the HTTP URL using curl
        curl -X PUT "$HTTP_URL" \
            -H "accept: */*" \
            -H "content-type: image/jpg" \
            -H "fingerprint: $FINGERPRINT" \
            -H "token: $TOKEN" \
            --data-binary "@/tmp/output.jpg" \
            --compressed

        # Reset delay to the normal value
        DELAY=$DELAY_SECONDS
    else
        echo "Octopi snapshot returned an error. Retrying after ${LONG_DELAY_SECONDS}s..."

        # Set delay to the longer value
        DELAY=$LONG_DELAY_SECONDS
    fi

    sleep "$DELAY"
done
```

And make it executable:

```bash
chmod +x /usr/local/bin/prusaconnect_upload_cam.sh
```

You can test it by running the script manually:

```bash
/usr/local/bin/prusaconnect_upload_cam.sh
```

You should see the camera image being uploaded to Prusa Connect.

Once done, you can move forward and run it in the background, create `/etc/systemd/system/prusaconnect_upload_cam.service`:

```ini
[Unit]
Description=Octocam to Prusa Connect

[Service]
ExecStart=/usr/local/bin/prusaconnect_upload_cam.sh

[Install]
WantedBy=multi-user.target
```

and then start and enable it:

```bash
systemctl enable --now prusaconnect_upload_cam.service
```

That's it! You should now see your OctoPrint camera feed in Prusa Connect.

## References

- https://gist.github.com/joltcan/bf31bd184b118ee8c983bc1a1fd642af
