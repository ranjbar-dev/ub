import React, { useEffect } from 'react';
let intervalForWidget: any;

export default function Tawk () {
  useEffect(() => {
    if (!window) {
      throw new Error('DOM is unavailable');
    }
    //@ts-ignore
    window.Tawk_API = window.Tawk_API || {};
    //@ts-ignore
    window.Tawk_LoadStart = new Date();

    const tawk = document.getElementById('tawkId');
    if (tawk) {
      //@ts-ignore
      return window.Tawk_API;
    }
    //

    const script = document.createElement('script');
    script.id = 'tawkId';
    script.async = true;
    script.src = 'https://embed.tawk.to/5fed96ffdf060f156a929c53/1eqs1blon';
    script.charset = 'UTF-8';
    script.setAttribute('crossorigin', '*');

    const first_script_tag = document.getElementsByTagName('script')[0];
    if (!first_script_tag || !first_script_tag.parentNode) {
      throw new Error('DOM is unavailable');
    }

    first_script_tag.parentNode.insertBefore(script, first_script_tag);
  }, []);
  //  const insertAfter = (newNode: any, referenceNode: any) => {
  //    referenceNode.parentNode.insertBefore(newNode, referenceNode.nextSibling)
  //  }
  useEffect(() => {
    let i = 0;
    intervalForWidget = setInterval(function () {
      //  if (i === 0) {
      const iframe = document.getElementsByTagName('iframe');
      if (iframe[0] && iframe[1]) {
        const contentContainer = iframe[0].contentWindow?.document.body.querySelectorAll(
          '#contentContainer',
        );
        if (contentContainer && contentContainer[0]) {
          i++;

          const chatSvg = iframe[1].contentWindow?.document.body.querySelectorAll(
            '#maximizeChat',
          );
          //@ts-ignore
          chatSvg[0].innerHTML = `			
			<svg id="appliIcon" data-name="Layer 1" xmlns="http://www.w3.org/2000/svg" style="width:38px;height:38px" viewBox="0 0 30 30">
<path id="Path_7551" data-name="Path 7551" d="M8.748,5.926l-3.719,3.7a4.763,4.763,0,0,0,0,6.758L9.38,20.72a.131.131,0,0,0,.186,0l10.18-10.128a.13.13,0,0,1,.223.093V14.5A2.605,2.605,0,0,1,19.2,16.35L9.6,25.932a.131.131,0,0,1-.186,0L2.51,19.016a8.519,8.519,0,0,1,0-12.091L9.3.177a.13.13,0,0,1,.223.093V4.084A2.605,2.605,0,0,1,8.748,5.926Z" transform="translate(0.007 0.029)" fill="white"></path>
<path id="Path_7552" data-name="Path 7552" d="M21.451,20l3.821-3.8a4.615,4.615,0,0,0,0-6.48l-4.49-4.481a.131.131,0,0,0-.186,0l-10.143,10.1a.13.13,0,0,1-.223-.093v-3.8A2.605,2.605,0,0,1,11,9.588L20.6,0a.131.131,0,0,1,.186,0l6.982,6.952a8.418,8.418,0,0,1,0,11.933l-6.852,6.823a.13.13,0,0,1-.223-.093V21.854A2.605,2.605,0,0,1,21.451,20Z" transform="translate(-0.713 0.042)" fill="white"></path>
</svg>						
			
			`;

          const aTags = contentContainer[0].getElementsByClassName('emojione');
          const branding = aTags[0].parentNode?.parentNode;

          //if (!(i > 1)) {
          //  let span = document.createElement('span')
          //  span.setAttribute(
          //    'style',
          //    `position: fixed;
          //  bottom: 55px;
          //  z-index: 99;
          //  left: calc(50% - 48px);
          //  `
          //  )
          //  let image = document.createElement('img')
          //  image.setAttribute('style', `width:95px;`)
          //  image.setAttribute(
          //    'src',
          //    'https://app.unitedbit.com/assets/images/email-images/logo.png'
          //  )
          //  span.appendChild(image)
          //  insertAfter(span, branding)
          //}
          //@ts-ignore
          branding.setAttribute(
            'style',
            `opacity:0;pointer-events:none;max-height:0px;`,
          );
          //clearInterval(intervalForWidget)
        }
      }
      //  }
    }, 500);
    return () => {
      clearInterval(intervalForWidget);
    };
  }, []);
  return <></>;
}
