import { useGoogleReCaptcha } from 'react-google-recaptcha-v3';
import { useEffect, useRef } from 'react';
import {
  MessageService,
  MessageNames,
  Subscriber,
} from 'services/message_service';
import { LocalStorageKeys } from 'services/constants';
import React from 'react';

export const GoogleRecaptchaComponent = () => {
  const canGet = useRef(true);
  const { executeRecaptcha } = useGoogleReCaptcha();
  const submitHandler = async () => {
    if (!executeRecaptcha) {
      return;
    }
    if (canGet.current === true) {
      canGet.current = false;
      setTimeout(() => {
        canGet.current = true;
      }, 1000);
      const result = await executeRecaptcha('login_page');
      localStorage[LocalStorageKeys.RECAPTCHA] = result;
      MessageService.send({
        name: MessageNames.SET_RECAPTCHA,
        payload: result,
      });
    }
  };
  useEffect(() => {
    const Subscription = Subscriber.subscribe((message: any) => {
      if (message.name === MessageNames.RESET_RECAPTCHA) {
        submitHandler();
      }
    });

    submitHandler();

    return () => {
      Subscription.unsubscribe();
    };
  }, [submitHandler]);

  return <></>;
};
