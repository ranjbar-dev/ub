/**
*
* ScanBlock
*
*/

import GridTabs from 'app/components/GridTabs/GridTabs';
import TitledContainer from 'app/components/titledContainer/TitledContainer';
import { FullWidthWrapper } from 'app/components/wrappers/FullWidthWrapper';
import { translations } from 'locales/i18n';
import React, { memo, useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { Subscriber, MessageNames } from 'services/messageService';
import styled from 'styled-components/macro';
import { useInjectReducer, useInjectSaga } from 'utils/redux-injectors';

import { scanBlockSaga } from './saga';
import ScanPage from './scanPage'
import { ScanBlockReducer, sliceKey } from './slice';


interface Props { }


export const ScanBlock = memo((props: Props) => {
  useInjectReducer({ key: sliceKey, reducer: ScanBlockReducer });
  useInjectSaga({ key: sliceKey, saga: scanBlockSaga });

  const { t } = useTranslation();

  const tabs = [

    {
      name: t(translations.CommonTitles.ScanBlock()),
      callObject: {
        page: 'scan',
      },
    },

  ];
  const [ActivePage, setActivePage] = useState('scan');

  const handleTopTabChange = (e: Record<string, unknown>) => {
    //console.log(e);
    setActivePage(e.page as string);
  };

  return (
    <FullWidthWrapper>
      <TitledContainer
        id="scanBlock"
        title={t(translations.CommonTitles.ScanBlock())}
      >
        <GridTabs onChange={handleTopTabChange} tabs={tabs} />
        {ActivePage === 'scan' && <ScanPage />}
      </TitledContainer>
    </FullWidthWrapper>

  );

});

const Wrapper = styled.div``;
