import { ValueFormatterParams, CellClassParams } from 'ag-grid-community';
import { DepositStatusStrings } from 'app/containers/Deposits/types';
import { ConfirmationStatus } from 'app/containers/UserAccounts/types';

export const stateStyler = (state: string): string | undefined => {
  const normalized = state.toLowerCase().replace(/_/g, ' ');

  if (
    ['completed', 'confirmed', 'successful', 'buy'].includes(normalized)
  ) {
    return '#369452';
  }
  if (
    ['created', 'in progress', 'notconfirmed', 'incomplete'].includes(normalized)
  ) {
    return '#B3B3B3';
  }
  if (
    ['rejected', 'reject', 'failed', 'sell'].includes(normalized)
  ) {
    return '#B16567';
  }
  return undefined;
};
export const cellColorAndNameFormatter = (fieldName: string) => {
  return {
    valueFormatter: (params: ValueFormatterParams) => {
      let tmpName = params.data[fieldName];
      if (
        tmpName === DepositStatusStrings.Created ||
        tmpName === ConfirmationStatus.Incomplete ||
        tmpName === ConfirmationStatus.NotConfirmed ||
        tmpName === 'in progress'
      ) {
        tmpName = 'Pending';
      }
      return tmpName ? tmpName.toLowerCase().replace(/_/g, ' ') : '';
    },
    cellStyle: (params: CellClassParams) => {
      return {
        color: stateStyler(params.data[fieldName]),
        textTransform: 'capitalize',
      };
    },
  };
};
