import React from 'react';
import lottie from 'lottie-web';
import animationData from './LOGO.json';
let animObj: any = null;
export default function SplashAnimation() {
  setTimeout(() => {
    animObj = lottie.loadAnimation({
      container: document.getElementById('SplashContainer') as HTMLElement, // the dom element that will contain the animation
      renderer: 'svg',
      loop: true,
      autoplay: true,
      animationData: animationData, // the path to the animation json
    });
  }, 0);

  return (
    <div>
      <div
        style={{
          width: 200,
          margin: '0 auto',
          height: '64px',
          overflow: 'hidden',
        }}
        id="SplashContainer"
      ></div>
    </div>
  );
}
