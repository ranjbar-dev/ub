import styled from 'styled-components/macro';
import {StyleConstants} from 'styles/StyleConstants';

export const FullWidthWrapper=styled.div`
  padding-top: ${StyleConstants.NAV_BAR_HEIGHT};
  width: calc(100% - ${StyleConstants.SIDE_NAV_WIDTH});
  background: ${p => p.theme.greyBackground};
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: auto;

  &.column {
    flex-direction: column;
  }
  &.noAlignment {
    display: block;
    align-items: unset;
    justify-content: unset;
  }
  box-shadow: inset 3px 0px 13px rgba(0, 0, 0, 0.06);

  div[col-id="address"],
  div[col-id="txId"],div[col-id="from"],div[col-id="to"] {
    -webkit-user-select: text;
    -moz-user-select: text;
    -ms-user-select: text;
    user-select: text;
}

`;
