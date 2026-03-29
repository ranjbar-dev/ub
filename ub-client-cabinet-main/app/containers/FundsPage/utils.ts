import {Balance} from './types';

export const BalanceArrayFormatter = (data: Balance[]): Balance[] =>
  data.map(b => ({ ...b, totalAmount: Number(b.totalAmount).toFixed(8) }));