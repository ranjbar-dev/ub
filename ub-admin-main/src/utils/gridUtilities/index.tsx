import { GridApi, RowNode } from 'ag-grid-community';
import { rowHeight } from 'app/constants';
export const headerHider = (data: {
  gridApi: GridApi;
  gridId?: string;
  hideBottomBorderFromTitledComponent?: boolean;
  wrapperId?: string;
}) => {
  let lastRendered = data.gridApi.getLastDisplayedRow();
  const header = (!data.gridId
    ? document.getElementsByClassName('ag-header')
    : document.getElementById(data.gridId)?.getElementsByClassName('ag-header'));
  const withheader = (!data.gridId
    ? document.getElementsByClassName('ag-header')
    : document
        .getElementById(data.gridId)
        ?.getElementsByClassName('with-ag-header'));
  let titledComponent = document.getElementById(
    data.wrapperId ?? 'withHidableBorder',
  );
  if (lastRendered === -1) {
    if (header && header[0]) {
      (header[0] as HTMLElement).style.opacity = '0';
    }
    if (withheader && withheader[0]) {
      (withheader[0] as HTMLElement).style.opacity = '0';
    }
    if (titledComponent && data.hideBottomBorderFromTitledComponent === true) {
      (titledComponent.getElementsByTagName('hr')[0] as HTMLElement).style.opacity = '1';
    }
  } else {
    if (header && header[0]) {
      (header[0] as HTMLElement).style.opacity = '1';
    }
    if (withheader && withheader[0]) {
      (withheader[0] as HTMLElement).style.opacity = '1';
    }
    if (titledComponent && data.hideBottomBorderFromTitledComponent === true) {
      (titledComponent.getElementsByTagName('hr')[0] as HTMLElement).style.opacity = '0';
    } else if (titledComponent) {
      (titledComponent.getElementsByTagName('hr')[0] as HTMLElement).style.opacity = '1';
    }
  }
};
export const whiteShaddowHider = (data: { gridApi: GridApi; isMini?: boolean }) => {
  if (!data.isMini) {
    if (
      data.gridApi.getLastDisplayedRow() >=
      data.gridApi.getModel().getRowCount() - 1
    ) {
      const whiteShaddow = document.getElementById('whiteShaddow');
      if (whiteShaddow) {
        whiteShaddow.style.opacity = '0';
      }
    } else {
      const whiteShaddow = document.getElementById('whiteShaddow');
      if (whiteShaddow) {
        whiteShaddow.style.opacity = '1';
      }
    }
  }
};

export const RowWithShaddow = (params: { data: { status?: string; isDetailsOpen?: boolean } }) => {
  let filter = 'unset';
  let pointerEvents = '';
  if (
    params.data.status &&
    (params.data.status === 'canceled' || params.data.status === 'cancel')
  ) {
    filter = 'opacity(0.5) grayscale(1)';
    pointerEvents = 'none';
  }
  let boxShadow = ' 0 0 0px 0px rgba(0,0,0,0.3)';
  let zIndex = 1;
  if (params.data.isDetailsOpen === true) {
    boxShadow = '0px 0px 6px 2px rgba(0, 0, 0, 0.08) ';
    zIndex = 2;
  }
  return { boxShadow, zIndex, filter, pointerEvents };
};
export const DetailsStyle = `

.detailsContainer {
  width: 800px;
  height: calc(100% - 40px);
  position: absolute;
  left: calc(-43vw + 45px);
  bottom: 0px;
  overflow: hidden;
}
.rotateAnimate200{
  transition: transform 0.2s;
  path{
    fill:var(--textBlue) !important;
  }
}
.rotated{
  transform:rotate(180deg);
  path{
    fill:var(--textGrey) !important;
  }
}
`;
export const ToggleDetail = (data: {
  gridApi: GridApi;
  params: { data: { id: string | number; isDetailsOpen?: boolean } };
  isOpen: boolean;
  height?: number;
  isMini?: boolean;
}) => {
  let el = document.getElementById('expandIcon' + data.params.data.id);
  let dtButton = document.getElementById('dtButton' + data.params.data.id);
  let txID = document.getElementById('txId' + data.params.data.id);
  if (data.isOpen === true) {
    if (el) {
      el.classList.remove('rotated');
    }
    if (dtButton) {
      dtButton.classList.remove('grey');
    }
    if (txID) {
      txID.style.opacity = '1';
    }
    data.gridApi.forEachNode((node: RowNode, index: number) => {
      if (node.data.id === data.params.data.id) {
        node.data.isDetailsOpen = false;
        data.gridApi.applyTransaction({
          update: [node.data],
        });
        node.setRowHeight(data.isMini ? 25 : 40);
        data.gridApi.onRowHeightChanged();
      }
    });
    // data.gridApi.resetRowHeights();
    return;
  }
  if (el) {
    el.classList.add('rotated');
  }
  if (dtButton) {
    dtButton.classList.add('grey');
  }
  if (txID) {
    txID.style.opacity = '0';
  }
  data.gridApi.forEachNode((node: RowNode, index: number) => {
    if (node.data.id === data.params.data.id) {
      node.data.isDetailsOpen = true;
      data.gridApi.applyTransaction({
        update: [node.data],
      });
      node.setRowHeight(data.height ? data.height : 120);
      data.gridApi.onRowHeightChanged();
    }
  });
  // data.gridApi.resetRowHeights();
};
export const resizeToTarget = (data: {
  elementId: string;
  resizeToElementWithId: string;
  additional?: number;
}) => {
  let target = document.getElementById(data.resizeToElementWithId);
  if (target) {
    let element = document.getElementById(data.elementId);
    if (element) {
      element.style.height =
        target.clientHeight -
        40 -
        (data.additional ? data.additional : 0) +
        'px';
    }
  }
};
export const filterHeight = 30;
export const getPageSize = (data?: {
  wrapperHeight?: number;
  gridHasFilter?: boolean;
  additional?: number;
}) => {
  let additional = data && data.additional ? data.additional : 0;
  if (data && data.gridHasFilter === true) {
    additional += filterHeight;
  }
  const containerHeight =
    data && data.wrapperHeight
      ? data.wrapperHeight - additional
      : window.innerHeight - 230 - additional;
  const rowNumber = (containerHeight / rowHeight).toFixed(0);
  return rowNumber;
};
export const randomColor = () => {
  let letters = '0123456789ABCDEF';
  let color = '#';
  for (let i = 0; i < 6; i++) {
    color += letters[Math.floor(Math.random() * 16)];
  }
  return color;
};
