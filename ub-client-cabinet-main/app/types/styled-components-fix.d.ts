// Fix styled-components 5 + TypeScript 5 + @types/react 17 JSX element type incompatibility
// See: https://github.com/styled-components/styled-components/issues/3738
import type { ReactElement } from 'react';

declare module 'react' {
  // Widen the JSX.Element key type to include number (matching the Key type)
  // This fixes: "Type 'Key | null' is not assignable to type 'string | null'"
  interface ReactElement<
    P = any,
    T extends string | React.JSXElementConstructor<any> =
      | string
      | React.JSXElementConstructor<any>,
  > {
    key: string | number | null;
  }
}
