import { ColDef } from 'ag-grid-community';
import { AgGridReact } from 'ag-grid-react';
import { MqttTopicsPrefixes } from 'containers/App/constants';
import React, { useEffect, useMemo, useState } from 'react';
import { useRef } from 'react';
import { MarketWatchSubscriber } from 'services/message_service';
import { MqttService } from 'services/MqttService2';
import styled from 'styles/styled-components';
import { zeroFixer } from 'utils/formatters';

const OrderList = () => {
  const mqtt2 = useRef(MqttService.getInstance());

  const receivedParData = useRef({});

  useEffect(() => {
    mqtt2.current.ConnectToSubject({
      subject: MqttTopicsPrefixes.MarketWatchAddress,
    });
    return () => {
      mqtt2.current.DisconnectFromSubject({
        subject: MqttTopicsPrefixes.MarketWatchAddress,
      });
    };
  }, []);

  const [apiData, setapiData] = useState<Array<any>>([]);

  useEffect(() => {
    const MarketWatchSubscription = MarketWatchSubscriber.subscribe(
      (message: any) => {
        receivedParData.current[message.payload.name] = message.payload;
        // console.log(receivedParData.current);
        const tmp = Object.values(receivedParData.current);
        setapiData(tmp);
      },
    );
    return () => {
      MarketWatchSubscription.unsubscribe();
    };
  }, []);
  // simple string comparator
  const staticRows = useRef<ColDef[]>([
    { headerName: 'id', field: 'id', hide: true },
    {
      headerName: 'Coin',
      field: 'name',
      suppressMenu: true,
      sortable: true,
      // maxWidth: 90,
    },
    {
      headerName: 'LastPrice',
      field: 'price',
      sort: 'desc',
      suppressMenu: true,
      sortable: true,
      valueFormatter: (params: any) => {
        return zeroFixer(params.data.price);
      },
      comparator: function (a, b) {
        return a - b;
      },
    },
    {
      headerName: 'Change',
      field: 'percentage',
      suppressMenu: true,
      sortable: true,
      cellStyle: function ({ data }) {
        return {
          color: +data.percentage > 0 ? 'var(--greenText)' : 'var(--redText)',
          textAlign: 'start',
        };
      },
      valueFormatter: ({ data }) => {
        return Number(data.percentage).toFixed(2) + '%';
      },
      comparator: function (a, b) {
        return a - b;
      },
    },
  ]);
  const gridConfig = useMemo(
    () => ({
      columnDefs: [...staticRows.current],
    }),
    [],
  );
  return (
    <GridWrapper className="ag-theme-balham">
      <AgGridReact
        animateRows={true}
        columnDefs={gridConfig.columnDefs}
        rowData={apiData}
        headerHeight={24}
        rowHeight={24}
        immutableData={true}
        getRowNodeId={(data) => {
          return data.id.toString();
        }}
      />
    </GridWrapper>
  );
};
const GridWrapper = styled.div`
  height: 450px;
  width: 600px;
  .ag-header {
    opacity: 1;
    span {
      font-size: 10px;
      font-weight: 600;
    }
    div[col-id='percentage'] {
      .ag-header-cell-label {
        justify-content: flex-start;
      }
    }
  }
  .ag-cell {
    line-height: 22px !important;
    font-size: 10px;
    color: var(--blackText);
    font-weight: 600;
  }
  .ag-row {
    cursor: pointer;
  }
  .search {
    position: absolute;
    z-index: 1;
    top: 1px;
    right: 10px;
    .searchInput {
      border: none;
      box-shadow: none;
      outline: none !important;
      border-radius: 5px;
      margin: 0 -27px;
      height: 16px;
      font-size: 10px !important;
      padding: 0 10px;
      width: 100px;
      font-weight: 500;
      background: var(--white);
      color: var(--blackText);
      &::placeholder {
        font-weight: 500;
        font-size: 10px !important;
        font-style: italic;
        color: var(--placeHolderColor);
      }
    }
    svg {
      pointer-events: none;
      transform: scale(0.8);
      margin-top: 0px;
      margin-left: 3px;
    }
  }
`;
export default OrderList;
