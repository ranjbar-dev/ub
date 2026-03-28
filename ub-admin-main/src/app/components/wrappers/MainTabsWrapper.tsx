import styled from 'styled-components/macro';

export const MainTabsWrapper = styled.div`
  position: relative;
  width: 100%;
  height: calc(100% - 20px);
  .mainTab {
    border-top-right-radius: 15px !important;
    border-top-left-radius: 15px !important;
    /*background: ${p => p.theme.white} !important;*/
    background: #F5F6F8 !important;
    border: 2px solid #dddddd;
	&.Mui-selected{
		  background: #E4EEF6 !important;
		  border-bottom:none;
		  .MuiTab-wrapper{
			  font-size:11px;
			  font-weight:600;
		  }
	}
  }
  &.verificateWindow{
	.mainTab {
    border-top-right-radius: 15px !important;
    border-top-left-radius: 15px !important;
    /*background: ${p => p.theme.white} !important;*/
    background: #F5F6F8 !important;
    border: 2px solid #dddddd;
	&.Mui-selected{
		  background: white !important;
		  border-bottom:none;
		  .MuiTab-wrapper{
			  font-size:11px;
			  font-weight:600;
		  }
	}
  }
  }
  .MuiAppBar-colorPrimary {
    color: ${p => p.theme.blackText} !important;
    background-color: transparent !important;
    box-shadow: none !important;
	max-height:38px !important;
	border-bottom:2px solid #E7E7E7;
  }
`;
