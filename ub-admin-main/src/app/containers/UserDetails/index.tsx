/**
 *
 * UserDetails
 *
 */

import AccountBalanceWalletOutlinedIcon from '@material-ui/icons/AccountBalanceWalletOutlined';
import CheckCircleOutlineOutlinedIcon from '@material-ui/icons/CheckCircleOutlineOutlined';
import PersonOutlineIcon from '@material-ui/icons/PersonOutline';
import { translations } from 'locales/i18n';
import React, { memo } from 'react';
import { useTranslation } from 'react-i18next';
import styled from 'styled-components/macro';
import { useInjectReducer, useInjectSaga } from 'utils/redux-injectors';

import UserDetailsTabs from './components';
import { InitialUserDetails } from '../UserAccounts/types';
import DetailsSegment from './components/DetailsSegment';
import WaletSegment from './components/WaletSegment';
import WhiteAddressesSegemnt from './components/WhiteAddressesSegemnt';
import { userDetailsSaga } from './saga';
import { UserDetailsReducer, sliceKey } from './slice';

interface Props {
  initialData: InitialUserDetails;
}

export const UserDetails = memo((props: Props) => {
  useInjectReducer({ key: sliceKey, reducer: UserDetailsReducer });
  useInjectSaga({ key: sliceKey, saga: userDetailsSaga });

  const { t } = useTranslation();

  return (
    <>
      <UserDetailsTabs
        options={[
          {
            title: t(translations.UserAccounts.Details()),
            component: <DetailsSegment data={props.initialData} />,
            icon: <PersonOutlineIcon />,
          },
          {
            title: t(translations.UserAccounts.Wallets()),
            component: <WaletSegment data={props.initialData} />,
            icon: <AccountBalanceWalletOutlinedIcon />,
          },
          {
            title: t(translations.UserAccounts.WhiteAddresses()),
            component: <WhiteAddressesSegemnt data={props.initialData} />,
            icon: <CheckCircleOutlineOutlinedIcon />,
          },
        ]}
      />
    </>
  );
});

const Wrapper = styled.div``;
