import { TextField } from '@material-ui/core';
import { GridApi, ColDef } from 'ag-grid-community';
import { FilterArrayElement } from 'locales/types';
import React, { CSSProperties, memo, useCallback, useEffect, useState } from 'react';
import { useTheme } from 'styled-components/macro';
import { Subscriber, MessageNames, BroadcastMessage } from 'services/messageService';
import styled from 'styled-components/macro';
import { under } from 'utils/formatters';
import { filterHeight } from 'utils/gridUtilities';

import CountryFilter from './CountryFilter';
import DateFilter from './DateFilter';
import DropDown from './DropDown';
interface Col {
	actualWidth: number;
	colDef: ColDef;
	colId: string;
}
interface Props {
	filters: FilterArrayElement;
	/** The AG Grid API instance, used to read column widths for filter positioning. */
	gridApi: GridApi | null;
	/** Called when any filter field changes. `immediate` skips debounce. */
	onFilterChange: (filters: Record<string, unknown>, immediate?: boolean) => void;
	gridId?: string;
}

/**
 * Horizontal filter bar that sits above an AG Grid and mirrors its column widths.
 * Renders text inputs, dropdowns, date pickers, or country selectors depending
 * on the FilterArrayElement configuration passed in.
 *
 * @example
 * ```tsx
 * <GridFilter
 *   filters={filterConfig}
 *   gridApi={gridApiRef.current}
 *   gridId="usersGrid"
 *   onFilterChange={(params, immediate) => loadData(params)}
 * />
 * ```
 */
const GridFilter = (props: Props) => {
	const theme = useTheme();
	let filterObject: Record<string, unknown> = {};
	const { gridApi, filters, onFilterChange, gridId } = props;
	const add00 = filters.add00ToDate
	let columns: Col[] = [];
	const disabledColumns = filters.disabledCols;
	const hiddenColumns = filters.hiddenCols;
	const dropDownColumns = filters.dropDownCols;
	const dateColumns = filters.dateCols;
	const countryColumns = filters.countryCols;
	const sortableColumns = filters.sortableCols;

	const [GridApi, setGridApi] = useState(gridApi);
	useEffect(() => {
		const Subscription = Subscriber.subscribe((message: BroadcastMessage) => {
			if (
				message.name === MessageNames.GRID_RESIZE &&
				message.payload && 
				typeof message.payload === 'object' &&
				'gridId' in message.payload &&
				message.payload.gridId === gridId &&
				'gridApi' in message.payload
			) {
				setGridApi((message.payload as Record<string, unknown>).gridApi as GridApi);

			}
		});
		return () => {
			Subscription.unsubscribe();
		};
	}, []);
	if (GridApi) {
		// Using internal AG Grid API - no public alternative available
		columns = (GridApi as unknown as { columnController: { displayedCenterColumns: Array<{ colDef: ColDef; actualWidth: number; colId: string }> } }).columnController.displayedCenterColumns.filter((item: { colDef: ColDef; actualWidth: number; colId: string }) => {
			return item.actualWidth > 30;
		});
	}


	const FilterInput = useCallback((props: { data: Col }) => {
		const { data } = props;

		if (sortableColumns && sortableColumns.includes(data.colId)) {
		}





		if (dropDownColumns) {
			for (let i = 0; i < dropDownColumns.length; i++) {
				if (dropDownColumns[i].id === data.colId) {
					return (
						<DropDown
							title={data.colDef.headerName ?? ''}
							options={dropDownColumns[i].options}
							onSelect={value => {
								let undered =
									dropDownColumns[i].substituteId ?? under(data.colId);
								if (value === '' && filterObject[undered]) {
									delete filterObject[undered];
								}
								else {
									filterObject[undered] = value;
								}
								onFilterChange(filterObject, true);
							}}
						/>
					);
				}
			}
		}
		if (dateColumns && dateColumns.includes(data.colId)) {
			return (
				<DateFilter
					title={data.colDef.headerName ?? ''}
					onDateSelect={e => {
						let undered = under(data.colId);

						if (e === '' && filterObject[undered]) {
							delete filterObject[undered];
						} else {
							filterObject[undered] = e + (add00 ? ' 00:00:00' : '');
						}
						onFilterChange(filterObject, true);
					}}
				/>
			);
		}
		if (countryColumns && countryColumns.includes(data.colId)) {
			return (
				<CountryFilter
					onClear={() => {
						delete filterObject['country_id'];
						onFilterChange(filterObject, true);
					}}
					onChange={e => {
						filterObject['country_id'] = e;
						onFilterChange(filterObject, true);
					}}
				/>
			);
		}
		return (
			<>
				{
					<TextField
						disabled={disabledColumns && disabledColumns.includes(data.colId)}
						variant="outlined"
						onChange={e => {
							let undered = under(data.colId);
							if (e.target.value === '' && filterObject[undered]) {
								delete filterObject[undered];
							} else {
								filterObject[undered] = e.target.value;
							}
							onFilterChange(filterObject);
						}}
						label={data.colDef.headerName}
					/>
				}
			</>
		);
	}, []);
	useEffect(() => {
		return () => {
			filterObject = {};
		};
	}, []);
	return (
		<>
			{GridApi && (
				<Wrapper id={'filter' + gridId}>
					{columns.map((item, index) => {
						let st: CSSProperties =
							(!hiddenColumns || !hiddenColumns.includes(item.colId)) === true
								? {}
								: {
									opacity: 0,
									pointerEvents: 'none',
								};
						return (
							<div
								key={item.colId}
								style={{
									width:
										item.actualWidth -
										(item.colId === 'loading' || item.colId === 'actions'
											? item.actualWidth
											: 20) +
										'px',
									height: filterHeight + 'px',

									position: 'relative',
									background:
										item.colId === 'loading' ? 'transparent' : theme.background,
									borderRadius: '3px',
									...st
								}}
								id={item.colId + gridId ?? ''}
							>
								{FilterInput({ data: item })}
								{/* <FilterInput data={item} /> */}
							</div>
						);
					})}
				</Wrapper>
			)}
		</>
	);
};
export default memo(GridFilter);
const Wrapper = styled.div`
  display: flex;
  min-height: ${filterHeight}px;
  min-height: 40px;
  padding-top: 10px;
  margin-top: -20px;
  padding-bottom: 10px;
  padding-left: 10px;
  position: absolute;
  z-index: 1;
  background: ${p => p.theme.filterBackground};

  width: calc(100% - 24px);
  margin-left: 12px;
  border-top-left-radius: 7px;
  border-top-right-radius: 7px;
  box-shadow: 12px 0px ${p => p.theme.white};
  gap: 20px;
  min-width: 1130px;
  input.MuiOutlinedInput-input {
    padding: 0 4px !important;
    box-shadow: none !important;
    border: none !important;
  }
  .MuiOutlinedInput-root,
  .MuiTextField-root {
    height: 100% !important;
    width: 100%;
  }
  .MuiInputLabel-outlined {
    color: ${p => p.theme.textGrey} !important;
    font-size: 12px;
    font-weight: 500;
    line-height: 0px;
    margin-top: -4px;
    &.Mui-focused {
      margin-top: 6px;
    }
  }
  legend {
    zoom: 0.9;
  }

  .MuiInputLabel-outlined.MuiInputLabel-shrink {
    margin-top: 6px;
  }
  input {
    background-color: transparent !important;
    border: none;
  }
  .select {
    font-size: 13px;
    font-weight: 500 !important;
    color: ${p => p.theme.textGrey} !important;
    width: 100%;
  }
  .MuiInputLabel-outlined {
    z-index: 0 !important;
  }
`;
