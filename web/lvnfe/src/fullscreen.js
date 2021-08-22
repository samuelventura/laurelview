//https://github.com/Darth-Knoppix/example-react-fullscreen/blob/master/src/utils/useFullscreenStatus.js
import React from "react";

export default function useFullscreenStatus() {
  const [isFullscreen, setIsFullscreen] = React.useState(calculateFullscreen());

  //safari iOS does not support fullscreen API
  const setFullscreen = (fullscreenEnabled) => {
    if (!document.fullscreenEnabled) return;
    if (fullscreenEnabled) {
      document.body
      .requestFullscreen()
      .then(() => {
        setIsFullscreen(calculateFullscreen());
      })
      .catch(() => {
        setIsFullscreen(false);
      });
    }
    else 
    {
      if (calculateFullscreen()) {
        document.exitFullscreen()
      }
    }
  };

  React.useLayoutEffect(() => {
    document.onfullscreenchange = () =>
      setIsFullscreen(calculateFullscreen())
    return () => (document.onfullscreenchange = undefined);
  });

  return [isFullscreen, setFullscreen];
}

function calculateFullscreen() {
  return document[getBrowserFullscreenElementProp()] === document.body
}

function getBrowserFullscreenElementProp() {
  if (typeof document.fullscreenElement !== "undefined") {
    return "fullscreenElement";
  } else if (typeof document.mozFullScreenElement !== "undefined") {
    return "mozFullScreenElement";
  } else if (typeof document.msFullscreenElement !== "undefined") {
    return "msFullscreenElement";
  } else if (typeof document.webkitFullscreenElement !== "undefined") {
    return "webkitFullscreenElement";
  } else {
    return "notfound"
  }
}
