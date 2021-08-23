//https://github.com/Darth-Knoppix/example-react-fullscreen/blob/master/src/utils/useFullscreenStatus.js
import React from "react";

import env from "./environ"

export default function useFullscreenStatus() {
  const [isFullscreen, setIsFullscreen] = React.useState(calculateFullscreen());

  //safari iOS does not support fullscreen API
  const setFullscreen = (fullscreenEnabled) => {

    if (!hasFullscreen()) return;

    if (fullscreenEnabled) {
      const key = requestFullscreenKey()
      if (document.body[key]) {
        document.body[key]()
      }
    }
    else {
      const key = exitFullscreenKey()
      if (document[key]) document[key]()
    }
  };

  function calculateFullscreen() {
    const elem = document[fullscreenElementKey()]
    return elem === document.body
  }

  function updateIsFullscreen() {
    const status = calculateFullscreen()
    env.log("setIsFullscreen", status)
    setIsFullscreen(status)
  }

  React.useLayoutEffect(() => {
    const key = onfullscreenchangeKey()
    document[key] = updateIsFullscreen
    return () => { document[key] = undefined }
  });

  env.log("hasFullscreen", hasFullscreen())
  return [hasFullscreen(), isFullscreen, setFullscreen];
}

function hasFullscreen() {
  env.log("onfullscreenchangeKey", onfullscreenchangeKey())
  env.log("fullscreenElementKey", fullscreenElementKey())
  env.log("exitFullscreenKey", exitFullscreenKey())
  env.log("requestFullscreenKey", requestFullscreenKey())
  return (typeof document[fullscreenElementKey()] !== "undefined")
    && (typeof document[exitFullscreenKey()] !== "undefined")
    && (typeof document[onfullscreenchangeKey()] !== "undefined")
    && (typeof document.body[requestFullscreenKey()] !== "undefined")
}

function onfullscreenchangeKey() {
  if (typeof document.onfullscreenchange !== "undefined") {
    return "onfullscreenchange";
  } else if (typeof document.onwebkitfullscreenchange !== "undefined") {
    return "onwebkitfullscreenchange";
  } else {
    return "notFound"
  }
}

function fullscreenElementKey() {
  if (typeof document.fullscreenElement !== "undefined") {
    return "fullscreenElement";
  } else if (typeof document.mozFullScreenElement !== "undefined") {
    return "mozFullScreenElement";
  } else if (typeof document.msFullscreenElement !== "undefined") {
    return "msFullscreenElement";
  } else if (typeof document.webkitFullscreenElement !== "undefined") {
    return "webkitFullscreenElement";
  } else {
    return "notFound"
  }
}

function exitFullscreenKey() {
  if (typeof document.exitFullscreen !== "undefined") {
    return "exitFullscreen";
  } else if (typeof document.mozExitFullscreen !== "undefined") {
    return "mozExitFullscreen";
  } else if (typeof document.msExitFullscreen !== "undefined") {
    return "msExitFullscreen";
  } else if (typeof document.webkitExitFullscreen !== "undefined") {
    return "webkitExitFullscreen";
  } else {
    return "notFound"
  }
}

function requestFullscreenKey() {
  if (typeof document.body.requestFullscreen !== "undefined") {
    return "requestFullscreen";
  } else if (typeof document.body.mozRequestFullscreen !== "undefined") {
    return "mozRequestFullscreen";
  } else if (typeof document.body.msRequestFullscreen !== "undefined") {
    return "msRequestFullscreen";
  } else if (typeof document.body.webkitRequestFullscreen !== "undefined") {
    return "webkitRequestFullscreen";
  } else {
    return "notFound"
  }
}
