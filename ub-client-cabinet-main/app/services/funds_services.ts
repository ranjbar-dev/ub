import { apiService } from './api_service';
import { RequestTypes } from './constants';
import { WithdrawModel } from 'containers/FundsPage/types';
export const getBalancesAPI = () => {
  return apiService.fetchData({
    data: { sort: 'desc' },
    url: 'user-balance/balance',
    requestType: RequestTypes.GET,
  });
};
export const getDepositAndWithdrawAPI = (code: string) => {
  return apiService.fetchData({
    data: { code },
    url: 'user-balance/withdraw-deposit',
    requestType: RequestTypes.GET,
  });
};
export const getTransactionHistoryAPI = (data: any) => {
  return apiService.fetchData({
    data: data,
    url: 'crypto-payment',
    requestType: RequestTypes.GET,
  });
};
export const getOrderDetailAPI = (data: { id: number }) => {
  return apiService.fetchData({
    data: data,
    url: 'crypto-payment/detail',
    requestType: RequestTypes.GET,
  });
};

export const getFormerWithdrawAddressesAPI = (code: string) => {
  return apiService.fetchData({
    data: { code },
    url: 'withdraw-address/former-addresses',
    requestType: RequestTypes.GET,
  });
};
export const withdrawAPI = (data: WithdrawModel) => {
  const sendingData = {
    code: data.code,
    amount: data.amount,
    address: data.address,
    label: data.label,
    ...(data.network && { network: data.network }),
    ...(data.email_code && { email_code: data.email_code }),
    '2fa_code': data.G2fa_code,
  };

  return apiService.fetchData({
    data: sendingData,
    url: 'crypto-payment/withdraw',
    requestType: RequestTypes.POST,
  });
};
export const preWithdrawAPI = (data: WithdrawModel) => {
  return apiService.fetchData({
    data: {
      code: data.code,
      amount: data.amount,
      address: data.address,
      ...(data.network && { network: data.network }),
    },
    url: 'crypto-payment/pre-withdraw',
    requestType: RequestTypes.POST,
  });
};
