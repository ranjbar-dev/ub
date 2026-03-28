import { Button } from '@material-ui/core';
import { AppPages } from 'app/constants';
import { replace } from 'connected-react-router';
import React from 'react';
import { useDispatch } from 'react-redux';
import styled from 'styled-components/macro';
import { StyleConstants } from 'styles/StyleConstants';

import { PageWrapper } from '../../components/PageWrapper';
import { LanguageSwitch } from '../LanguageSwitch';
import { ThemeSwitch } from '../ThemeSwitch';

export function NavBar() {
  const dispatch = useDispatch();
  const handleLogoutClick = () => {
    localStorage.clear();
    dispatch(replace(AppPages.RootPage));
  };
  return (
    <Wrapper>
      <div className="mainNavTitle">UB Admin</div>
      <Button
        onClick={handleLogoutClick}
        variant="outlined"
        className="logoutButton"
      >
        logout
      </Button>
      <PageWrapper>
        <ThemeSwitch />
        <LanguageSwitch />
      </PageWrapper>
    </Wrapper>
  );
}

const Wrapper = styled.header`
  box-shadow: 0 1px 0 0 ${p => p.theme.borderLight};
  height: ${StyleConstants.NAV_BAR_HEIGHT};
  display: flex;
  position: fixed;
  top: 0;
  width: 100%;

    align-items: center;
  background-color: ${p => p.theme.navBackground};
  z-index: 2;

  /* @supports (backdrop-filter: blur(10px)) {
    backdrop-filter: blur(10px);
    background-color: ${p =>
      p.theme.background.replace(
        /rgba?(\(\s*\d+\s*,\s*\d+\s*,\s*\d+)(?:\s*,.+?)?\)/,
        'rgba$1,0.75)',
      )};
  } */

  ${PageWrapper} {
    display: flex;
    align-items: center;
    justify-content: space-between;
  }
  .logoutButton{
	position: absolute;
    right: 20px;
	top: 15px;
	color:white;
	/*border-color:*/
  }
  .mainNavTitle{
	color: white;
    font-weight: 700;
    font-size: 20px;
    padding-left: 43px;
    font-family: 'Open Sans';
  }
`;
