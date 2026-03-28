import React, { memo, useEffect, useState } from 'react';
import { Subscriber, MessageNames, BroadcastMessage } from 'services/messageService';
import NewWindow from 'utils/NW/new-window';
// import NewWindow from 'react-new-window';

interface Props { }
export interface INewWindowPayload {
  id: string;
  component: React.ReactNode;
  fullTitle: string;
  windowWidth?: number;
  windowHeight?: number;
}
export interface INewWindow {
  name: MessageNames;
  payload: INewWindowPayload;
}
function NewWindowContainer(props: Props) {
  const [Windows, setWindows] = useState<Record<string, INewWindowPayload>>({});
  const { } = props;
  useEffect(() => {
    const Subscription = Subscriber.subscribe((message: BroadcastMessage) => {
      const data = message as unknown as INewWindow;
      if (data.name === MessageNames.OPEN_WINDOW && !Windows[data.payload.id]) {
        let windows = { ...Windows };
        windows[data.payload.id] = data.payload;
        setWindows(windows);
      }
      else if (message.name === MessageNames.CLOSE_WINDOW) {
        let windows: Record<string, unknown> = {};
        for (const key in Windows) {
          if (Object.prototype.hasOwnProperty.call(Windows, key)) {
            const element = Windows[key];
            if (data.payload.id !== key) {
              windows[key] = Windows[key]
            }
          }
        }
        setWindows(windows as Record<string, INewWindowPayload>);
      }
    });
    return () => {
      Subscription.unsubscribe();
    };
  }, [Windows]);
  const windowsCreator = (): React.ReactNode[] => {
    let wins: React.ReactNode[] = [];
    for (const key in Windows) {
      if (Object.prototype.hasOwnProperty.call(Windows, key)) {
        const element = (
          <NewWindow
            //  center="parent"
            key={Windows[key].id}
            title={Windows[key].fullTitle}
            features={{
              height: Windows[key].windowHeight ?? 650,
              width: Windows[key].windowWidth ?? 1200,
            }}
            onUnload={() => handleUnLoad(key)}
          >
            {Windows[key].component}
          </NewWindow>
        );
        wins.push(element);
      }
    }
    return wins;
  };
  const handleUnLoad = (key: string) => {
    let windows = { ...Windows };
    delete windows[key];
    setWindows(windows);
  };
  return <div className="cont">{windowsCreator()}</div>;
}

export default memo(NewWindowContainer);
