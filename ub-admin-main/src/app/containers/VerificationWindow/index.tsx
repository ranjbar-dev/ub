/**
 *
 * VerificationWindow
 *
 */

import React, { memo, useEffect } from 'react';
import { Subscriber, BroadcastMessage } from 'services/messageService';
import styled from 'styled-components/macro';
import { useInjectReducer, useInjectSaga } from 'utils/redux-injectors';

import MainVerificationWrapper from './components/MainVerificationWrapper';
import { verificationWindowSaga } from './saga';
import { VerificationWindowReducer, sliceKey } from './slice';
import { InitialUserDetails } from '../UserAccounts/types';

interface Props {
  initialData: InitialUserDetails;
}

export const VerificationWindow = memo((props: Props) => {
  const { initialData } = props;
  useInjectReducer({ key: sliceKey, reducer: VerificationWindowReducer });
  useInjectSaga({ key: sliceKey, saga: verificationWindowSaga });

  useEffect(() => {
    const Subscription = Subscriber.subscribe((message: BroadcastMessage) => {});
    return () => {
      Subscription.unsubscribe();
    };
  }, []);

  return (
    <>
      <Wrapper className="NWindow">
        <MainVerificationWrapper initialData={initialData} />
      </Wrapper>
    </>
  );
});

const Wrapper = styled.div``;
