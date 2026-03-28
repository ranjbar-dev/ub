/*
 *
 * AddressManagementPage
 *
 */

import React, { memo, useEffect } from 'react';

import BreadCrumb from 'components/BreadCrumb';
import { MaxWidthWrapper } from 'components/wrappers/maxWidthWrapper';
import { AppPages } from 'containers/App/constants';
import { Helmet } from 'react-helmet';
import { FormattedMessage } from 'react-intl';
import { useDispatch, useSelector } from 'react-redux';
import { createStructuredSelector } from 'reselect';
import styled from 'styles/styled-components';
import { useInjectReducer } from 'utils/injectReducer';
import { useInjectSaga } from 'utils/injectSaga';

import { Card } from '@material-ui/core';

import translate from './messages';
import reducer from './reducer';
import saga from './saga';
import {
  makeSelectIsLoading,
  makeSelectCurrencies,
  makeSelectAddresses,
  // makeSelectIsAddingAddress,
} from './selectors';
import { initialAction, addNewAddressAction } from './actions';
import AnimateChildren from 'components/AnimateChildren';
import { Currency } from './types';
import CreateWrapper from './components/createWrapper';
import GridFilters from './components/gridFilters';
import DataGrid from './components/dataGrid';
import TitledComponent from 'components/titled';

const stateSelector = createStructuredSelector({
  isLoading: makeSelectIsLoading(),
  currencies: makeSelectCurrencies(),
  addresses: makeSelectAddresses(),
});

interface Props {}

function AddressManagementPage(props: Props) {
  useInjectReducer({ key: 'AddressManagementPage', reducer: reducer });
  useInjectSaga({ key: 'AddressManagementPage', saga: saga });

  const { isLoading, currencies, addresses } = useSelector(stateSelector);
  const dispatch = useDispatch();
  const currencyList: Currency[] = currencies;
  useEffect(() => {
    if (currencyList.length === 0) {
      dispatch(initialAction());
    }
    return () => {};
  }, []);

  const onCreateClick = (data: {
    address: string;
    label: string;
    code: string;
    network?: string;
  }) => {
    dispatch(addNewAddressAction(data));
  };

  return (
    <>
      <Helmet>
        <title>AddressManagement</title>
        <meta
          name="description"
          content="Description of AddressManagementPage"
        />
      </Helmet>
      <MaxWidthWrapper>
        <BreadCrumb
          links={[
            { pageName: 'home', pageLink: AppPages.HomePage },
            {
              pageName: 'acountAndSecurity',
              pageLink: AppPages.AcountPage,
            },
            {
              pageName: 'AddressManagement',
              pageLink: AppPages.AddressManagement,
              last: true,
            },
          ]}
        />
        <AnimateChildren isLoading={isLoading}>
          <TopWrapper>
            <TitledComponent
              title={<FormattedMessage {...translate.Createaddress} />}
            >
              <CreateWrapper
                currencyList={currencyList}
                onCreateClick={onCreateClick}
              />
            </TitledComponent>
          </TopWrapper>
          <BottomWrapper>
            <TitledComponent
              title={<FormattedMessage {...translate.Withdrawaddress} />}
            >
              <GridFilters currencyList={currencyList} />
              <DataGrid data={addresses} />
            </TitledComponent>
          </BottomWrapper>
        </AnimateChildren>
      </MaxWidthWrapper>
    </>
  );
}

export default memo(AddressManagementPage);

const TopWrapper = styled(Card)`
  border-radius: 10px !important;
  box-shadow: none !important;
  height: 12vh;
  min-height: 120px;
  max-height: 120px;
  display: flex;
  flex-direction: column;
  align-items: center;
  margin-bottom: 1vh;
`;

const BottomWrapper = styled(Card)`
  border-radius: 10px !important;
  box-shadow: none !important;
  height: calc(98vh - 245px);
  display: flex;
  flex-direction: column;
  align-items: center;
`;
