import { RefObject, useState, useRef, useEffect, useCallback } from 'react';

export const observerErr =
  "💡dimensions: the browser doesn't support Resize Observer, please use polyfill: https://github.com/wellyshen/react-cool-dimensions#resizeobserver-polyfill";
export const borderBoxWarn =
  "💡dimensions: the browser doesn't support border-box size, fallback to content-box size. Please see: https://github.com/wellyshen/react-cool-dimensions#border-box-size-measurement";
interface ResizeObserverSize {
  inlineSize: number;
  blockSize: number;
}

interface ResizeObserverEntry {
  readonly target: Element;
  readonly contentRect: DOMRectReadOnly;
  readonly borderBoxSize: readonly ResizeObserverSize[];
  readonly contentBoxSize: readonly ResizeObserverSize[];
}

interface State {
  currentBreakpoint: string;
  width: number;
  height: number;
  entry?: ResizeObserverEntry;
}
interface Event extends State {
  entry: ResizeObserverEntry;
  observe: () => void;
  unobserve: () => void;
}
interface OnResize {
  (event: Event): void;
}
type Breakpoints = { [key: string]: number };
export interface Options<T> {
  ref?: RefObject<T>;
  useBorderBoxSize?: boolean;
  breakpoints?: Breakpoints;
  onResize?: OnResize;
  polyfill?: typeof ResizeObserver;
}
interface Return<T> {
  ref: RefObject<T>;
  readonly currentBreakpoint: string;
  readonly width: number;
  readonly height: number;
  readonly entry?: ResizeObserverEntry;
  readonly observe: () => void;
  readonly unobserve: () => void;
}

const getCurrentBreakpoint = (bps: Breakpoints, w: number): string => {
  let curBp = '';
  let max = -1;

  Object.keys(bps).forEach((key: string) => {
    const val = bps[key];
    if (w >= val && val > max) {
      curBp = key;
      max = val;
    }
  });

  return curBp;
};

const useDimensions = <T extends HTMLElement>({
  ref: refOpt,
  useBorderBoxSize = false,
  breakpoints,
  onResize,
  polyfill,
}: Options<T> = {}): Return<T> => {
  const [state, setState] = useState<State>({
    currentBreakpoint: '',
    width: 0,
    height: 0,
  });
  const prevSizeRef = useRef<{ width?: number; height?: number }>({});
  const prevBreakpointRef = useRef<string>();
  const isObservingRef = useRef<boolean>(false);
  const observerRef = useRef<ResizeObserver | null>(null);
  const warnedRef = useRef<boolean>(false);
  const onResizeRef = useRef<OnResize | null>(null);
  const refVar = useRef<T>(null);
  const ref = refOpt || refVar;

  useEffect(() => {
    onResizeRef.current = onResize ?? null;
  }, [onResize]);

  const observe = useCallback((): void => {
    if (isObservingRef.current || !observerRef.current || !ref.current) return;

    observerRef.current.observe(ref.current);
    isObservingRef.current = true;
  }, [ref]);

  const unobserve = useCallback((): void => {
    if (!isObservingRef.current || !observerRef.current) return;

    observerRef.current.disconnect();
    isObservingRef.current = false;
  }, []);

  useEffect(() => {
    // @ts-expect-error — ref.current type might not satisfy condition check
    if (!ref.current) return (): void => null;

    if (
      (!('ResizeObserver' in window) || !('ResizeObserverEntry' in window)) &&
      !polyfill
    ) {
      console.error(observerErr);
      // @ts-expect-error — ref.current type might not satisfy condition check
      return (): void => null;
    }

    observerRef.current = new (window.ResizeObserver || polyfill)(
      (entries: ResizeObserverEntry[]) => {
        const entry = entries[0];
        const { contentBoxSize, borderBoxSize, contentRect } = entry;

        const resolvedBoxSize = useBorderBoxSize
          ? (borderBoxSize
            ? (Array.isArray(borderBoxSize) ? borderBoxSize[0] : borderBoxSize)
            : (() => { if (!warnedRef.current) { console.warn(borderBoxWarn); warnedRef.current = true; } return undefined; })())
          : (Array.isArray(contentBoxSize) ? contentBoxSize[0] : contentBoxSize);

        const width = resolvedBoxSize ? resolvedBoxSize.inlineSize : contentRect.width;
        const height = resolvedBoxSize ? resolvedBoxSize.blockSize : contentRect.height;

        if (
          width === prevSizeRef.current.width &&
          height === prevSizeRef.current.height
        )
          return;

        prevSizeRef.current = { width, height };

        const e = {
          currentBreakpoint: '',
          width,
          height,
          entry,
          observe,
          unobserve,
        };

        if (breakpoints) {
          e.currentBreakpoint = getCurrentBreakpoint(breakpoints, width);

          if (e.currentBreakpoint !== prevBreakpointRef.current) {
            if (onResizeRef.current) onResizeRef.current(e);
            prevBreakpointRef.current = e.currentBreakpoint;
          }
        } else if (onResizeRef.current) {
          onResizeRef.current(e);
        }

        setState({
          currentBreakpoint: e.currentBreakpoint,
          width,
          height,
          entry,
        });
      },
    );

    observe();

    return (): void => {
      unobserve();
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [ref, JSON.stringify(breakpoints), observe, unobserve]);

  return { ref, ...state, observe, unobserve };
};

export default useDimensions;
