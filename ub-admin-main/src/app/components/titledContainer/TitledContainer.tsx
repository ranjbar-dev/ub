import React, { memo } from 'react';
import styled from 'styled-components/macro';
import { StyleConstants } from 'styles/StyleConstants';

/**
 * Titled card container with a sticky header bar.
 * Used as the top-level page layout wrapper for most admin pages.
 *
 * @example
 * ```tsx
 * <TitledContainer title="Users" id="usersPage">
 *   <SimpleGrid ... />
 * </TitledContainer>
 * ```
 */
function TitledContainer(props: { title: string; children: React.ReactNode; id?: string }) {
  return (
    <Wrapper id={props.id ? props.id : 'titledContainerID'}>
      <div className="whiteHead"></div>
      <div className="title">{props.title}</div>
      <div className="container">{props.children}</div>
    </Wrapper>
  );
}

export default memo(TitledContainer);
const Wrapper = styled.div`
  height: 95%;
  background: ${p => p.theme.white};
  width: 95%;
  border-radius: ${StyleConstants.CARD_BORDER_RADIUS};
  position:relative;
.whiteHead{
	min-width: 1165px;
    height: 20px;
    position: absolute;
    top: 31px;
    left: 0;
}
  @media screen and (max-width: 1471px) {
background:transparent;
.whiteHead{
	background:white;
}
  }
  &#withdrawals,&#balances,&#marketTicks,&#scanBlock{
	   background: transparent;
	   
  }
  .title {
    color: ${p => p.theme.blackText};
    font-size: 18px;
    font-weight: 600;
    padding: 10px 0px;
	margin-top: -15px;
    margin-bottom: 15px;
    /*background: ${p => p.theme.greyHead};*/
	background: ${p => p.theme.greyBackground};
    border-top-left-radius: ${StyleConstants.CARD_BORDER_RADIUS};
    border-top-right-radius: ${StyleConstants.CARD_BORDER_RADIUS};
  }
  .container {
    padding: 0 !important;
    padding-top: 15px !important;
	/*background: ${p => p.theme.greyBackground};*/
	background:white;
  min-width: 1165px !important;

  }
`;
