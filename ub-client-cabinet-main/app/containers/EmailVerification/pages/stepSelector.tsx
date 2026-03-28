import React, { useState, useEffect } from 'react';
import { EmailVerificationPages } from '../constants';
import Verificating from './step1';
import EmailVerified from './step2';
import { Subscriber, MessageNames } from 'services/message_service';
import ErrorStep from './errorStep';

interface Props {}

function StepSelector(props: Props) {
  const [ActivePage, setActivePage]: [EmailVerificationPages, any] = useState(
    EmailVerificationPages.Loading,
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
  const step = (activePage: EmailVerificationPages) => {
    switch (activePage) {
      case EmailVerificationPages.Loading:
        return <Verificating />;
      case EmailVerificationPages.Verified:
        return <EmailVerified />;
      case EmailVerificationPages.Error:
        return <ErrorStep />;
      default:
        return <Verificating />;
    }
  };

  return <>{step(ActivePage)}</>;
}

export default StepSelector;
