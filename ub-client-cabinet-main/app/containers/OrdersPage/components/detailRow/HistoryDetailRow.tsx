import React from 'react';
import { OrderDetail, Order } from 'containers/OrdersPage/types';
import DataRow from 'components/dataRow';
import TextLoader from 'components/textLoader';
import styled from 'styles/styled-components';
import { CurrencyFormater } from 'utils/formatters';

const HistoryDetailRow = (props: { details: OrderDetail; mainData: Order }) => {
  const details = props.details;
  const mainData = props.mainData;
  return (
    <Wrapper>
      <div className='top'>
        <DataRow
          small
          dense
          title='Last updated: '
          titleWidth={'60px'}
          value={
            mainData ? (
              mainData.updatedAt ?? mainData.createdAt
            ) : (
              <TextLoader width={80} />
            )
          }
        />
        <DataRow
          small
          dense
          title='Trading Price: '
          titleWidth={'120px'}
          value={
            mainData ? (
              CurrencyFormater(mainData.price)
            ) : (
              <TextLoader width={80} />
            )
          }
        />
        <DataRow
          small
          dense
          title='Filled: '
          titleWidth={'60px'}
          value={
            mainData ? (
              CurrencyFormater(mainData.executed)
            ) : (
              <TextLoader width={80} />
            )
          }
        />
        <DataRow
          small
          dense
          title='Fee: '
          titleWidth={'60px'}
          value={
            details ? CurrencyFormater(details.fee) : <TextLoader width={80} />
          }
        />
        <DataRow
          small
          dense
          title='Total: '
          titleWidth={'60px'}
          value={
            mainData ? (
              CurrencyFormater(mainData.amount)
            ) : (
              <TextLoader width={80} />
            )
          }
        />
      </div>
    </Wrapper>
  );
};
export default HistoryDetailRow;
const Wrapper = styled.div`
  width: calc(100% - 24px);
  height: calc(100% - 40px);
  position: fixed;
  left: 12px;
  bottom: 0px;
  overflow: hidden;
  pointer-events: none;
  .top {
    /* border-bottom: 1px solid var(--lightGrey); */
    padding-bottom: 10px;
    margin-top: 14px;
  }
  .top,
  .bottom {
    max-height: 37px;
    width: 100%;
    display: flex;
    justify-content: space-between;
    .RowData {
      flex: 1;
      max-width: 200px;
    }
  }
`;
