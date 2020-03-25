window.onload = function() {

  let videoElement = document.querySelector('video');
  
  // Check that browser supports Media Source Extensions API
  if (window.MediaSource) {
    let mediaSource = new MediaSource();
    videoElement.src = URL.createObjectURL(mediaSource);
    mediaSource.addEventListener('sourceopen', sourceOpen);
  } else {
    console.log("Media Source Extensions API is NOT supported");
  }
  
  function sourceOpen(e) {
    URL.revokeObjectURL(videoElement.src);

    let mediaSource = e.target;

    // remote pushes media segments via websocket
    ws = new WebSocket("ws://" + location.hostname + (location.port ? ":"+location.port : "" ) + "/websocket");
    ws.binaryType = "arraybuffer";

    // opened websocket to remote. request initial file segment.
    ws.onopen = function(event) {
      ws.send("please send media segment")
    }

    // The six hexadecimal digit suffix after avc1 is the H.264
    // profile, flags, and level (respectively, one byte each). See
    // ITU-T H.264 specification for details.
    let mime = 'video/mp4; codecs="avc1.4D4028, mp4a.40.2"';
    let sourceBuffer = mediaSource.addSourceBuffer(mime);
    sourceBuffer.addEventListener('updateend', function(e) {
      ws.send("please send media segment");
    });

    // received file or media segment
    ws.onmessage = function(event) {
      sourceBuffer.appendBuffer(event.data);
    }

    // remote closed websocket. end-of-stream.
    ws.onclose = function(event) {
      mediaSource.endOfStream();
    }
  }
};
