/**
 *
 * FinanceMethods
 *
 */

import ConstructiveModal from 'app/components/ConstructiveModal/ConstructiveModal';
import PopupModal from 'app/components/materialModal/modal';
import { SimpleGrid } from 'app/components/SimpleGrid/SimpleGrid';
import TitledContainer from 'app/components/titledContainer/TitledContainer';
import { FullWidthWrapper } from 'app/components/wrappers/FullWidthWrapper';
import { translations } from 'locales/i18n';
import React, { memo, useMemo, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useDispatch, useSelector } from 'react-redux';
import { MessageNames, GridNames } from 'services/messageService';
import { useInjectReducer, useInjectSaga } from 'utils/redux-injectors';

import { financeMethodsSaga } from './saga';
import {
  FinanceMethodsReducer,
  sliceKey,
  FinanceMethodsActions,
} from './slice';
import { selectFinanceMethodsData } from './selectors';

interface Props { }
export interface IConstructiveModalElement {
  name: string;
  field: keyof IFinanceMethod;
  type?: string;
  editable: boolean;
  options?: {
    name: string;
    value: boolean | string;
  }[];
}
export interface IFinanceMethod {
  code: string;
  depositFee: number;
  id: number;
  maximumWithdraw: string;
  minimumWithdraw: string;
  name: string;
  show_digits?: string;
  supportsDeposit: boolean;
  supportsWithdraw: boolean;
  withdrawalFee: number;
  // currency Pairs
  botOrdersAggregationTime: number;
  botSpread: string;
  botSpreadType: string;

  isActive: boolean;
  makerFee: number;
  maxOurExchangeLimit: string;
  minimumOrderAmount: string;

  ohlcSpread: number;
  showDigits: number;
  takerFee: number;
  tradeStatus: string;
  //
}
export interface ICurrencyPairs { }
export const FinanceMethods = memo((props: Props) => {
  useInjectReducer({ key: sliceKey, reducer: FinanceMethodsReducer });
  useInjectSaga({ key: sliceKey, saga: financeMethodsSaga });
  const [IsModalOpen, setIsModalOpen] = useState(false);
  const [ModalData, setModalData] = useState({});
  const dispatch = useDispatch();
  const { t } = useTranslation();
  const financeMethodsData = useSelector(selectFinanceMethodsData);

  const staticRows = useMemo(
    () => [
      {
        headerName: t(translations.CommonTitles.Name()),
        field: 'name',
      },
      {
        headerName: t(translations.CommonTitles.Currency()),
        field: 'code',
      },
      {
        headerName: t(translations.CommonTitles.ShowDigits()),
        field: 'showDigits',
      },
      {
        headerName: t(translations.CommonTitles.Withdrawals()),
        field: 'supportsWithdraw',
      },
      {
        headerName: t(translations.CommonTitles.Deposits()),
        field: 'supportsDeposit',
      },
      {
        headerName: t(translations.CommonTitles.MinWithdraw()),
        field: 'minimumWithdraw',
      },
      {
        headerName: t(translations.CommonTitles.MaxWithdraw()),
        field: 'maximumWithdraw',
      },
      {
        headerName: t(translations.CommonTitles.WithdrawFee()),
        field: 'withdrawalFee',
      },
      {
        headerName: t(translations.CommonTitles.DepositFee()),
        field: 'depositFee',
      },
    ],

    [t],
  );

  const modalFields: IConstructiveModalElement[] = [
    {
      name: 'Name',
      field: 'name',
      editable: false,
    },
    {
      name: 'Currency',
      field: 'code',
      editable: false,
    },
    {
      name: 'Withdrawal',
      field: 'supportsWithdraw',
      type: 'dropDown',
      editable: true,
      options: [
        {
          name: 'Enabled',
          value: true,
        },
        {
          name: 'Disabled',
          value: false,
        },
      ],
    },
    {
      name: 'Deposit',
      field: 'supportsDeposit',
      editable: true,
      type: 'dropDown',
      options: [
        {
          name: 'Enabled',
          value: true,
        },
        {
          name: 'Disabled',
          value: false,
        },
      ],
    },
    {
      name: 'Show Digits',
      field: 'showDigits',
      editable: true,
    },
    {
      name: 'Min Withdraw',
      field: 'minimumWithdraw',
      editable: true,
    },
    {
      name: 'Max Withdraw',
      field: 'maximumWithdraw',
      editable: true,
    },
    {
      name: 'Withdraw Fee',
      field: 'withdrawalFee',
      editable: true,
    },
    {
      name: 'Deposit Fee',
      field: 'depositFee',
      editable: true,
    },
  ];
  const handleRowClick = ({ data }: { data: Record<string, unknown> }) => {
    setModalData(data);
    setIsModalOpen(true);
  };
  const handleSubmit = (data: Record<string, unknown>) => {
    dispatch(FinanceMethodsActions.UpdateFinanceMethod(data));
  };
  return (
    <FullWidthWrapper>
      <PopupModal
        onClose={() => {
          setIsModalOpen(false);
        }}
        isOpen={IsModalOpen}
      >
        <ConstructiveModal
          onCancel={() => {
            setIsModalOpen(false);
          }}
          onSubmit={handleSubmit}
          // @ts-expect-error — ModalData type (empty object initially) doesn't match ConstructiveModal prop type
          initialData={ModalData}
          modalFields={modalFields}
        />
      </PopupModal>
      <TitledContainer
        id="financeMethods"
        title="Finance Methods"
      >
        <SimpleGrid
          containerId="financeMethods"
          additionalInitialParams={{}}
          gridName={GridNames.FINANCE_METHODS_PAGE}
          arrayFieldName="data"
          flashCellUpdate={true}
          immutableId="id"
          filters={{}}
          onRowClick={handleRowClick}
          initialAction={FinanceMethodsActions.GetFinanceMethods}
          messageName={MessageNames.SET_FINANCEMETHODS_DATA}
          externalData={financeMethodsData}
          staticRows={staticRows}
        />
      </TitledContainer>
    </FullWidthWrapper>
  );
});
