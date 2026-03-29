import React, {
  useEffect,
  useState,
  useRef,
  memo,
  useCallback,
  useLayoutEffect,
} from 'react';
import { LocalStorageKeys } from 'services/constants';
import {
  LayoutWidth,
  layoutMargin,
  LayoutHeight,
} from 'containers/TradePage/utils/tradeUtilities';
import TradeHeader from 'containers/TradePage/components/TradeHeader';
import NewOrder from 'containers/TradePage/components/newOrder';
import GridLayout from 'react-grid-layout';
import MiniTitledComponent from 'components/miniTitledComponent';
import { FormattedMessage } from 'react-intl';
import translate from 'containers/TradePage/messages';
import MarketWatchGrid from 'containers/TradePage/components/marketWatch/marketWatchGrid';
import MarketTradeGrid from 'containers/TradePage/components/marketTradeGrid';
import OrderBookGrid from 'containers/TradePage/components/orderBook/orderBookGrid';
import {
  MessageNames,
  MessageService,
  Subscriber,
} from 'services/message_service';
import useForceUpdate from 'utils/hooks/useForceUpdate';
import styled from 'styles/styled-components';
import Orders from 'containers/TradePage/components/orders';
import TradeChart from 'containers/TradePage/components/TradeChart';
import PopupModal from 'components/materialModal/modal';
import LoginPage from 'containers/LoginPage/Loadable';
import { layout, LayoutContainers } from './layout';
import { CentrifugoChannels } from 'containers/App/constants';

let canSave = true;
let dragTimeout;

const onDragStart = e => {
  canSave = false;
  dragTimeout = setTimeout(() => {
    canSave = true;
  }, 100);
  //to prevent main chart being a bugger when dragging stuff on it
  const chartWrapper = document.getElementById('mainChartWrapper');
  if (chartWrapper) {
    chartWrapper.style.pointerEvents = 'none';
  }
};

const redoLayOutAnimation = () => {
  const elements = document.querySelectorAll('.react-grid-item');
  elements.forEach(item => {
    item.classList.remove('animated');
  });
  setTimeout(() => {
    const elements = document.querySelectorAll('.react-grid-item');
    elements.forEach(item => {
      item.classList.add('animated');
    });
  }, 1000);
};

const TradeLayout = () => {
  const handleOnResize = useCallback((e, layout) => {
    localStorage[LocalStorageKeys.TRADELAYOUT] = JSON.stringify(e);
    MessageService.send({ name: MessageNames.LAYOUT_RESIZE, payload: layout });
    MessageService.send({ name: MessageNames.LAYOUT_CHANGE });
  }, []);
  const handleDrag = useCallback((e, layout) => {
    clearTimeout(dragTimeout);
    if (canSave === true) {
      localStorage[LocalStorageKeys.TRADELAYOUT] = JSON.stringify(e);
      MessageService.send({ name: MessageNames.LAYOUT_CHANGE });
    } else {
      canSave = true;
    }
  }, []);

  const forceUpdate = useForceUpdate();
  const [isLoginPopupOpen, setisLoginPopupOpen] = useState(false);
  const [LayoutName, setLayoutName] = useState(
    localStorage[LocalStorageKeys.LAYOUT_NAME],
  );
  const recalculateSize = useCallback((data: { components: string[] }) => {
    data.components.forEach(i => {
      MessageService.send({
        name: MessageNames.LAYOUT_RESIZE,
        payload: { i },
      });
    });
  }, []);
  useEffect(() => {
    const Subscription = Subscriber.subscribe((message: any) => {
      if (message.name === MessageNames.RESIZE) {
        redoLayOutAnimation();
        setTimeout(() => {
          forceUpdate();
        }, 50);
      } else if (message.name === MessageNames.CHANGE_LAYOUT) {
        redoLayOutAnimation();
        localStorage[LocalStorageKeys.LAYOUT_NAME] = message.payload;
        setLayoutName(message.payload);
        setTimeout(() => {
          recalculateSize({
            components: [
              LayoutContainers.ORDERBOOK,
              LayoutContainers.MARKETTRADE,
              LayoutContainers.ORDERS,
              LayoutContainers.MARKETWATCH,
            ],
          });
        }, 225);
      } else if (message.name === MessageNames.CLOSE_MODAL) {
        setisLoginPopupOpen(false);
      } else if (message.name === MessageNames.OPEN_LOGIN_POPUP) {
        setisLoginPopupOpen(true);
      }
    });
    return () => {
      Subscription.unsubscribe();
    };
  }, []);
  useLayoutEffect(() => {
    setTimeout(() => {
      const elements = document.querySelectorAll('.react-grid-item');
      elements.forEach(item => {
        item.classList.add('animated');
      });
    }, 300);
    return () => {};
  }, []);

  const enable = true;
  // const memoizedLayout = useMemo(() => layout(LayoutName), [LayoutName]);
  const MarketWatch = useRef(
    <div key={LayoutContainers.MARKETWATCH} id='MARKETWATCHContainer'>
      <MiniTitledComponent
        title={<FormattedMessage {...translate.MARKETWATCH} />}
      >
        <MarketWatchGrid
          enabled={enable}
          uniqueId={LayoutContainers.MARKETWATCH}
          subject={CentrifugoChannels.TickerChannel}
        />
      </MiniTitledComponent>
    </div>,
  );

  const MarketTrade = useRef(
    <div key={LayoutContainers.MARKETTRADE} id='MARKETTRADEContainer'>
      <MiniTitledComponent
        title={<FormattedMessage {...translate.MARKETTRADE} />}
      >
        <MarketTradeGrid
          enabled={enable}
          uniqueId={LayoutContainers.MARKETTRADE}
          subject={CentrifugoChannels.MarketTradePrefix}
        />
      </MiniTitledComponent>
    </div>,
  );

  const OrderBook = useRef(
    <div key={LayoutContainers.ORDERBOOK} id='ORDERBOOKContainer'>
      <MiniTitledComponent
        title={<FormattedMessage {...translate.ORDERBOOK} />}
      >
        <OrderBookGrid
          enabled={enable}
          uniqueId={LayoutContainers.ORDERBOOK}
          subject={CentrifugoChannels.OrderBookPrefix}
        />
      </MiniTitledComponent>
    </div>,
  );

  const Tradeheader = useRef(
    <div key={LayoutContainers.TRADEHEADER} id='TRADEHEADERContainer'>
      <TradeHeader subject={CentrifugoChannels.TickerChannel} />
    </div>,
  );

  const Neworder = useRef(
    <div key={LayoutContainers.NEWORDER} id='NEWORDERContainer'>
      <MiniTitledComponent title={<FormattedMessage {...translate.NEWORDER} />}>
        <NewOrder />
      </MiniTitledComponent>
    </div>,
  );

  const orders = useRef(
    <div key={LayoutContainers.ORDERS} id='ORDERSContainer'>
      <Orders />
    </div>,
  );

  const Tradechart = useRef(
    <div key={LayoutContainers.TRADECHART} id='TRADECHARTContainer'>
      <TradeChart enabled={enable} />
    </div>,
  );

  useLayoutEffect(() => {
    //to re-let :)) the interaction with main chart
    const elements = document.querySelectorAll('.react-grid-item');
    elements.forEach(item => {
      item.addEventListener('mouseup', () => {
        const chartWrapper = document.getElementById('mainChartWrapper');
        if (chartWrapper) {
          chartWrapper.style.pointerEvents = 'initial';
        }
      });
    });
    return () => {
      const elements = document.querySelectorAll('.react-grid-item');
      elements.forEach(item => {
        item.removeEventListener('mouseup', () => {});
      });
    };
  }, []);
  return (
    <Wrapper>
      <PopupModal
        isOpen={isLoginPopupOpen}
        onClose={() => {
          setisLoginPopupOpen(false);
        }}
      >
        <LoginPage isPopup />
      </PopupModal>

      <GridLayout
        onResizeStop={handleOnResize}
        onDragStop={handleDrag}
        onDragStart={onDragStart}
        draggableHandle='.dragHandle'
        className='DraGGableGridsContainer'
        layout={layout(LayoutName)}
        cols={48}
        //@ts-ignore
        margin={layoutMargin()}
        rowHeight={LayoutHeight()}
        width={LayoutWidth()}
      >
        {MarketWatch.current}
        {MarketTrade.current}
        {OrderBook.current}
        {Tradeheader.current}
        {Neworder.current}
        {orders.current}
        {Tradechart.current}
      </GridLayout>
    </Wrapper>
  );
};
export default memo(TradeLayout, () => true);
const Wrapper = styled.div`
  .react-resizable-handle {
    ::after {
      border-right: 2px solid var(--textGrey) !important;
      border-bottom: 2px solid var(--textGrey) !important;
    }
  }
  .react-grid-item.react-grid-placeholder {
    background: var(--textBlue) !important;
    border-radius: 3px !important;
    opacity: 0.5 !important;
  }

  .react-grid-item {
    transition: none !important;
    overflow: hidden;
    transition-property: left, top;
    &.animated {
      transition: all 200ms ease !important;
    }
  }
  .ag-header {
    font-size: 11px !important;
    font-weight: 500 !important;
    span {
      font-size: 11px !important;
      font-weight: 500 !important;
    }
  }
  .ag-cell {
    font-size: 11px !important;
    font-weight: 600 !important;
    color: var(--tradeGridTextColor);
    text-shadow: 0 0 1px rgba(0, 0, 0, 0.25);
    span {
      font-size: 11px !important;
      font-weight: 600 !important;
      text-shadow: 0 0 1px rgba(0, 0, 0, 0.25);
    }
  }
  .dragHandle {
    cursor: all-scroll;
  }
  .ag-theme-balham
    .ag-ltr
    .ag-header-cell:not(.ag-numeric-header)
    .ag-header-label-icon {
    margin-left: unset;
  }
`;
