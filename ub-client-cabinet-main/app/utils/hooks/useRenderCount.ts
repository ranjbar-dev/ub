
import {useRef} from 'react';

export const useRenderCount = (componentName: string) => {
    const renderTime = useRef(1);
    console.log(componentName + ' rendered: ' + renderTime.current++, renderTime.current == 2 ? 'time' : 'times');
};
