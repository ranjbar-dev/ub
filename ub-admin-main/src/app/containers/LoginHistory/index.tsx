/**
 *
 * LoginHistory
 *
 */

import { ColDef } from 'ag-grid-community';
import { SimpleGrid } from 'app/components/SimpleGrid/SimpleGrid';
import TitledContainer from 'app/components/titledContainer/TitledContainer';
import { FullWidthWrapper } from 'app/components/wrappers/FullWidthWrapper';
import { translations } from 'locales/i18n';
import React, { memo, useMemo } from 'react';
import { useTranslation } from 'react-i18next';
import { useSelector } from 'react-redux';
import styled from 'styled-components/macro';
import { useInjectReducer, useInjectSaga } from 'utils/redux-injectors';
import { cellColorAndNameFormatter } from 'utils/stylers';

import { loginHistorySaga } from './saga';
import { selectLoginHistoryData } from './selectors';
import { LoginHistoryReducer, sliceKey, LoginHistoryActions } from './slice';


interface Props {}
export const LoginHistory = memo((props: Props) => {
  useInjectReducer({ key: sliceKey, reducer: LoginHistoryReducer });
  useInjectSaga({ key: sliceKey, saga: loginHistorySaga });

  const { t } = useTranslation();
  const loginHistoryData = useSelector(selectLoginHistoryData);

  const staticRows: ColDef[] = useMemo(
    () => [
      {
        headerName: t(translations.CommonTitles.UserName()),
        field: 'email',
      },
      {
        headerName: t(translations.CommonTitles.IP()),
        field: 'ip',
      },
      {
        headerName: t(translations.CommonTitles.Date()),
        field: 'createdAt',
      },
      {
        headerName: t(translations.CommonTitles.State()),
        field: 'type',
        ...cellColorAndNameFormatter('type'),
      },
    ],
    [],
  );

  return (
    <FullWidthWrapper>
      <TitledContainer
        id={'loginHistory'}
        title={t(translations.CommonTitles.LoginHistory())}
      >
        <SimpleGrid
          containerId="loginHistory"
          additionalInitialParams={{}}
          arrayFieldName="userLoginHistory"
          immutableId="id"
          filters={{
            dateCols: ['createdAt'],
            dropDownCols: [
              {
                id: 'type',
                options: [
                  { name: 'successful', value: 'successful' },
                  { name: 'failed', value: 'failed' },
                ],
              },
            ],
          }}
          //  onRowClick={handleRowClick}
          initialAction={LoginHistoryActions.GetLoginHistory}
          externalData={loginHistoryData}
          staticRows={staticRows}
        />
      </TitledContainer>
    </FullWidthWrapper>
  );
});

const Wrapper = styled.div``;
