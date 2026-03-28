/**
 *
 * Orders
 *
 */

import FlipToBackOutlinedIcon from '@material-ui/icons/FlipToBackOutlined';
import ReceiptOutlinedIcon from '@material-ui/icons/ReceiptOutlined';
import RestoreOutlinedIcon from '@material-ui/icons/RestoreOutlined';
import UserDetailsTabs from 'app/containers/UserDetails/components';
import { translations } from 'locales/i18n';
import React, { memo } from 'react';
import { useTranslation } from 'react-i18next';
import styled from 'styled-components/macro';
import { useInjectReducer, useInjectSaga } from 'utils/redux-injectors';


import OpenOrders from './components/OpenOrders';
import { InitialUserDetails } from '../UserAccounts/types';
import OrderHistory from './components/OrderHistory';
import TradeHistory from './components/TradeHistory';
import { ordersSaga } from './saga';
import { OrdersReducer, sliceKey } from './slice';

interface Props {
  initialData: InitialUserDetails;
}

export const Orders = memo((props: Props) => {
  useInjectReducer({ key: sliceKey, reducer: OrdersReducer });
  useInjectSaga({ key: sliceKey, saga: ordersSaga });
  const { initialData } = props;
  const { t } = useTranslation();

  return (
    <Wrapper>
      <UserDetailsTabs
        options={[
          {
            title: t(translations.CommonTitles.OpenOrders()),
            component: <OpenOrders data={initialData} />,
            icon: <FlipToBackOutlinedIcon />,
          },
          {
            title: t(translations.CommonTitles.OrderHistory()),
            component: <OrderHistory data={initialData} />,
            icon: <RestoreOutlinedIcon />,
          },
          {
            title: t(translations.CommonTitles.TradeHistory()),
            component: <TradeHistory data={initialData} />,
            icon: <ReceiptOutlinedIcon />,
          },
        ]}
      />
    </Wrapper>
  );
});

const Wrapper = styled.div``;
