import { ColDef } from 'ag-grid-community';
import { SimpleGrid } from 'app/components/SimpleGrid/SimpleGrid';
import { InitialUserDetails } from 'app/containers/UserAccounts/types';
import { FilterArrayElement } from 'locales/types';
import React, { memo } from 'react';
import { useSelector } from 'react-redux';

import { BillingActions } from '../slice';
import { selectBillingAllTransactionsData } from '../selectors';

interface Props {
  data: InitialUserDetails;
  staticRows: ColDef[];
  filters: FilterArrayElement;
}
function BillingAllTransactionsDataGrid(props: Props) {
  const { data, staticRows, filters } = props;
  const allTransactionsData = useSelector(selectBillingAllTransactionsData);

  return (
    <div style={{ width: '100%' }}>
      <SimpleGrid
        containerId="UserDetailsWindow"
        additionalInitialParams={{ user_id: data.id }}
        arrayFieldName="payments"
        immutableId="id"
        filters={filters}
        userId={data.id}
        //onRowClick={handleRowClick}
        initialAction={BillingActions.GetBillingAllTransactionsDataAction}
        externalData={allTransactionsData}
        staticRows={staticRows}
      />
    </div>
  );
}

export default memo(BillingAllTransactionsDataGrid);
