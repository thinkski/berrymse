# μMSE

Light-weight low-latency live video streaming via the [Media Source Extensions API](//developer.mozilla.org/en-US/docs/Web/API/Media_Source_Extensions_API).

Note: As of March 2020, Safari on iOS devices still does not support this API (excluding iOS 13 on iPad devices, which do support the API).

## Overview

The server receives a H.264 Network Abstraction Layer (NAL) unit byte-stream
and wraps the NAL units within MPEG-4 ISO BMFF (ISO/IEC 14496-12) compliant
boxes.

## Notes

Making a MP4 file for streaming via MSE:

	ffmpeg -pix_fmt uyvy422 -f avfoundation -i "0:0" -t 10 -c:v libx264 -profile:v baseline -level:v 31 -pix_fmt yuv420p -movflags empty_moov+default_base_moof+frag_keyframe sample.mp4
