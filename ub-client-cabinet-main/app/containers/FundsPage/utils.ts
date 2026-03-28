import {Balance} from './types';

export const BalanceArrayFormatter=(data: Balance[]) => {
	const tmp=data;
	for(let i=0;i<tmp.length;i++) {
		tmp[i].totalAmount=Number(tmp[i].totalAmount).toFixed(8);
	}
	return tmp;
};