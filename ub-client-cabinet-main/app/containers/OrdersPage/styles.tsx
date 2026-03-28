import { css } from 'styles/styled-components';

export const miniGridStyles = css`
  &.miniGrid {
    height: 50px;
    .ag-row {
      border-bottom: 2px solid var(--white) !important;
    }
    .ag-row:not(:hover) {
      background-color: var(--ordersMiniGridRowBackground) !important;
    }
    .ag-cell {
      line-height: 20px !important;
      font-size: 11px;
    }
    .statusBadge {
      margin-top: 0px !important;
      &.mini {
        margin-left: -15px;
      }
    }
    .ag-row .detailButton {
      top: -6px;
    }
    .value.small {
      line-height: 34px !important;
    }
    .cancelButtons {
      background: var(--miniCancelButtonBackground) !important;
      border-radius: 5px !important;
      span {
        font-size: 10px !important;
        color: var(--miniCancelButtonTextColor) !important;
      }
      border-radius: 5px !important;
      position: absolute;
      top: 2px !important;
      max-height: 16px;
      padding: 0 !important;
      line-height: 0px;
      min-height: 16px;
      left: 5px;
      min-width: 44px;
    }

  }
`;
