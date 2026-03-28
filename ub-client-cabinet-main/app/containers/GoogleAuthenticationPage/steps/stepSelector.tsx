import React, { useState, useEffect } from 'react';
import { QrCode } from '../types';
import Step1 from './step1';
import Step2 from './step2';
import Step3 from './step3';
import { Subscriber, MessageNames } from 'services/message_service';
import DisableStep from './disableStep';

export default function StepSelector (props: {
  qrCode: QrCode;
  userData: any;
  isAuthenticated?: boolean;
}) {
  const goToStep = (step: number, code?: string) => {
    if (code) {
      setG2faCode(code);
    }
    setActiveStep(step);
  };

  useEffect(() => {
    const Subscription = Subscriber.subscribe((message: any) => {
      if (message.name === MessageNames.SET_STEP) {
        goToStep(message.payload);
      }
    });
    return () => {
      Subscription.unsubscribe();
    };
  }, []);
  const [G2faCode, setG2faCode] = useState('');
  const [ActiveStep, setActiveStep] = useState(
    props.isAuthenticated === true ? 3 : 0,
  );
  const selectStep = (step: number) => {
    switch (step) {
      case 0:
        return (
          <Step1
            onNextClick={(code: string) => goToStep(1, code)}
            qrCode={props.qrCode}
          />
        );
      case 1:
        return (
          <Step2
            code={G2faCode}
            userData={props.userData}
            onSubmit={() => {
              goToStep(2);
            }}
            onCancel={() => {
              goToStep(0);
            }}
          />
        );
      case 2:
        return <Step3 />;
      case 3:
        return <DisableStep />;
      default:
        return <Step1 onNextClick={() => goToStep(1)} qrCode={props.qrCode} />;
    }
  };

  return <>{selectStep(ActiveStep)}</>;
}
