import { ColDef } from 'ag-grid-community';
import PopupModal from 'app/components/materialModal/modal';
import { SimpleGrid } from 'app/components/SimpleGrid/SimpleGrid';
import { InitialUserDetails } from 'app/containers/UserAccounts/types';
import { INewWindow } from 'app/NewWindowContainer';
import { FilterArrayElement } from 'locales/types';
import React, {
  memo,
  useMemo,
  useCallback,
  useState,
  useEffect,
  useRef,
} from 'react';
import { useDispatch, useSelector } from 'react-redux';
import {
  MessageNames,
  GridNames,
  MessageService,
} from 'services/messageService';

import { BillingActions } from '../slice';
import { selectBillingWithdrawalsData, selectWithdrawalItemDetails } from '../selectors';
import WithdrawModal from './WithdrawModal';


interface Props {
  data: InitialUserDetails;
  staticRows: ColDef[];
  filters: FilterArrayElement;
}
function BillingWithdrawalssDataGrid(props: Props) {
  //  const modalData = useRef({});
  const { data, staticRows, filters } = props;
  const dispatch = useDispatch();
  const withdrawalsData = useSelector(selectBillingWithdrawalsData);
  const withdrawalItemDetails = useSelector(selectWithdrawalItemDetails);
  const isFirstRender = useRef(true);

  const handleRowClick = useCallback(e => {
    //setInitialModalData(e);
    //modalData.current = e.data;
    dispatch(
      BillingActions.GetBillingWithdrawDetailsAction({
        ...e.data,
        user_id: data.id,
      }),
    );
    //setIsModalOpen(true);
  }, []);

  useEffect(() => {
    if (isFirstRender.current) {
      isFirstRender.current = false;
      return;
    }
    if (withdrawalItemDetails) {
      const messageData: INewWindow = {
        name: MessageNames.OPEN_WINDOW,
        payload: {
          component: <WithdrawModal initialModalData={withdrawalItemDetails} />,
          fullTitle:
            withdrawalItemDetails.details.txId +
            withdrawalItemDetails.details.userId +
            withdrawalItemDetails.details.id,
          id:
            withdrawalItemDetails.details.txId +
            withdrawalItemDetails.details.userId +
            withdrawalItemDetails.details.id,
        },
      };
      MessageService.send(messageData);
    }
  }, [withdrawalItemDetails]);

  return (
    <div style={{ width: '100%' }}>
      {/*<PopupModal
        isOpen={IsModalOpen}
        fullScreen={true}
        onClose={() => {
          setIsModalOpen(false);
        }}
      >
        <WithdrawModal initialModalData={InitialModalData} />
      </PopupModal>*/}
      {useMemo(
        () => (
          <SimpleGrid
            containerId="UserDetailsWindow"
            additionalInitialParams={{ user_id: data.id }}
            arrayFieldName="payments"
            immutableId="id"
            gridName={GridNames.Billing_Withdraw}
            filters={filters}
            userId={data.id}
            onRowClick={handleRowClick}
            initialAction={BillingActions.GetBillingWithdrawalsDataAction}
            externalData={withdrawalsData}
            staticRows={staticRows}
          />
        ),
        [data],
      )}
    </div>
  );
}

export default memo(BillingWithdrawalssDataGrid);
