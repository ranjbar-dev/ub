import {
  ColDef,
  GridReadyEvent,
  RowClickedEvent,
  GridApi,
} from 'ag-grid-community';
import { AgGridReact } from 'ag-grid-react';
import { CellRenderer } from 'app/components/renderer';
import { GridWrapper } from 'app/components/wrappers/GridWrapper';
import { rowHeight } from 'app/constants';
import { translations } from 'locales/i18n';
import React, { memo, useCallback, useRef, useMemo, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { cellColorAndNameFormatter } from 'utils/stylers';

import { ProfileImageData } from '../types';

interface Props {
  images: ProfileImageData[];
  selectedImage: ProfileImageData | null;
  onImageSelect: Function;
}

function ImagesGrid(props: Props) {
  const gridApi = useRef<GridApi | null>(null);
  const { images, onImageSelect, selectedImage } = props;
  const { t } = useTranslation();
  const staticRows: ColDef[] = [
    {
      headerName: t(translations.CommonTitles.FileName()),
      field: 'originalFileName',
      suppressMenu: true,
      sortable: true,
    },
    {
      headerName: t(translations.CommonTitles.UploadeDate()),
      field: 'createdAt',
      suppressMenu: true,
      sortable: true,
    },
    {
      headerName: t(translations.CommonTitles.LastResponseDate()),
      field: 'updatedAt',
      suppressMenu: true,
      sortable: true,
    },
    {
      headerName: t(translations.CommonTitles.Status()),
      field: 'confirmationStatus',
      suppressMenu: true,
      sortable: true,
      ...cellColorAndNameFormatter('confirmationStatus'),
    },
    {
      headerName: t(translations.CommonTitles.ReasonMessage()),
      field: 'rejectionReason',
      suppressMenu: true,
      sortable: true,
    },
    {
      headerName: t(translations.CommonTitles.Side()),
      field: 'isBack',
      suppressMenu: true,
      sortable: true,
      cellRenderer: params =>
        CellRenderer(
          <span>{params.data.isBack === true ? 'Back' : 'Front'}</span>,
        ),
    },
  ];
  //select first row(blue left border)
  const selectFirst = useCallback(() => {
    gridApi.current?.forEachNode((node: { data: ProfileImageData & { isSelected?: boolean } }, index: number) => {
      if (index === 0) {
        let data = { ...node.data, isSelected: true };
        node.data = data;
        gridApi.current?.redrawRows();
        return;
      }
    });
  }, []);
  //  const selectFirstRow = () => {
  //    let rows = document.querySelectorAll('.ag-row-first');
  //    console.log(rows);
  //  };
  const gridRendered = useCallback((e: GridReadyEvent) => {
    setTimeout(() => {
      gridApi.current = e.api;
      gridApi.current.sizeColumnsToFit();
      //  selectFirstRow();
      selectFirst();
    }, 0);
  }, []);
  //  if (gridApi.current) {
  //    selectFirstRow();
  //  }
  //check if first row is not selected
  if (gridApi.current) {
    let selected = false;
    gridApi.current.forEachNode((node: { data: ProfileImageData & { isSelected?: boolean } }, index: number) => {
      if (node.data.isSelected === true) {
        selected = true;
      }
    });
    if (selected === false) {
      selectFirst();
    }
  }
  const gridConfig: { columnDefs: ColDef[] } = {
    columnDefs: [...staticRows],
  };
  const onRowClicked = useCallback((e: RowClickedEvent) => {
    onImageSelect(e.data);
    setTimeout(() => {
      gridApi.current?.forEachNode((node: { data: ProfileImageData & { isSelected?: boolean } }, index: number) => {
        if (node.data.id === e.data.id) {
          let data = { ...node.data, isSelected: true };
          node.data = data;
        } else {
          let data = { ...node.data, isSelected: false };
          node.data = data;
        }
      });
      gridApi.current?.redrawRows();
    }, 0);
  }, []);

  return (
    <>
      <GridWrapper
        className={`ag-theme-balham clickableRows imageGridView`}
        style={{ height: 255 + 'px', width: '100%' }}
      >
        {selectedImage ? (
          <AgGridReact
            onGridReady={gridRendered}
            animateRows={true}
            headerHeight={29}
            rowHeight={rowHeight}
            onRowClicked={onRowClicked}
            getRowClass={params => {
              return params.node.data.isSelected === true ? 'selectedRow' : '';
            }}
            columnDefs={gridConfig.columnDefs}
            rowData={images}
            immutableData={true}
            enableCellChangeFlash={true}
            getRowNodeId={(data: ProfileImageData) => {
              return data.id.toString();
            }}
            //  onRowClicked={onRowClocked}
            //  overlayNoRowsTemplate={ReactDOMServer.renderToString(
            //    <div>{images.length === 0 && 'no rows'}</div>,
            //  )}
          ></AgGridReact>
        ) : (
          ''
        )}
      </GridWrapper>
    </>
  );
}

export default memo(ImagesGrid);
