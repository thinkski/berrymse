# μMSE

Light-weight low-latency live video streaming via the [Media Source Extensions API](//developer.mozilla.org/en-US/docs/Web/API/Media_Source_Extensions_API).

Note: As of March 2020, Safari on iOS devices still does not support this API (excluding iOS 13 on iPad devices, which do support the API).

## Overview

The server receives a H.264 Network Abstraction Layer (NAL) unit byte-stream
and wraps the NAL units within MPEG-4 ISO BMFF (ISO/IEC 14496-12) compliant
boxes.

## Quickstart

This demo requires a Raspberry Pi with Camera Module (USB Video Class devices
not currently supported).

To build:

    GOARCH=arm GOOS=linux go build

To run, copy the `micromse` executable together with the `web/` directory to
the Raspberry Pi and run:

	./micromse -l 192.168.1.123:8000

The above is only an example -- use the Raspberry Pi's actual IP address. The
webpage will show a live video stream with about one group-of-pictures (GoP) of
latency (GoP size is 30 frames in the demo, or about one second). The browser
will buffer frames, providing a lookback window.


## Notes

Making a MP4 file for streaming via MSE:

	ffmpeg -pix_fmt uyvy422 -f avfoundation -i "0:0" -t 10 -c:v libx264 -profile:v baseline -level:v 31 -pix_fmt yuv420p -movflags empty_moov+default_base_moof+frag_keyframe sample.mp4
