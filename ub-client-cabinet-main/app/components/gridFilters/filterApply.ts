// import { GridFilterTypes } from 'containers/App/constants';
import { FilterModel } from 'containers/OrdersPage/types';
import { GridApi } from 'ag-grid-community';

export const FilterApplier = (data: { gridApi: GridApi; message: any }) => {
  const { gridApi, message } = data;
  //!uncomment for live filters
  // let message = data.message;
  const filters: FilterModel = message.payload;
  // if (filters.pair_currency_name) {
  //   let filterInstance = gridApi.getFilterInstance('pair');
  //   let tmp = filters.pair_currency_name.split('-');
  //   if (tmp[0] === tmp[1]) {
  //     filterInstance?.setModel({
  //       type: GridFilterTypes.Contains,
  //       filter: '',
  //     });
  //     gridApi.onFilterChanged();
  //   } else if (tmp[0] == 'all') {
  //     filterInstance?.setModel({
  //       type: GridFilterTypes.Ends_With,
  //       filter: tmp[1],
  //     });
  //     gridApi.onFilterChanged();
  //   } else if (tmp[1] == 'all') {
  //     filterInstance?.setModel({
  //       type: GridFilterTypes.Starts_With,
  //       filter: tmp[0],
  //     });
  //     gridApi.onFilterChanged();
  //   } else {
  //     filterInstance?.setModel({
  //       type: GridFilterTypes.Contains,
  //       filter: filters.pair_currency_name,
  //     });
  //     gridApi.onFilterChanged();
  //   }
  // }
  // if (filters.type) {
  //   let filterInstance =
  //     gridApi.getFilterInstance('side') || gridApi.getFilterInstance('type');
  //   if (filters.type != 'all') {
  //     filterInstance?.setModel({
  //       type: GridFilterTypes.Contains,
  //       filter: filters.type,
  //     });
  //     gridApi.onFilterChanged();
  //   } else {
  //     filterInstance?.setModel({
  //       type: GridFilterTypes.Contains,
  //       filter: '',
  //     });
  //     gridApi.onFilterChanged();
  //   }
  // }
  // if (filters.code) {
  //   let filterInstance = gridApi.getFilterInstance('code');
  //   if (filters.code === 'all') {
  //     gridApi.destroyFilter('code');
  //   } else {
  //     filterInstance?.setModel({
  //       type: GridFilterTypes.Contains,
  //       filter: filters.code,
  //     });
  //     gridApi.onFilterChanged();
  //   }
  // }
  // if (filters.address) {
  //   let filterInstance = gridApi.getFilterInstance('address');
  //   if (filters.address == 'all') {
  //     gridApi.destroyFilter('address');
  //   } else {
  //     filterInstance?.setModel({
  //       type: GridFilterTypes.Contains,
  //       filter: filters.address,
  //     });
  //     gridApi.onFilterChanged();
  //   }
  // }
  // if (filters.dwType) {
  //   let filterInstance = gridApi.getFilterInstance('type');
  //   if (filters.dwType == 'all') {
  //     gridApi.destroyFilter('type');
  //   } else {
  //     filterInstance?.setModel({
  //       type: GridFilterTypes.Contains,
  //       filter: filters.dwType,
  //     });
  //     gridApi.onFilterChanged();
  //   }
  // }
  // if (filters.start_date === filters.end_date && filters.start_date === '') {
  //   gridApi.destroyFilter('createdAt');
  // } else if (filters.start_date && !filters.end_date) {
  //   let filterInstance = gridApi.getFilterInstance('createdAt');
  //   filterInstance?.setModel({
  //     type: GridFilterTypes.Greater_Than_or_Equal,
  //     filter: Number(filters.start_date.replace(/-/g, '').split(' ')[0]),
  //     filterTo: null,
  //   });
  //   gridApi.onFilterChanged();
  // } else if (filters.end_date && !filters.start_date) {
  //   let filterInstance = gridApi.getFilterInstance('createdAt');
  //   filterInstance?.setModel({
  //     type: GridFilterTypes.Less_Than_or_Equal,
  //     filter: Number(filters.end_date.replace(/-/g, '').split(' ')[0]),
  //     filterTo: null,
  //   });
  //   gridApi.onFilterChanged();
  // } else if (filters.end_date && filters.start_date) {
  //   let filterInstance = gridApi.getFilterInstance('createdAt');
  //   filterInstance?.setModel({
  //     condition1: {
  //       type: GridFilterTypes.Less_Than_or_Equal,
  //       filter: Number(filters.end_date.replace(/-/g, '').split(' ')[0]),
  //       filterTo: null,
  //     },
  //     condition2: {
  //       type: GridFilterTypes.Greater_Than_or_Equal,
  //       filter: Number(filters.start_date.replace(/-/g, '').split(' ')[0]),
  //       filterTo: null,
  //     },
  //     operator: 'AND',
  //   });
  //   gridApi.onFilterChanged();
  // }
  if (filters.hideCancelledOrders === false) {
    gridApi.setQuickFilter('');
  } else if (filters.hideCancelledOrders === true) {
    gridApi.setQuickFilter('showNotCancelled');
  }
};
