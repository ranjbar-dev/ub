import { rowHeight } from 'app/constants';
import styled from 'styled-components/macro';

export const GridWrapper = styled.div`
  .ag-header-cell::after {
    display: none !important;
  }
  .ag-header-cell-text {
    color: ${p => p.theme.textGrey} !important;
    font-size: 13px;
  }

  .ag-header {
    /*border-top-left-radius: 7px;
    border-top-right-radius: 7px;*/
    background-color: ${p => p.theme.greyBackground} !important;
    border-bottom: 1px solid ${p => p.theme.greyBackground} !important;
  }
  .ag-cell-focus {
    border: 1px solid transparent !important;
  }
  .ag-cell {
    /* display: flex !important;
		align-items: center !important; */
    /*line-height: 33px !important;*/

    font-size: 12px;
    color: ${p => p.theme.blackText};
    font-weight: 500;
    text-shadow: 0 0 1px #0000005c;
    line-height: ${rowHeight - 5}px;
    border-left: 1px solid #eeeeee !important;
    padding-left: 9px !important;
    /*display: flex;
    align-items: center;*/
  }
  div[col-id='loading'] {
    width: 0 !important;
    &.ag-cell {
      border-left: 0px solid #eeeeee !important;
      /*border-right: 0px solid #eeeeee !important;*/
    }
  }
  .ag-root-wrapper {
    background-color: ${p => p.theme.white} !important;
    border: none !important;
    border-radius: 2px;
    font-family: 'Open Sans' !important;
    padding: 0 12px;
  }
  .ag-row-animation .ag-row {
    -webkit-transition: top 0.4s, height 0.4s, opacity 0.4s, box-shadow 0.4s,
      -webkit-transform 0.4s !important;
    transition: top 0.4s, height 0.4s, opacity 0.4s, box-shadow 0.4s,
      -webkit-transform 0.4s !important;
    transition: transform 0.4s, top 0.4s, height 0.4s, box-shadow 0.4s,
      opacity 0.4s !important;
    transition: transform 0.4s, top 0.4s, height 0.4s, box-shadow 0.4s,
      opacity 0.4s, -webkit-transform 0.4s !important;
  }
  .ag-body-horizontal-scroll-viewport {
    display: none;
  }

  .ag-row {
    /*border-color: transparent !important;*/
    border-radius: 0px !important;
    background-color: ${p => p.theme.white} !important;
    &:hover {
      background-color: ${p => p.theme.lightBlue} !important;
    }
    border: 0px solid transparent;
    border-bottom: 1px solid #e4e4e4;
    border-right: 1px solid #eeeeee;
  }

  .ag-row-odd {
    background-color: ${p => p.theme.oddRows} !important;
    &:hover {
      background-color: ${p => p.theme.lightBlue} !important;
    }
  }
`;
