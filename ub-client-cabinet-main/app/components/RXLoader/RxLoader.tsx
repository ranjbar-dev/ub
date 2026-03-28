import React, { memo, useEffect, useState } from 'react';
import { GridLoading } from 'components/grid_loading/gridLoading';
import { Subscriber, MessageNames } from 'services/message_service';

interface Props {
  id: string;
  style?: React.CSSProperties;
}

function RxPageLoading(props: Props) {
  const { id, style } = props;
  const [ShowLoading, setShowLoading] = useState(false);
  useEffect(() => {
    const Subscription = Subscriber.subscribe((message: any) => {
      if (
        message.name === MessageNames.SET_PAGE_LOADING_WITH_ID &&
        message.id === id
      ) {
        setShowLoading(message.payload);
      }
    });
    return () => {
      Subscription.unsubscribe();
    };
  }, [id]);

  return <div>{ShowLoading === true ? <GridLoading style={style} /> : ''}</div>;
}

export default memo(RxPageLoading);
