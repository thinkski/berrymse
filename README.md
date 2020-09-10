# 🍓 BerryMSE

Simple low-latency live video streaming from a Raspberry Pi&trade; via the [Media Source Extensions API](//developer.mozilla.org/en-US/docs/Web/API/Media_Source_Extensions_API).

Note: As of March 2020, Safari on iOS devices still does not support this API (excluding iOS 13 on iPad devices, which do support the API).

## Overview

H.264 Network Abstraction Layer (NAL) units are read from `/dev/video0`, a
Video4Linux2 compatible camera interface. Each unit corresponds to one frame.
Frames are packaged into MPEG-4 ISO BMFF (ISO/IEC 14496-12) compliant
fragments and sent via a websocket to the browser client. The client appends
each received buffer to the media source for playback.

## Quickstart

This demo requires a Raspberry Pi with Camera Module (USB Video Class devices
not currently supported). As it uses the Video4Linux2 interface to access the
camera, the Broadcom v4l2 driver must be installed and the camera must be
enabled in `/boot/config.txt`.

To fetch dependencies:

    GOOS=linux go get -v ./...
    go get github.com/markbates/pkger/cmd/pkger

To build:

    make

To run, copy the appropriate `berrymse` executable to the Raspberry Pi and run:

	./berrymse -l <raspberry pi ip address>:2020

For example:

    ./berrymse -l 192.168.2.1:2020

The Raspberry Pi Zero uses the `armv6l/berrymse` executable. Other models use
the `armv7l/berrymse` executable.

The webpage will show a live video stream with approximately 200ms of latency.
The browser will buffer frames, providing a lookback window.
