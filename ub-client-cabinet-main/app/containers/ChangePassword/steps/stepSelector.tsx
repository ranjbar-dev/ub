import React, { useState, useEffect } from 'react';
import Step1 from './step1';
import Step2 from './step2';
import { Subscriber, MessageNames } from 'services/message_service';
import Anime from 'react-anime';

export default function StepSelector () {
  useEffect(() => {
    const Subscription = Subscriber.subscribe((message: any) => {
      if (message.name === MessageNames.SET_STEP) {
        setActiveStep(message.payload);
      }
    });
    return () => {
      Subscription.unsubscribe();
    };
  }, []);

  const [ActiveStep, setActiveStep] = useState(0);
  const selectStep = (stepIndex: number) => {
    switch (stepIndex) {
      case 0:
        return <Step1 />;
      case 1:
        return (
          <Anime
            duration={300}
            easing='easeOutCirc'
            scale={[0.5, 1]}
            opacity={[0, 1]}
          >
            <Step2 />
          </Anime>
        );

      default:
        return <Step1 />;
    }
  };
  return <>{selectStep(ActiveStep)}</>;
}
