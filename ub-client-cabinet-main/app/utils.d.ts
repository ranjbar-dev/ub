declare module '*.svg';

declare module '*.jpg';

declare module '*.png';
declare module '*.gif';

declare namespace JSX {
  interface IntrinsicElements {
    'vaadin-date-picker': any;
  }
}
declare namespace JSX {
  interface IntrinsicElements {
    'dom-module': any;
  }
}
declare module '*.json' {
  const value: any;
  export default value;
}
