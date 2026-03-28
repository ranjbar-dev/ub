import { ColDef, RowClickedEvent } from 'ag-grid-community';
import PopupModal from 'app/components/materialModal/modal';
import { SimpleGrid } from 'app/components/SimpleGrid/SimpleGrid';
import { InitialUserDetails } from 'app/containers/UserAccounts/types';
import { FilterArrayElement } from 'locales/types';
import React, { memo, useMemo, useCallback, useState } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { GridNames } from 'services/messageService';

import { BillingActions } from '../slice';
import { selectBillingDepositsData } from '../selectors';
import DepositModal from './DepositModal';
import { DepositSaveData } from '../types';


interface Props {
  data: InitialUserDetails;
  staticRows: ColDef[];
  filters: FilterArrayElement;
}
function BillingDepositsDataGrid(props: Props) {
  const { data, staticRows, filters } = props;
  const dispatch = useDispatch();
  const depositsData = useSelector(selectBillingDepositsData);
  const handleRowClick = useCallback(e => {
    setModalData(e);
    setIsModalOpen(true);
  }, []);
  const [IsModalOpen, setIsModalOpen] = useState(false);
  const [ModalData, setModalData] = useState<RowClickedEvent | undefined>();
  const handleUpdate = (data: DepositSaveData) => {
    dispatch(BillingActions.UpdateDepositsAction(data));
  };
  return (
    <div style={{ width: '100%', position: 'relative' }}>
      {IsModalOpen === true && (
        <PopupModal
          onClose={() => {
            setIsModalOpen(false);
          }}
          isOpen={IsModalOpen}
        >
          <DepositModal
            onSave={handleUpdate}
            onCancel={() => {
              setIsModalOpen(false);
            }}
            row={ModalData!}
          />
        </PopupModal>
      )}
      {useMemo(
        () => (
          <SimpleGrid
            containerId="UserDetailsWindow"
            additionalInitialParams={{ user_id: data.id }}
            arrayFieldName="payments"
            immutableId="id"
            gridName={GridNames.BILLING_DEPOSIT}
            filters={filters}
            userId={data.id}
            onRowClick={handleRowClick}
            initialAction={BillingActions.GetBillingDepositsDataAction}
            externalData={depositsData}
            staticRows={staticRows}
          />
        ),
        [data],
      )}
    </div>
  );
}

export default memo(BillingDepositsDataGrid);
