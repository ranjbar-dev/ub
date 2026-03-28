import React, { useState, useEffect } from 'react';
import Step1 from './step1';
import Anime from 'react-anime';
import Step2 from './step2';
import { Subscriber, MessageNames } from 'services/message_service';

export default function StepSelector() {
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
  const selectStep = (step: number) => {
    switch (step) {
      case 0:
        return (
          <Anime
            duration={500}
            easing="easeOutCirc"
            scale={[0.1, 1]}
            opacity={[0, 1]}
          >
            <Step1 />
          </Anime>
        );
      case 1:
        return (
          <Anime
            duration={500}
            easing="easeOutCirc"
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
