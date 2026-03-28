import { ActionType } from 'typesafe-actions';
import * as actions from './actions';
import { ApplicationRootState } from 'types';

/* --- STATE --- */
interface TradeChartState {
  readonly default: any;
  readonly chartConfig: any;
}

/* --- ACTIONS --- */
type TradeChartActions = ActionType<typeof actions>;

/* --- EXPORTS --- */
type RootState = ApplicationRootState;
type ContainerState = TradeChartState;
type ContainerActions = TradeChartActions;
export interface ChartContainerProps {
  symbol: ChartingLibraryWidgetOptions['symbol'];
  interval: ChartingLibraryWidgetOptions['interval'];
  // BEWARE: no trailing slash is expected in feed URL
  datafeedUrl?: string;
  libraryPath: ChartingLibraryWidgetOptions['library_path'];
  chartsStorageUrl: ChartingLibraryWidgetOptions['charts_storage_url'];
  chartsStorageApiVersion: ChartingLibraryWidgetOptions['charts_storage_api_version'];
  clientId: ChartingLibraryWidgetOptions['client_id'];
  userId: ChartingLibraryWidgetOptions['user_id'];
  fullscreen: ChartingLibraryWidgetOptions['fullscreen'];
  autosize: ChartingLibraryWidgetOptions['autosize'];
  studiesOverrides: ChartingLibraryWidgetOptions['studies_overrides'];
  containerId?: ChartingLibraryWidgetOptions['container_id'];
}
export interface ChartSummary {
  baseVolume?: number;
  closePrice?: number;
  highPrice?: number;
  lowPrice?: number;
  ohlcCloseTime?: string;
  ohlcStartTime?: string;
  openPrice?: number;
  quoteVolume?: number;
  takerBuyBaseVolume?: number;
  takerBuyQuoteVolume?: number;
}
export { RootState, ContainerState, ContainerActions };
