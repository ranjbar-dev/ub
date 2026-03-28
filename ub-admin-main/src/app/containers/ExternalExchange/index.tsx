/**
 *
 * ExternalExchange
 *
 */

import { SimpleGrid } from 'app/components/SimpleGrid/SimpleGrid';
import TitledContainer from 'app/components/titledContainer/TitledContainer';
import { FullWidthWrapper } from 'app/components/wrappers/FullWidthWrapper';
import { translations } from 'locales/i18n';
import React, { memo, useMemo } from 'react';
import { useTranslation } from 'react-i18next';
import { useSelector } from 'react-redux';
import { MessageNames } from 'services/messageService';
import styled from 'styled-components/macro';
import { useInjectReducer, useInjectSaga } from 'utils/redux-injectors';

import { externalExchangeSaga } from './saga';
import { selectExternalExchangeData } from './selectors';
import {
  ExternalExchangeReducer,
  sliceKey,
  ExternalExchangeActions,
} from './slice';

interface Props {}

export const ExternalExchange = memo((props: Props) => {
  useInjectReducer({ key: sliceKey, reducer: ExternalExchangeReducer });
  useInjectSaga({ key: sliceKey, saga: externalExchangeSaga });

  const { t } = useTranslation();
  const externalExchangeData = useSelector(selectExternalExchangeData);

  const staticRows = useMemo(
    () => [
      {
        headerName: t(translations.CommonTitles.ExchangeName()),
        field: 'name',
      },
      {
        headerName: t(translations.CommonTitles.ExchangeStatus()),
        field: 'status',
      },
      {
        headerName: t(translations.CommonTitles.ExchangeType()),
        field: 'type',
      },
    ],
    [],
  );

  return (
    <FullWidthWrapper>
      <TitledContainer
        id="externalExchange"
        title="External exchange"
      >
        <SimpleGrid
          containerId="externalExchange"
          additionalInitialParams={{}}
          arrayFieldName="data"
          immutableId="id"
          filters={{}}
          initialAction={ExternalExchangeActions.GetExternalExchange}
          messageName={MessageNames.SET_EXTERNAL_EXCHANGE_DATA}
          externalData={externalExchangeData}
          staticRows={staticRows}
        />
      </TitledContainer>
    </FullWidthWrapper>
  );
});

const Wrapper = styled.div``;
