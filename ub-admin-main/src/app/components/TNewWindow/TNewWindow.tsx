import React from 'react';
import ReactDOM from 'react-dom';

interface Props {
  title: string; // The title of the TNewWindow window
  closeWindow: () => void; // Callback to close the TNewWindow
  /** Window.open features string, e.g. "width=1175,height=745". */
  features: string;
}

interface State {
  externalWindow: Window | null; // The TNewWindow window
  containerElement: HTMLElement | null; // The root element of the TNewWindow window
}

/**
 * Class component that opens a real browser window and portals its children into it.
 * Copies the parent page's stylesheets so the child window matches the app's theme.
 *
 * @example
 * ```tsx
 * <TNewWindow title="User Details" features="width=1175,height=745" closeWindow={() => setOpen(false)}>
 *   <UserDetailsWindow initialData={userData} />
 * </TNewWindow>
 * ```
 */
export default class TNewWindow extends React.Component<Props, State> {
  constructor(props: Props) {
    super(props);
    this.state = {
      externalWindow: null,
      containerElement: null,
    };
  }

  // When we create this component, open a new window
  public componentDidMount() {
    const features = this.props.features;
    const externalWindow = window.open('', '', features);

    let containerElement: HTMLDivElement | null = null;
    if (externalWindow) {
      containerElement = externalWindow.document.createElement('div');
      externalWindow.document.body.appendChild(containerElement);

      // Copy the app's styles into the new window
      const stylesheets = Array.from(document.styleSheets);
      stylesheets.forEach(stylesheet => {
        const css = stylesheet as CSSStyleSheet;

        if (stylesheet.href) {
          const newStyleElement = document.createElement('link');
          newStyleElement.rel = 'stylesheet';
          newStyleElement.href = stylesheet.href;
          externalWindow.document.head.appendChild(newStyleElement);
        } else if (css && css.cssRules && css.cssRules.length > 0) {
          const newStyleElement = document.createElement('style');
          Array.from(css.cssRules).forEach(rule => {
            newStyleElement.appendChild(document.createTextNode(rule.cssText));
          });
          externalWindow.document.head.appendChild(newStyleElement);
        }
      });

      externalWindow.document.title = this.props.title;

      // Make sure the window closes when the component unloads
      externalWindow.addEventListener('beforeunload', () => {
        this.props.closeWindow();
      });
    }

    this.setState({
      externalWindow: externalWindow,
      containerElement: containerElement,
    });
  }

  // Make sure the window closes when the component unmounts
  public componentWillUnmount() {
    if (this.state.externalWindow) {
      this.state.externalWindow.close();
    }
  }

  public render() {
    if (!this.state.containerElement) {
      return null;
    }

    // Render this component's children into the root element of the TNewWindow window
    return ReactDOM.createPortal(
      this.props.children,
      this.state.containerElement,
    );
  }
}
