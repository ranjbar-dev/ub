import React, { useState, useEffect } from 'react';
import { UpdatePasswordPages } from '../constants';
import UpdatedPasswordPage from './step2';
import { Subscriber, MessageNames } from 'services/message_service';
import ErrorStep from './errorStep';
import Step1 from './step1';

function StepSelector(props: { email: string; code: string }) {
  const [ActivePage, setActivePage]: [UpdatePasswordPages, any] = useState(
    UpdatePasswordPages.ResetPage,
  );
  useEffect(() => {
    const Subscription = Subscriber.subscribe((message: any) => {
      if (message.name === MessageNames.SET_STEP) {
        setActivePage(message.payload);
      }
    });
    return () => {
      Subscription.unsubscribe();
    };
  }, []);
  const step = (activePage: UpdatePasswordPages) => {
    switch (activePage) {
      case UpdatePasswordPages.ResetPage:
        return <Step1 {...props} />;
      case UpdatePasswordPages.UpdatedPage:
        return <UpdatedPasswordPage />;
      case UpdatePasswordPages.Error:
        return <ErrorStep />;
      default:
        return <Step1 {...props} />;
    }
  };

  return <>{step(ActivePage)}</>;
}

export default StepSelector;
