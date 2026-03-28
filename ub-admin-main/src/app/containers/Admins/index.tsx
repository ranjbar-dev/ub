/**
 *
 * Admins
 *
 */

import React, { memo, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import { Subscriber, MessageNames, BroadcastMessage } from 'services/messageService';
import styled from 'styled-components/macro';
import { useInjectReducer, useInjectSaga } from 'utils/redux-injectors';

import { adminsSaga } from './saga';
import { AdminsReducer, sliceKey } from './slice';

interface Props {}

export const Admins = memo((props: Props) => {
  useInjectReducer({ key: sliceKey, reducer: AdminsReducer });
  useInjectSaga({ key: sliceKey, saga: adminsSaga });


  const { t } = useTranslation();
  useEffect(() => {
    const Subscription = Subscriber.subscribe((message: BroadcastMessage) => {});
    return () => {
      Subscription.unsubscribe();
    };
  }, []);

  return (
    <>
      <Wrapper>{t('')}</Wrapper>
    </>
  );
});

const Wrapper = styled.div``;
