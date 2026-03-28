import { createGlobalStyle, css } from 'styles/styled-components';

const GlobalStyle = createGlobalStyle`

:root {
    /* styles */
    --maxHeight: 850px;
    /* colors */
	--expandedShaddowColor:rgba(0,0,0,0.1);
	--interactiveTooltipBackground:#292A2C !important;
	--textGrey: #818181;
	--tradeGridTextColor:var(--blackText);
    --appHeaderSelectedColor: var(--primary);
    --blackText: #3D3E40;
    --blueSelect: #b5d5f7;
    --BoxShaddow: rgba(0, 0, 0, 0.2);
    --cardBorderRadius: 5px;
    --darkGrey: #C1C1C1;
    --darkRed: #BA3C3C;
    --darkTrans: rgba(0, 0, 0, 0.4);
    --dateShadow: 4px 5px 8px 5px #d8d8d8 !important;
    --disabledInput: rgba(0, 0, 0, 0.1);
    --disabledText: #C1C1C1;
    --dragIconColor: var(--blackText);
    --dropDownIconBackGround: #D8D8D8;
    --greenStatus: #86E0B9;
    --greenText: #06BA61;
    --greyBackground: #f0f1f3;
    --inputBorderColor: #D0D0D0;
    --lightBlue: #F9FAFE;
    --lightGreen: #F5FFF5;
    --lightGrey: #D8D8D8;
    --lightRed: #FFE5EC;
    --lumo-border-radius-m: 7px;
    --miniCancelButtonBackground: #e6e6e6;
    --miniCancelButtonTextColor: #28292B;
    --oddRows: #F8F8F8;
    --orange: #F49806;
    --ordersMiniGridRowBackground:#f7fafd;
    --placeHolderColor: #C1C1C1;
    --primary: #396DE0;
    --redText: #E64141;
    --selectHover: #C1C1C1;
    --statusTextColor: #707070;
    --swithBackBround: #C1C1C1;
    --textBlue: #396DE0;
    --tradeCellFontWeight: 600;
    --white: #ffffff;
    --yearScroll: #C1C1C1;
	.WhiteHeader{
		background:white !important;
		.MuiSelect-selectMenu{
			background:white !important;
		}
	}
    --verticalWhiteGradient: rgba(255, 255, 255, 1) 0%, rgba(255, 255, 255, 0.6) 45%, rgba(255, 255, 255, 0.1) 90%, rgba(255, 255, 255, 0) 100%
}

.darkTheme {
	--dragIconColor: #787b86;
	--expandedShaddowColor:rgba(255,255,255,0.1);
	--interactiveTooltipBackground:#292A2C !important;
	--selectHover: #16161A;
	--tradeGridTextColor:#cacaca;
    --appHeaderSelectedColor: var(--primary);
    --blackText: #F2F2F2;
    --blueSelect: #396de0;
    --BoxShaddow: rgba(255, 255, 255, 0.2);
    --darkGrey: #525252;
    --darkRed: #BA3C3C;
    --darkTrans: rgba(255, 255, 255, 0.4);
    --dateShadow: 4px 5px 8px 5px #111111 !important;
    --disabledInput: rgba(255, 255, 255, 0.1);
    --disabledText: #363636;
    --dropDownIconBackGround: #707070;
    --greenText: #06BA61;
    --greyBackground: #292A2C;
    --inputBorderColor: #585858;
    --lightBlue: #23232B;
    --lightGreen: #25302B;
    --lightGrey: #313131;
    --lightRed: #321825;
    --lumo-base-color: var(--white);
    --lumo-body-text-color: var(--textGrey);
    --lumo-header-text-color: var(--textGrey);
    --lumo-primary-text-color: var(--textBlue);
    --lumo-shade-5pct: #4e4e4e;
    --lumo-tertiary-text-color: #848282;
    --middleIconWidth: 336px;
    --miniCancelButtonBackground: #28292B;
    --miniCancelButtonTextColor: #BFBFBF;
    --oddRows: #16161A;
    --orange: #F49806;
    --ordersMiniGridRowBackground:#141421;
    --placeHolderColor: #636363;
    --redText: #E64141;
    --statusTextColor: var(--white);
    --swithBackBround: #707070;
    --textBlue: #00A7FF;
    --textGrey: #BFBFBF;
    --tradeCellFontWeight: 500;
    --white: #1C1C21;
    --yearScroll: #28292B;
	.WhiteHeader{
		background:#1C1C21 !important;
		.MuiSelect-selectMenu{
			background:#1C1C21 !important;
		}
	}
    --verticalWhiteGradient: rgba(28, 28, 33, 1) 0%, rgba(28, 28, 33, 0.6) 45%, rgba(28, 28, 33, 0.1) 90%, rgba(28, 28, 33, 0) 100%
}

.MuiSkeleton-wave::after {
    background: linear-gradient( 90deg, transparent, var(--lightGrey), transparent) !important;
    animation: MuiSkeleton-keyframes-wave 1s linear 0.5s infinite !important;
}


/* buttons */

.MuiButton-contained.Mui-disabled {
    color: var(--disabledText) !important;
    box-shadow: none;
    span {
        color: var(--disabledText) !important;
    }
    background-color: rgba(0, 0, 0, 0.12);
}

.MuiIconButton-root:hover {
    background-color: var(--oddRows) !important;
}


/* buttons end */


/* listItems */

.MuiMenuItem-root {
    color: var(--blackText) !important;
}

.MuiListItem-root:hover {
    background-color: var(--selectHover) !important;
}

.MuiListItem-root.Mui-selected {
    background-color: var(--lightBlue) !important;
}


/* listItems end */


/* input */

input:-webkit-autofill,
input:-webkit-autofill:hover,
input:-webkit-autofill:focus,
textarea:-webkit-autofill,
textarea:-webkit-autofill:hover,
textarea:-webkit-autofill:focus,
select:-webkit-autofill,
select:-webkit-autofill:hover,
select:-webkit-autofill:focus {
    border: none !important;
    -webkit-text-fill-color: var(--blackText) !important;
    -webkit-box-shadow: none !important;
    box-shadow: none !important;
    background: var(--white);
    transition: background-color 5000s ease-in-out 0s !important;
}


/* Chrome, Safari, Edge, Opera */

input::-webkit-outer-spin-button,
input::-webkit-inner-spin-button {
    -webkit-appearance: none;
    margin: 0;
}


/* Firefox */

input[type=number] {
    -moz-appearance: textfield;
}


/* input end */


/* filter wrapper */

.MuiMenu-paper {
    li {
        font-size: 14px !important;
        font-weight: 500;
    }
}


/* filter wrapper end */

.red {
    color: var(--redText);
}

.MuiPaper-root {
    background-color: var(--white) !important;
}

.MuiTab-textColorPrimary {
    color: var(--textGrey) !important;
}

.MuiTab-textColorPrimary.Mui-selected {
    color: var(--primary) !important;
}

.MuiButton-outlined.Mui-disabled {
    border: 0px solid var(--oddRows) !important;
    color: var(--textGrey) !important;
    background: var(--oddRows) !important;
    font-size: 13px !important;
    border-radius: 7px !important;
}

.MuiButton-outlinedPrimary {
    font-size: 13px !important;
    border-radius: 7px !important;
}

.MuiInputBase-root {
    background: var(--white) !important;
    color: var(--blackText) !important;
    legend {
        zoom: 0.72;
        margin-left: 0px;
    }
}

.MuiInputBase-root.select {
    background: transparent !important;
    legend {
        zoom: 1;
        margin-left: -1px;
    }
}

.MuiList-root {
    background: var(--white) !important;
}

.MuiOutlinedInput-notchedOutline {
    border-color: var(--inputBorderColor) !important;
    transition: border-color 0.2s !important;
}

.MuiSelect-outlined.MuiSelect-outlined {
    font-weight: 600;
}

.MuiFormControl-root,
.MuiInputBase-root.select {
    &:hover {
        .MuiOutlinedInput-notchedOutline {
            border-color: var(--textBlue) !important;
        }
    }
}

.MuiSvgIcon-root {
    path {
        fill: var(--darkGrey) !important;
    }
}

.MuiOutlinedInput-root.Mui-focused .MuiOutlinedInput-notchedOutline {
    border-color: var(--textBlue) !important;
    border-width: 1px !important;
}

.MuiOutlinedInput-root.Mui-error .MuiOutlinedInput-notchedOutline {
    border-color: var(--redText) !important;
}


/* .MuiFormControl-fullWidth{
	  min-width:278px !important;
	  max-width:276px !important
	} */

.MuiTextField-root {
    label {
        span {
            color: var(--placeHolderColor) !important;
            font-size: 13px;
        }
        &.Mui-focused {
            span {
                color: var(--textBlue) !important;
            }
        }
        &.Mui-error {
            span {
                color: var(--redText) !important;
            }
        }
    }
}

.Mui-error {
    .MuiOutlinedInput-notchedOutline {
        border-color: var(--redText) !important;
    }
    &:hover {
        .MuiOutlinedInput-notchedOutline {
            border-color: var(--redText) !important;
        }
    }
}

.expandMore {
    width: 0;
    height: 0;
    position: relative;
    svg {
        min-width: 20px;
        min-height: 20px;
    }
}

.badge {
    background: var(--lightBlue);
    color: var(--textBlue);
    margin-right: 10px;
    padding: 5px 12px;
    font-size: 10px;
    font-weight: 700;
    border-radius: 29px;
}

.ag-theme-balham .ag-root-wrapper {
    background-color: var(--white) !important;
}

.ag-header-cell-text {
    color: var(--textGrey) !important;
    font-size: 13px;
}

div[col-id="address"],
div[col-id="txId"],
.value {
    -webkit-user-select: text;
    -moz-user-select: text;
    -ms-user-select: text;
    user-select: text;
}

.whiteShaddow {
    transition: opacity 0.5s;
    background: linear-gradient( 0deg, var(--verticalWhiteGradient));
}

html,
body {
    height: 100%;
    width: 100%;
    line-height: 1.5;
    min-width: 100%;
    width: fit-content;
}

html {
    background: var(--white);
    &.htmldark{
      background-color:#292A2C;
    }
}

body {
    font-family: 'Open Sans', sans-serif, 'ar', 'Helvetica Neue', Helvetica, Arial, sans-serif;
    font-size: 14px;
}

body.fontLoaded {
    font-family: 'Open Sans', sans-serif, 'ar', 'Helvetica Neue', Helvetica, Arial;
}

body.arabic {
    font-family: 'ar', 'Open Sans';
    direction: rtl;
}

#unitedBit {
    background-color: var(--greyBackground);
    min-height: 100%;
    min-width: 100%;
    &.darkTheme {
        background-color: #292A2C;
    }
}

p,
label {
    font-family: 'Open Sans', 'ar', Georgia, Times, 'Times New Roman', serif;
    /* line-height: 1.5em; */
}

 ::-webkit-scrollbar {
    width: 5px;
    height: 5px;
}


/*vertical Track */

::-webkit-scrollbar-track:vertical {
    /* box-shadow: inset 0 0 3px grey; */
    border-radius: 10px;
}


/*vertical Handle */

::-webkit-scrollbar-thumb:vertical {
    background: var(--darkGrey);
    border-radius: 10px;
}


/*horizontal Track */

::-webkit-scrollbar-track:horizontal {
    /* box-shadow: inset 0 0 3px grey; */
    border-radius: 10px;
}


/*horizontal Handle */

::-webkit-scrollbar-thumb:horizontal {
    background: var(--darkGrey);
    border-radius: 10px;
}


/* Handle on hover */

::-webkit-scrollbar-thumb:hover {
    background: var(--darkGrey);
}

::selection {
    background: var(--textBlue);
    color: var(--white);
}

::-moz-selection {
    background: var(--textBlue);
    color: var(--white);
}

* {
    scrollbar-width: thin;
}

.grecaptcha-badge {
    visibility: hidden;
}

.react-slideshow-container+div.indicators {
    direction: ltr;
}

.react-slideshow-container {
    direction: ltr;
}

.Toastify__toast {
    border-radius: 5px !important;
    font-family: 'Open Sans' !important;
    font-size: small;
    box-shadow: 3px 4px 9px 4px var(--BoxShaddow);
}

.Toastify__toast-body {
    text-align: center;
    display: flex;
    align-items: center;
}

.Toastify__toast-container--top-right {
    top: 60px;
}

.MuiOutlinedInput-root {
    border-radius: 7px !important;
}

.MuiButton-root {
    letter-spacing: 0 !important;
    border-radius: 7px !important;
    box-shadow: none !important;
    /* span{
	  font-weight:600;
	} */
}

.MuiButton-text {
    padding: 4px 12px !important;
}

.MuiButton-containedPrimary:hover {
    background-color: #063295 !important;
}

.MuiDivider-root {
    background-color: var(--greyBackground) !important;
}

.MuiSelect-select:focus {
    border-radius: 7px !important;
    background-color: transparent !important;
}

.MuiListItemText-root {
    margin-top: 6px !important;
}

span {
    text-transform: none !important;
    font-size: 14px;
}

.black {
    color: var(--blackText);
    span {
        color: var(--blackText);
    }
}

.centerHor {
    display: flex;
    justify-content: center;
    align-items: center;
}

.centerText {
    text-align: center;
}

.fl1 {
    flex: 1;
    align-items: center;
    max-height: 30px;
}

.flexSpacer1 {
    flex: 1;
    max-height: 12px;
    &.mh40 {
        max-height: 40px;
    }
}

.fl2 {
    flex: 2;
}

.fl3 {
    flex: 3;
}

.fl4 {
    flex: 4;
}

.bold {
    font-weight: 600 !important;
}

.p5 {
    padding: 5px;
}

.restrictedInput {
    max-width: 274px;
    min-width: 278px !important;
}

.restrictedInputWrapper {
    display: flex;
    justify-content: center;
    flex-direction: column;
    align-items: center;
    &.spaceAround {
        justify-content: space-around;
    }
}

.simpleRoundButton {
    border-radius: 40px !important;
    padding: 4px 12px !important;
    background: var(--lightBlue) !important;
    border: 1px solid var(--lightBlue) !important;
    span {
        color: var(--textBlue);
    }
    &:hover {
        border: 1px solid var(--textBlue) !important;
    }
}

.transParentRoundButton {
    /* padding: 0 12px !important; */
    border-radius: 50px !important;
    &.MuiButton-textPrimary {
        span {
            color: var(--textBlue);
        }
        path {
            fill: var(--textBlue) !important;
        }
    }
}

.cancelButton {
    span {
        color: var(--textGrey) !important;
        font-weight: 600;
    }
    &:hover {
        /*background:  var(--oddRows) !important;
		  */
        background: transparent !important;
        span {
            text-decoration: underline;
        }
    }
    @media screen and (min-height: 700px) {
      margin-top:18px;
    }
}

.underlined {
    &:hover {
        background: transparent !important;
        span {
            text-decoration: underline;
        }
    }
}

.greenOutlined {
    border: 1px solid rgba(6, 186, 97, 1) !important;
    padding: 0 !important;
    background: rgba(144, 255, 141, 0.1) !important;
    border-radius: var(--cardBorderRadius) !important;
    &:hover {
        background: rgba(144, 255, 141, 0.3) !important;
    }
    span {
        color: var(--blackText);
        font-weight: 600;
        font-family: 'Open Sans';
        font-size: 12px;
    }
}

.simpleGreyButton {
    height: 32px;
    padding: 0 12px !important;
    background: var(--greyBackground) !important;
    &:hover {
        background: var(--lightGrey) !important;
    }
    span {
        color: var(--blackText) !important;
    }
}

.roundedRedButton {
    border-radius: 40px !important;
    height: 32px;
    padding: 0 12px !important;
    width: fit-content;
    background: var(--redText) !important;
    span {
        color: white !important;
    }
    &:hover {
        background: var(--darkRed) !important;
    }
}

.ubButton {
    margin: 0.5vh 0 !important;
    padding: 7px 25px 5px 25px !important;
    min-height: 40px;
}

.blue {
    color: var(--textBlue)
}

.ag-root {
    border: none !important;
}

.ag-header-cell::after {
    display: none !important;
}

.ag-header {
    border-top-left-radius: 7px;
    border-top-right-radius: 7px;
    background-color: var(--greyBackground) !important;
    border-bottom: 1px solid var(--greyBackground) !important;
}

.ag-row {
    border-color: transparent !important;
    border-radius: 7px !important;
    background-color: var(--white) !important;
    &:hover {
        background-color: var(--lightBlue) !important;
    }
}

.ag-row-odd {
    background-color: var(--oddRows) !important;
    &:hover {
        background-color: var(--lightBlue) !important;
    }
}

.ag-cell-focus {
    border: 1px solid transparent !important;
}

.ag-cell {
    /* display: flex !important;
		align-items: center !important; */
    line-height: 33px !important;
    font-size: 13px;
    color: var(--blackText);
    font-weight: 600;
}

.ag-theme-balham .ag-root-wrapper {
    border: none !important;
    border-radius: 2px;
    button {
        background: transparent !important;
        span {
            color: var(--textBlue) !important;
        }
        &.black {
            span {
                color: var(--blackText) !important;
            }
            path {
                fill: var(--blackText) !important;
            }
        }
        .rotated {
            path {
                fill: var(--textGrey) !important;
            }
        }
        &.grey {
            color: var(--textGrey) !important;
            span {
                color: var(--textGrey) !important;
            }
        }
    }
}

.ag-theme-balham {
    font-family: 'Open Sans' !important;
}

.ag-icon {
    color: var(--blackText);
}

.ag-row-animation .ag-row {
    -webkit-transition: top 0.4s, height 0.4s,  opacity 0.4s, box-shadow 0.4s, -webkit-transform 0.4s !important;
    transition: top 0.4s, height 0.4s,  opacity 0.4s, box-shadow 0.4s, -webkit-transform 0.4s !important;
    transition: transform 0.4s, top 0.4s, height 0.4s, box-shadow 0.4s, opacity 0.4s !important;
    transition: transform 0.4s, top 0.4s, height 0.4s,box-shadow 0.4s, opacity 0.4s, -webkit-transform 0.4s !important;
}

.ag-theme-balham .ag-root-wrapper {
    background-color: var(--white);
}

.withExpandableRows {
    .ag-row {
        width: calc(100% - 10px)!important;
        margin-left: 3px !important;
        .detailButton {
            position: absolute;
            top: 2px;
            max-width: 30px;
            max-height: 35px;
            padding: 0 8px !important;
            min-height: 32px;
            left: 2px;
            span {
                font-size: 13px;
            }
        }
    }
}

.statusBadge {
    height: 20px;
    background: var(--greenStatus);
    max-width: fit-content;
    min-width: fit-content;
    padding: 1px 8px;
    line-height: 17px;
    margin-top: 10px;
    font-size: 11px !important;
    border-radius: 4px;
    color: var(--statusTextColor);
    font-size: 12px;
    min-width: 61px;
    text-align: center;
    &.expired {
        background: #e08686;
    }
    &.mini {
        padding: 1px 0;
    }
}

button:focus {
    box-shadow: none !important;
}

.noRowsButtonWrapper {
    text-align: center;
    .noRowButton {
        bottom: calc(40vh - 140px);
    }
}

.mb12 {
    margin-bottom: 12px !important;
}

.mb24 {
    margin-bottom: 24px !important;
}

.mb2vh {
    margin-bottom: 2vh !important;
}

.pt1vh {
    padding-top: 1vh;
}

.MuiDialog-paper {
    border-radius: 10px !important;
}

.MuiDialogContent-root {
    padding: 0 !important;
}

.MuiDialog-paperWidthSm {
    max-width: 655px !important;
}

.alertConfirmWrapper {
    color: var(--textGrey);
    padding: 40px 60px 24px 60px !important;
    min-height: 90px !important;
}

.alertButtonsWrapper {
    display: flex;
    justify-content: center;
    padding-bottom: 24px;
    span {
        color: var(--blackText)!important;
    }
    .separator {
        height: 25px;
        width: 1px;
        background: var(--textGrey);
        margin: 0px 9px;
        margin-top: 4px;
    }
}

input,
textarea,
.MuiSelect-selectMenu {
    font-size: 13px !important;
}

.MuiSelect-select,
li.MuiMenuItem-root {
    font-size: 13px !important;
    font-weight: 600;
    span {
        font-size: 13px !important;
        font-weight: 600;
        color: var(--textGrey);
    }
}

.MuiMenu-paper {
    /* max-height:250px !important; */
    span {
        color: var(--blackText);
        font-size: 13px;
        font-weight: 600;
    }
    border-radius: 8px !important;
}

.MuiOutlinedInput-inputMarginDense {
    padding-top: 10.5px;
    padding-bottom: 10.5px;
    max-height: 40px;
    min-height: 19px;
}

.MuiButton-outlinedPrimary {
    color: var(--textBlue) !important;
    border: 1px solid !important;
}

.MuiInputLabel-outlined.MuiInputLabel-marginDense {
    margin-top: -0.5px !important;
    line-height: 15px;
}

.countryContainer {
    cursor: pointer;
    img {
        max-width: 23px;
        max-height: 15px;
        border-radius: 3px;
    }
    span {
        color: var(--textGrey);
        font-weight: 600;
        margin: 0 5px;
        font-size: 13px !important;
    }
    input {
        font-weight: 600;
    }
}

.mt12 {
    margin-top: 12px !important;
}

.MuiAutocomplete-paper {
    border-radius: 7px !important;
    min-width: 278px !important;
    margin-bottom: -35vh !important;
    max-height: 20vh !important;
    box-shadow: 0px 5px 5px -3px rgba(0, 0, 0, 0.2), 0px 8px 10px 1px rgba(0, 0, 0, 0.14), 0px 3px 14px 2px rgba(0, 0, 0, 0.12) !important;
    .MuiAutocomplete-listbox {
        max-height: 20vh !important;
    }
    .MuiAutocomplete-noOptions {
        color: var(--textGrey) !important;
    }
}

.autoComplete {
    width: 100%;
    input {
        font-weight: 600;
    }
}

input {
    &::placeholder {
        font-size: 13px !important;
        color: var(--textGrey) !important;
        font-weight: 600;
    }
}

.MuiAutocomplete-inputRoot[class*="MuiOutlinedInput-root"] .MuiAutocomplete-input {
    padding: 7.5px 4px !important;
&.filled{

		 padding-left:30px !important;

}
}
.countryInput.filled{
	.MuiAutocomplete-inputRoot[class*="MuiOutlinedInput-root"] .MuiAutocomplete-input {
		 padding-left:30px !important;
}
}

.MuiAutocomplete-inputRoot {
    max-height: 40px;
    input {
        margin-top: -6px;
    }
}

.ag-body-horizontal-scroll-viewport {
    overflow-x: scroll;
    display: none;
}

.ag-header {
    opacity: 0;
    /* transition-delay:0.2s; */
}

.hidenTextArea {
    opacity: 0;
    position: absolute;
}

.MuiSwitch-thumb {
    box-shadow: none !important;
}

.MuiButton-endIcon {
    margin-top: 1px !important;
    margin-left: 0px !important;
}

.checkIcon {
    .MuiCheckbox-root {
        background: transparent !important;
    }
    &:hover {
        path {
            fill: var(--textBlue);
        }
    }
}

.addressCoin {
    span {
        color: var(--blackText) !important;
    }
}

.NoRowsWrapper {
    .noRowsText {
        font-weight: 600;
    }
}

.miniGrid {
    .NoRowsWrapper {
        transform: scale(0.7);
        @media screen and (max-width: 1600px) {
            transform: scale(0.5);
            .noRowsText {
                span {
                    font-size: 18px;
                }
            }
        }
        margin-top: 10px;
        .noRowsText {
            margin-top: -10px !important;
        }
        &.miniNoRows {
            filter: grayscale(1);
        }
		.noRowsImage{

			transform: scale(0.7) !important;
    margin-bottom: -15px;

		}
    }
    .ag-header {
        background: transparent !important;
        font-size: 11px;
    }
    .ag-header-cell-text {
        font-size: 11px !important;
    }
    .whiteShaddow {
        display: none;
    }
    .statusBadge {
        background: unset;
        color: var(--blackText);
    }
}

.middleIcon {
    svg {
        max-width: var(--middleIconWidth);
    }
}

.layoutItem.Mui-selected {
    background-color: var(--blueSelect) !important;
}
.MuiTooltip-popperInteractive{
	.MuiTooltip-tooltip{
		background-color:var(--interactiveTooltipBackground);
		box-shadow:0 0 9px rgba(0, 0, 0, 0.77) !important;

		span{
			font-size:13px;
			font-weight: 500;
    color: #bfbfbf ;

		}

	}
	.MuiTooltip-arrow{
			color:#292A2C  !important;
		}
	hr{
		background-color:#5b5c5f !important;
	}
}
`;

export default GlobalStyle;
export const NarrowInputs = css`
  .MuiInputBase-root {
    max-height: 32px;
    &.select {
      margin-top: 8px;
      margin-bottom: -5px;
    }
  }

  .MuiInputLabel-outlined.MuiInputLabel-marginDense {
    transform: translate(14px, 7px) scale(1);
  }

  .MuiInputLabel-outlined.MuiInputLabel-shrink {
    transform: translate(14px, -6px) scale(0.75);
  }

  .MuiOutlinedInput-input {
    font-size: 14px;
  }
`;
