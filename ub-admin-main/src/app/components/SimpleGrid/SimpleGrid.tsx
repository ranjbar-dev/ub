import { ActionCreatorWithPayload } from '@reduxjs/toolkit';
import {
	GridReadyEvent,
	ColDef,
	RowClickedEvent,
	CellClickedEvent,
	GridApi,
	RowNode,
} from 'ag-grid-community';
import { AgGridReact } from 'ag-grid-react';
import { GridLoading } from 'app/components/grid_loading/gridLoading';
import PaginationComponent from 'app/components/PaginationComponent/PaginationComponent';
import { GridWrapper } from 'app/components/wrappers/GridWrapper';
import { rowHeight } from 'app/constants';
import { FilterArrayElement, GridTopTab } from 'locales/types';
import React, {
	memo,
	useEffect,
	useState,
	useCallback,
	useMemo,
	useRef,
} from 'react';
import ReactDOMServer from 'react-dom/server';
import { useDispatch } from 'react-redux';
import {
	Subscriber,
	MessageNames,
	GridNames,
	MessageService,
	BroadcastMessage,
} from 'services/messageService';
import styled from 'styled-components/macro';
import {
	getPageSize,
	resizeToTarget,
	filterHeight,
	randomColor,
} from 'utils/gridUtilities';

import GridFilter from '../GridFilter/GridFilter';
import GridTabs from '../GridTabs/GridTabs';

interface Props {
	messageName?: MessageNames;
	/** Redux action creator dispatched to load/paginate grid data. */
	// eslint-disable-next-line @typescript-eslint/no-explicit-any
	initialAction: ActionCreatorWithPayload<any>;
	arrayFieldName: string;
	pagination?: boolean;
	staticRows: ColDef[];
	immutableId: string;
	gridName?: GridNames;
	onPageSizeSet?: (page_size: string) => void
	onRowClick?: (event: RowClickedEvent) => void;
	onDataReceived?: (data: Record<string, unknown>) => void;
	onCellClick?: (event: CellClickedEvent) => void;
	filters?: FilterArrayElement;
	additionalInitialParams: Record<string, unknown>;
	additionalColumns?: ColDef[];
	containerId?: string;
	flashCellUpdate?: boolean;
	userId?: number;
	topTabs?: GridTopTab[];
	additionalHeight?: number;
	pullUpPaginationBy?: number;
	heightWhenHasNoContainerId?: number;
	/**
	 * When provided, data comes from Redux (via useSelector) instead of
	 * the MessageService subscription. Pass `null` while loading, the page
	 * response object once data arrives.
	 */
	// eslint-disable-next-line @typescript-eslint/no-explicit-any
	externalData?: Record<string, unknown> | null;
}

/**
 * Paginated AG Grid wrapper for admin panel data tables.
 * Subscribes to a message bus for data updates, handles pagination,
 * column filters, top-tab switching, and row-level loading states.
 *
 * @example
 * ```tsx
 * <SimpleGrid
 *   messageName={MessageNames.USERS_DATA}
 *   initialAction={UserAccountsActions.getUsers}
 *   arrayFieldName="users"
 *   staticRows={columnDefs}
 *   immutableId="id"
 *   additionalInitialParams={{}}
 * />
 * ```
 */
export const SimpleGrid = memo((props: Props) => {
	const {
		messageName,
		initialAction,
		pagination = true,
		arrayFieldName,
		additionalColumns,
		staticRows,
		additionalInitialParams,
		immutableId,
		containerId,
		topTabs,
		onDataReceived,
		gridName,
		userId,
		pullUpPaginationBy,
		flashCellUpdate,
		additionalHeight,
		filters,
		onRowClick,
		onCellClick,
		onPageSizeSet,
		heightWhenHasNoContainerId,
		externalData,
	} = props;
	// True when the parent has opted into Redux-based data (externalData prop was passed).
	const useExternalData = externalData !== undefined;
	const cols = additionalColumns ? additionalColumns : [];
	const dispatch = useDispatch();
	const Initialized = useRef(false);
	const Timeout = useRef<ReturnType<typeof setTimeout> | undefined>(undefined);
	const gridApi = useRef<GridApi | null>(null);
	const topFilterObj = useRef<Record<string, unknown>>({});
	const topTabObj = useRef<Record<string, unknown>>({});
	const lastPage = useRef(1);
	const lastGridWidth = useRef('0px');
	const containerRef = useRef<HTMLDivElement>(null);
	const gridId = useRef(randomColor());
	const resizeTimer = useRef<ReturnType<typeof setInterval> | undefined>(undefined);
	//   ref: containerRef,
	//   onResize: ({ width, height, entry, unobserve, observe }) => {
	//     if (gridApi.current) {
	//       // const grridApi:GridApi=gridApi.current;
	//       // grridApi.co
	//       // grridApi.
	//       console.debug(width);

	//     }
	//   },
	// });

	const pageSize = useMemo(() => {
		let additional = -40;
		let addHeight = additionalHeight ?? 0;
		if (containerId === 'UserDetailsWindow') {
			additional += 160 + addHeight;
		}
		if (topTabs) {
			additional += 30 + addHeight;
		}
		let pagesize = Number(
			getPageSize({
				gridHasFilter: filters != null ? true : false,
				additional: additional,
			}),
		) -
			(pullUpPaginationBy ? 1 : 0) +
			''
		onPageSizeSet && onPageSizeSet(pagesize)
		return pagesize
	}, []);
	const [PageData, setPageData] = useState<Record<string, unknown> | undefined>();

	const [IsLoading, setIsLoading] = useState(true);
	useEffect(() => {
		dispatch(
			initialAction({
				page: 1,
				page_size: pageSize,
				...additionalInitialParams,
			}),
		);
		const Subscription = Subscriber.subscribe((message: BroadcastMessage) => {
			if (
				message.name === messageName &&
				!useExternalData &&
				(!userId || userId === message.userId)
			) {
				Initialized.current = true;
				setIsLoading(false);

				setPageData(message.payload as Record<string, unknown>);
				gridApi.current!.setRowData((message.payload as Record<string, unknown>)[arrayFieldName] as unknown[]);
				onDataReceived && onDataReceived(message.payload as Record<string, unknown>)
			}
			if (
				gridName &&
				message.name === MessageNames.SET_ROW_LOADING &&
				(message.payload as Record<string, unknown>).gridName === gridName &&
				(message.payload as Record<string, unknown>).userId === userId
			) {
				gridApi.current!.forEachNode((node: RowNode) => {
					if (node.data.id === (message.payload as Record<string, unknown>).rowId) {
						node.data.IsLoading = (message.payload as Record<string, unknown>).value;
						gridApi.current!.applyTransaction({
							update: [node.data],
						});
						gridApi.current!.redrawRows({ rowNodes: [node] });
						return;
						//  node.setRowHeight(data.isMini ? 25 : 40);
						//  data.gridApi.onRowHeightChanged();
					}
				});
			}
			if (
				message.name === MessageNames.UPDATE_GRID_ROW &&
				gridName &&
				gridName === message.gridName
			) {
				gridApi.current!.forEachNode((node: RowNode) => {
					if (node.data.id === message.rowId) {
						const payloadObj = message.payload as Record<string, unknown>;
						for (const key in payloadObj) {
							if (Object.prototype.hasOwnProperty.call(payloadObj, key)) {
								node.data[key] = payloadObj[key];
							}
						}
						gridApi.current!.applyTransaction({
							update: [node.data],
						});
						gridApi.current!.redrawRows({ rowNodes: [node] });
						return;
					}
				});
			}
			if (
				message.name === MessageNames.REFRESH_GRID &&
				gridName &&
				message.gridName === gridName
			) {

				const callObj = {
					page: lastPage.current,
					page_size: pageSize,
					...additionalInitialParams,
					...topTabObj.current,
					...topFilterObj.current,
				}

				dispatch(
					initialAction(callObj),
				);
			}
			if (
				message.name === MessageNames.APPLY_PARAMS_TO_GRID &&
				gridName &&
				message.gridName === gridName
			) {
				handleTopTabChange(message.payload as Record<string, unknown>);
			}
		});
		return () => {
			Initialized.current = false;
			Timeout.current = undefined;
			clearInterval(resizeTimer.current);
			Subscription.unsubscribe();
		};
	}, []);

	// When the parent provides data via Redux (externalData prop), update the grid.
	useEffect(() => {
		if (externalData !== undefined && externalData !== null && gridApi.current) {
			Initialized.current = true;
			setIsLoading(false);
			setPageData(externalData);
			gridApi.current.setRowData((externalData[arrayFieldName] as unknown[]) || []);
			onDataReceived && onDataReceived(externalData);
		}
	}, [externalData]);

	const gridConfig: { columnDefs: ColDef[] } = {
		columnDefs: [...staticRows, ...cols],
	};
	const gridRendered = useCallback((e: GridReadyEvent) => {
		gridApi.current = e.api;
		if (containerId) {
			let toReduce = 0;
			if (containerId == 'UserDetailsWindow') {
				toReduce = 50;
			}
			if (topTabs) {
				toReduce += 50;
			}
			resizeToTarget({
				elementId: containerId + 'SimpleGridWrapper',
				resizeToElementWithId: containerId,
				additional:
					(additionalHeight ?? 0) +
					(filters != null ? filterHeight + toReduce : 0),
			});
		}
		gridApi.current!.sizeColumnsToFit();
		MessageService.send({
			name: MessageNames.GRID_RESIZE,
			payload: {
				gridApi: gridApi.current,
				gridId: gridId.current,
			},
		});
		// @ts-expect-error — containerRef.current may be null
		lastGridWidth.current = containerRef.current?.offsetWidth;
		// did this shit because some times new window is a jerk while importing app css
		resizeTimer.current = setInterval(() => {
			// @ts-expect-error — containerRef.current may be null
			if (lastGridWidth.current != containerRef.current.offsetWidth) {
				gridApi.current!.sizeColumnsToFit();
				MessageService.send({
					name: MessageNames.GRID_RESIZE,
					payload: {
						gridApi: gridApi.current,
						gridId: gridId.current,
					},
				});
				// @ts-expect-error — containerRef.current may be null
				lastGridWidth.current = containerRef.current?.offsetWidth;
			}
		}, 500);
	}, []);
	const handlePageChange = useCallback((pageNumber: number) => {
		setIsLoading(true);
		lastPage.current = pageNumber;

		const callObj = {
			page: pageNumber,
			page_size: pageSize,
			...additionalInitialParams,
			...topTabObj.current,
			...topFilterObj.current,
		}

		dispatch(
			initialAction(callObj),
		);
	}, []);
	const handleFilterChange = (e: Record<string, unknown>, immidiate?: boolean) => {


		setIsLoading(true);

		clearTimeout(Timeout.current);

		topFilterObj.current = e;

		const callObj = {
			page: 1,
			page_size: pageSize,
			...additionalInitialParams,
			...topTabObj.current,
			...e,
		}


		Timeout.current = setTimeout(
			() => {
				dispatch(
					initialAction(callObj),
				);
			},
			immidiate === true ? 0 : 500,
		);
	};
	const handleTopTabChange = useCallback((callObject: Record<string, unknown>) => {
		setIsLoading(true);
		for (const key in callObject) {
			if (Object.prototype.hasOwnProperty.call(callObject, key)) {
				if (callObject[key] === null) {
					delete callObject[key];
					if (additionalInitialParams[key]) {
						delete additionalInitialParams[key];
					}
				}
			}
		}
		topTabObj.current = callObject;

		const callObj = {
			page: 1,
			page_size: pageSize,
			...additionalInitialParams,
			...topTabObj.current,
		}

		dispatch(
			initialAction(callObj),
		);
	}, []);
	const filterContent = useMemo(
		() =>
			filters ? (
				<GridFilter
					onFilterChange={handleFilterChange}
					filters={filters}
					gridId={gridId.current}
					gridApi={gridApi.current}
				/>
			) : null,
		[gridApi.current],
	);
	return (
		<>
			{IsLoading === true && <GridLoading />}
			<Wrapper className="simpleGrid">
				{topTabs && <GridTabs onChange={handleTopTabChange} tabs={topTabs} />}
				<GridWrapper
					ref={containerRef}
					id={(containerId ? containerId : 'noContainer') + 'SimpleGridWrapper'}
					style={
						heightWhenHasNoContainerId
							? {
								height: heightWhenHasNoContainerId + 'px',
								...(additionalHeight ? { marginBottom: '28px' } : {}),
							}
							: { ...(additionalHeight ? { marginBottom: '28px' } : {}) }
					}
					className={`ag-theme-balham  ${onRowClick != null ? 'clickableRows' : ''
						} ${filters != null ? 'withFilter' : ''}`}
				>
					<div
						style={{
							opacity: Initialized.current === true ? '1' : 0,
							minHeight: filters ? filterHeight + 'px' : '0px',
						}}
					>
						{filterContent}
					</div>
					{useMemo(
						() => (
							<AgGridReact
								onGridReady={gridRendered}
								animateRows={true}
								enableCellChangeFlash={flashCellUpdate ?? false}
								headerHeight={filters ? 0 : 32}
								enableCellTextSelection
								rowHeight={rowHeight}
								columnDefs={gridConfig.columnDefs}
								tooltipShowDelay={0}
								defaultColDef={{ suppressMenu: true, sortable: true, tooltipComponent: 'customTooltip' }}
								rowData={
									PageData && PageData[arrayFieldName]
										? (PageData[arrayFieldName] as unknown[])
										: []
								}

								immutableData={true}
								getRowNodeId={(data: Record<string, unknown>) => {
									return String(data[immutableId] ?? '');
								}}
								onRowClicked={onRowClick ? onRowClick : e => { }}
								onCellClicked={onCellClick ? onCellClick : e => { }}
								overlayNoRowsTemplate={ReactDOMServer.renderToString(
									<div>{IsLoading === false && 'no data'}</div>,
									//  <AnimatedNoRows
									//    isMini={props.isMini}
									//    icon={<NoOpenOrdersIcon />}
									//    texts={[
									//      <span className="black">
									//        {Translate({ message: 'noOrder' })}
									//      </span>,
									//    ]}
									//  />,
								)}
							></AgGridReact>
						),

						[gridApi.current],
					)}
				</GridWrapper>

				{PageData &&
					PageData.count && pagination &&
					Number(((PageData.count as number) / Number(pageSize)).toFixed(0)) > 1 ? (
					<PaginationComponent
						style={{
							marginTop: pullUpPaginationBy
								? '-' + pullUpPaginationBy + 'px'
								: containerId === 'UserDetailsWindow'
									? '-15px'
									: 0,
						}}
						size={
							Number(((PageData.count as number) / Number(pageSize)).toFixed(0)) +
							((PageData.count as number) % Number(pageSize) > 0 ? 1 : 0)
						}
						onPageChange={handlePageChange}
					/>
				) : (
					''
				)}
			</Wrapper>
		</>
	);
});

const Wrapper = styled.div`
  box-shadow: 0px 6px 7px rgba(0, 0, 0, 0.18);
  border-radius: 6px;
  border-top-left-radius: 0px;
  border-top-right-radius: 0px;
  .bold{
	  font-weight:600;
  }
  .MuiPaginationItem-ellipsis {
    position: relative !important;
  }
`;
