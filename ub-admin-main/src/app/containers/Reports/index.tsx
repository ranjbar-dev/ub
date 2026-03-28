/**
 *
 * Reports
 *
 */

import AssignmentIndOutlinedIcon from '@material-ui/icons/AssignmentIndOutlined';
import TextsmsOutlinedIcon from '@material-ui/icons/TextsmsOutlined';
import ViewListOutlinedIcon from '@material-ui/icons/ViewListOutlined';
import UserDetailsTabs from 'app/containers/UserDetails/components';
import { translations } from 'locales/i18n';
import React, { memo } from 'react';
import { useTranslation } from 'react-i18next';
import styled from 'styled-components/macro';
import { useInjectReducer, useInjectSaga } from 'utils/redux-injectors';

import { reportsSaga } from './saga';
import { ReportsReducer, sliceKey } from './slice';
import { InitialUserDetails } from '../UserAccounts/types';
import AdminReports from './components/AdminReports';
import UserLogs from './components/UserLogs';
import WithdrawalComments from './components/WithdrawalComments';

interface Props {
  initialData: InitialUserDetails;
}

export const Reports = memo((props: Props) => {
  const { initialData } = props;
  useInjectReducer({ key: sliceKey, reducer: ReportsReducer });
  useInjectSaga({ key: sliceKey, saga: reportsSaga });

  const { t } = useTranslation();

  return (
    <>
      <Wrapper>
        <UserDetailsTabs
          options={[
            {
              title: t(translations.CommonTitles.AdminReports()),
              component: <AdminReports data={initialData} />,
              icon: <AssignmentIndOutlinedIcon />,
            },
            {
              title: t(translations.CommonTitles.WithdrawalComments()),
              component: <WithdrawalComments data={initialData} />,
              icon: <TextsmsOutlinedIcon />,
            },
            {
              title: t(translations.CommonTitles.UserLogs()),
              component: <UserLogs data={initialData} />,
              icon: <ViewListOutlinedIcon />,
            },
          ]}
        />
      </Wrapper>
    </>
  );
});

const Wrapper = styled.div``;
