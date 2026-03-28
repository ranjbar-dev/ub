import { GridApi } from 'ag-grid-community';
export const headerHider = (data: {
  gridApi: any;
  gridId?: string;
  hideBottomBorderFromTitledComponent?: boolean;
  wrapperId?: string;
}) => {
  const lastRendered = data.gridApi.getLastDisplayedRow();
  const header: any = !data.gridId
    ? document.getElementsByClassName('ag-header')
    : document.getElementById(data.gridId)?.getElementsByClassName('ag-header');
  const withheader: any = !data.gridId
    ? document.getElementsByClassName('ag-header')
    : document
      .getElementById(data.gridId)
      ?.getElementsByClassName('with-ag-header');
  const titledComponent = document.getElementById(
    data.wrapperId ?? 'withHidableBorder',
  );
  if (lastRendered === -1) {
    if (header && header[0]) {
      header[0].style.opacity = 0;
    }
    if (withheader && withheader[0]) {
      withheader[0].style.opacity = 0;
    }
    if (titledComponent && data.hideBottomBorderFromTitledComponent === true) {
      titledComponent.getElementsByTagName('hr')[0].style.opacity = '1';
    }
  } else {
    if (header && header[0]) {
      header[0].style.opacity = 1;
    }
    if (withheader && withheader[0]) {
      withheader[0].style.opacity = 1;
    }
    if (titledComponent && data.hideBottomBorderFromTitledComponent === true) {
      titledComponent.getElementsByTagName('hr')[0].style.opacity = '0';
    } else if (titledComponent) {
      titledComponent.getElementsByTagName('hr')[0].style.opacity = '1';
    }
  }
};
export const whiteShaddowHider = (data: { gridApi: any; isMini?: boolean }) => {
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

export const RowWithShaddow = (params: any) => {
  let filter = 'unset';
  let pointerEvents = '';
  if (
    params.data.status &&
    (params.data.status == 'canceled' || params.data.status == 'cancel')
  ) {
    filter = 'opacity(0.5) grayscale(1)';
    pointerEvents = 'none';
  }
  let boxShadow = ' 0 0 0px 0px var(--expandedShaddowColor)';
  let zIndex = 1;
  if (params.data.isDetailsOpen === true) {
    boxShadow = '0px 0px 6px 2px var(--expandedShaddowColor) ';
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
  params: any;
  isOpen: boolean;
  height?: number;
  isMini?: boolean;
}) => {
  const el = document.getElementById('expandIcon' + data.params.data.id);
  const dtButton = document.getElementById('dtButton' + data.params.data.id);
  const txID = document.getElementById('txId' + data.params.data.id);
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
    data.gridApi.forEachNode((node, index) => {
      if (node.data.id === data.params.data.id) {
        const { data: tmpData } = node;
        tmpData.isDetailsOpen = false;
        data.gridApi.applyTransaction({
          update: [tmpData],
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
  data.gridApi.forEachNode((node, index) => {
    if (node.data.id === data.params.data.id) {
      const { data: tmpData } = node;
      tmpData.isDetailsOpen = true;
      data.gridApi.applyTransaction({
        update: [tmpData],
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
  const target = document.getElementById(data.resizeToElementWithId);
  if (target) {
    const element = document.getElementById(data.elementId);
    if (element) {
      element.style.height =
        target.clientHeight -
        40 -
        (data.additional ? data.additional : 0) +
        'px';
    }
  }
};
