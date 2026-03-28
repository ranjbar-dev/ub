import React, { useState, useEffect } from 'react';
import { Button } from '@material-ui/core';
import IsLoadingWithText from 'components/isLoadingWithText/isLoadingWithText';
import { Subscriber, MessageNames } from 'services/message_service';

export default function StreamLoadingButton(props) {
  const buttonProps = props;
  const [IsLoading, setIsLoading] = useState(false);
  useEffect(() => {
    const Subscription = Subscriber.subscribe((message: any) => {
      if (message.name === MessageNames.SETLOADING) {
        setIsLoading(message.payload);
      }
    });
    return () => {
      Subscription.unsubscribe();
    };
  }, []);
  return (
    <Button {...buttonProps}>
      <IsLoadingWithText text={props.text} isLoading={IsLoading} />
    </Button>
  );
}
