import { createGlobalStyle } from 'styled-components/macro';
/* istanbul ignore next */
export const GlobalStyle = createGlobalStyle`
 html,
body,
.UnitedBitAdmin {
    height: 100%;
    width: 100%;

}

html{
	overflow-x: hidden;
}
body {
    font-family: 'Helvetica Neue', Helvetica, Arial, sans-serif;
    background-color: ${p => p.theme.background};
    padding:0;
}

* {
    font-family: 'Open Sans', sans-serif, 'Helvetica Neue', Helvetica, Arial !important;
}

.submitButton {}

.redButton {
    background: #E18E8F !important;
}

.blackButton {
    background: #A3AFBB !important;
}

.skyBlueButton {
    background: #7DA6E9 !important;
}

.greenButton {
    background: #39A559 !important;
}

.lightGreenButton {
    background: #90E8C5 !important;
}

.lightYellowButton {
    background: #F6F3E4 !important;
}

.veryLightGreenButton {
    background: #DCF6EC !important;
}

.veryLightBlueButton {
    background: #F1F5FC !important;
}

.redButton,
.blackButton,
.skyBlueButton,
.greenButton,
.lightGreenButton {
    font-size: 12px;
    font-weight: 600;
    color: white !important;
}

.lightYellowButton,
.veryLightGreenButton,
.veryLightBlueButton {
    span {
        font-size: 12px;
        color: black;
        padding: 4px 0;
    }
    border:1px solid #eeee
}

.greyButton {}

.MuiButton-contained {
    box-shadow: none !important;
}

.MuiButton-root {
    text-transform: none !important;
    font-family: "Open Sans" !important;
}

.expandIcon {
    transform: rotate(-90deg);
    &.rotated {
        transform: rotate(0deg);
    }
}

.MuiTypography-root {
    color:${p => p.theme.blackText};
}

div[col-id="ip"],
div[col-id="id"],
div[col-id="email"],
div[col-id="firstName"],
div[col-id="lastName"],
/* div[col-id="address"], */
.value {
    -webkit-user-select: text;
    -moz-user-select: text;
    -ms-user-select: text;
    user-select: text;
}

::-webkit-scrollbar {
    width: 6px;
    height: 6px;
}

.ag-header-cell-sorted-none {
    pointer-events: none !important;
    cursor: none;
}


/*vertical Track */

::-webkit-scrollbar-track:vertical {
    /* box-shadow: inset 0 0 3px grey; */
    border-radius: 10px;
}


/*vertical Handle */

::-webkit-scrollbar-thumb:vertical {
    background: ${p => p.theme.darkGrey};
    border-radius: 10px;
}


/*horizontal Track */

::-webkit-scrollbar-track:horizontal {
    /* box-shadow: inset 0 0 3px grey; */
    border-radius: 10px;
}


/*horizontal Handle */

::-webkit-scrollbar-thumb:horizontal {
    background: ${p => p.theme.darkGrey};
    border-radius: 10px;
}


/* Handle on hover */

::-webkit-scrollbar-thumb:hover {
    background: ${p => p.theme.darkGrey};
}

::selection {
    background: ${p => p.theme.textBlue};
    color: ${p => p.theme.white};
}

::-moz-selection {
    background: ${p => p.theme.textBlue};
    color: ${p => p.theme.white};
}

* {
    scrollbar-width: thin;
}

.clickableRows {
    .ag-cell {
        cursor: pointer;
    }
}

.NWindow {
    .ag-cell {
        font-size: 12px !important;
    }
}

.loadingGridRow {
    position: fixed;
    top: 0px;
    width: 100%;
    left: 0px;
    overflow: hidden;
    background: rgba(17, 83, 126, 0.3);
    height: 100%;
    border-radius: 5px;
    &:active:after {
        opacity: 0;
    }
    &:after {
        animation: shine 8s ease-in-out infinite;
        animation-fill-mode: forwards;
        content: "";
        position: absolute;
        top: -110%;
        left: -100%;
        width: 200%;
        height: 200%;
        opacity: 0;
        transform: rotate(30deg);
        background: rgba(255, 255, 255, 0.13);
        background: linear-gradient( to right, rgba(255, 255, 255, 0.5) 0%, rgba(255, 255, 255, 1) 77%, rgba(255, 255, 255, 0.5) 92%, rgba(255, 255, 255, 0.0) 100%);
    }
}

@keyframes shine {
    10% {
        opacity: 1;
        top: -30%;
        left: 10%;
        transition-property: left, top, opacity;
        transition-duration: 0.7s, 0.7s, 0.15s;
        transition-timing-function: ease;
    }
    100% {
        opacity: 0;
        top: -30%;
        left: 10%;
        transition-property: left, top, opacity;
    }
}

.Toastify__toast--info {
    background: ${p => p.theme.white};
}

.Toastify__toast--error {
    background: ${p => p.theme.white};
}

.Toastify__toast--warning {
    background: ${p => p.theme.white};
}

.Toastify__toast--success {
    background: ${p => p.theme.white};
}

.Toastify__toast-body {
    margin: auto 0;
    -ms-flex: 1;
    flex: 1;
    color: ${p => p.theme.textGrey}!important;
}

.Toastify__close-button {
    color:${p => p.theme.textGrey};
    font-weight: bold;
    font-size: 14px;
    background: transparent;
    outline: none;
    border: none;
    padding: 0;
    cursor: pointer;
    opacity: 0.7;
    transition: 0.3s ease;
    -ms-flex-item-align: start;
    align-self: flex-start;
}

.img-fork {
    position: absolute;
    width: 130px;
    top: 0;
    right: 0;
}

.container-fluid {
    padding-right: 15px;
    padding-left: 15px;
    margin-right: auto;
    margin-left: auto;
    display: flex;
    align-items: center;
    height: 100%;
}

.navbar-brand {
    color: white;
    font-size: 2rem;
    margin-right: 24px;
}

.bagde {
    margin-left: 6px;
    display: flex;
    align-items: center;
}

.github {
    position: absolute;
    right: 15px;
}

.container {
    padding: 1.5em 2em 2em 2em;
}

.wrap {
    display: flex;
}

.img-list-wrap {
    flex: 1;
}

.img-list {
    display: flex;
    justify-content: flex-start;
    &.hide {
        display: none;
    }
}

.img-item>img {
    width: 100%;
    height: 100%;
    cursor: pointer;
}

.footer {
    position: fixed;
    bottom: 0;
    background-color: #0b2f3d;
    width: 100%;
}

.container-footer {
    padding: 24px;
    text-align: center;
}

.signature {
    color: white;
}

.container {
    position: relative;
}

.inline-container {
    display: none;
    max-width: 600px;
    // height: 400px;
    &.show {
        display: block;
    }
}

.options {
    margin-top: 12px;
    width: 100%;
    max-width: 250px;
    margin-right: 48px;
}

.options-list {
    margin-top: 12px;
    height: 440px;
    overflow: auto;
}

@keyframes zoo {
    0% {
        opacity: 0;
        transform: scale(0.2);
    }
    100% {
        opacity: 1;
        transform: scale(1);
    }
}

.MuiDialog-root {
    background: rgba(0, 0, 0, 0.2) !important;
}

.UbToast {
    position: fixed;
    z-index: 1;
    pointer-events: none;
    transition-property: bottom, opacity;
    transition-duration: 0.3s, 0.3s;
    transition-timing-function: ease;
    width: 100%;
    display: flex;
    justify-content: center;
    bottom: -40px;
    opacity: 0;
    &.show {
        bottom: 40px;
        opacity: 1;
        &.secound {
            bottom: 100px;
        }
    }
    .content {
        padding: 10px 20px;
        max-width: fit-content;
        border-radius: 5px;
        color: black;
        &.success {
            background: rgb(195, 236, 187);
        }
        &.error {
            background: rgb(221, 81, 69);
        }
    }
}

.withdrawModal {
    width: 100%;
    height: 100%;
    display: flex;
    flex-direction: column;
    .wrapperTitle {
        width: 100%;
        height: 40px;
        display: flex;
        align-items: center;
        background: rgb(213, 217, 230);
        padding: 0 10px;
        font-size: 14px;
        font-weight: 600;
    }
    .content {
        display: flex;
        width: 100%;
        flex: 1;
        background: #f5f6f8;
        justify-content: space-between;
        .MuiFormControlLabel-label {
            font-size: 13px !important;
        }
        .detailsContainer {
            background: rgb(245, 246, 248);
            flex: 5;
        }
        .actionsContainer {
            flex: 3;
            display: flex;
            flex-direction: column;
            align-items: center;
            background: #f5f6f8;
            max-width: 400px;
            .inp1,
            .inp2 {
                width: 60%;
                margin: 15px 0 0 0;
                input {
                    background: white;
                }
            }
            .inp2 {
                display: flex;
                align-items: center;
                margin: 9px 0 0 0;
                .fee {
                    min-width: 50px;
                    font-size: 12px;
                }
            }
            .sendMoney {
                button {
                    min-width: 238px;
                }
            }
            .rejectCancel {
                width: 246px;
                display: flex;
                justify-content: space-around;
                margin: 8px 0;
                button {
                    min-width: 116px;
                }
            }
            .autoTransfer {
                margin-left: 69px;
                min-height: 44px;
            }
            .bottomActions {
                display: flex;
                justify-content: space-around;
                margin: 50px 0px 0px 0;
                button {
                    margin: 0px 8px;
                    padding: 0 10px;
                }
                .loadingCircle {
                    top: 5px !important;
                }
            }
            .actionDetailRows {
                .detailRow {
                    background: #f5f6f8;
                }
            }
        }
        .detailsContainer {
            padding: 10px 20px;
            max-width: 550px;
            .detailRowsContainer {}
        }
        .detailRow {
            display: flex;
            min-height: 35px;
            align-items: center;
            background: white;
            border-bottom: 1px solid #eee;
            min-width: 320px;
            &.last {
                border-bottom: none;
            }
            .title {
                min-width: 150px;
                min-width: 150px;
                font-size: 12px;
                padding: 0 12px;
                font-weight: 600;
                color: #6c757e;
            }
            .value {
                flex: 1;
                min-width: 150px;
                min-width: 150px;
                font-size: 12px;
                padding: 0 12px;
                font-weight: 600;
            }
            &.small {
                .value {
                    text-align: end;
                }
            }
        }
        .actionDetailRows {
            display: flex;
            flex-direction: column;
            align-items: center;
        }
    }
}

.MuiDialog-paperScrollPaper {
    background: white;
    border-radius: 7px;
}

.rejectPopup {
    width: 600px;
    height: 340px;
    display: flex;
    flex-direction: column;
    padding: 10px 10px;
    .title {
        font-size: 13px;
        font-weight: 600;
        color: #5d5d5d;
    }
    .inp {
        margin-top: 60px;
    }
    .buttons {
        flex: 1;
        display: flex;
        place-items: flex-end;
        align-self: flex-end;
    }
}

.countryContainer {
    cursor: pointer;
    img {
        max-width: 23px;
        max-height: 15px;
        border-radius: 3px;
    }
    span {
        color: ${p => p.theme.textGrey};
        font-weight: 600;
        margin: 0 5px;
        font-size: 13px !important;
    }
    input {
        font-weight: 600;
    }
}

.MuiAutocomplete-hasPopupIcon.MuiAutocomplete-hasClearIcon .MuiAutocomplete-inputRoot[class*="MuiOutlinedInput-root"] {
    padding: 7px 0px !important;
}

.MuiAutocomplete-paper {
    min-width: 320px !important;
}

.MuiAutocomplete-inputRoot[class*="MuiOutlinedInput-root"] .MuiAutocomplete-input {
    font-size: 14px !important;
    font-weight: 500 !important;
    color: ${p => p.theme.textGrey};
    &::placeholder {
        color: ${p => p.theme.textGrey} !important;
    }
}

.ag-theme-balham button {
    box-shadow: none !important;
    padding: 0 5px !important;
    margin-left: 4px !important;
    span {
        line-height: 26px !important;
        font-size: 11px !important;
    }
    .loadingCircle {
        top: 2px !important;
    }
}

.AdminReports__Wrapper {
    width: 100%;
    padding: 24px;
    height: 100%;
}

div.adminReportsInput {
    width: 100%;
    height: fit-content;
    display: flex;
    flex-direction: column;
    align-items: flex-end;
}

div.adminReportsComments {
    width: 100%;
    height: calc(100% - 130px);
    overflow: auto;
    .mainCommentWrapper {
        width: 100%;
        margin-bottom: 24px;
        padding: 0 10px;
        .bold.grey {
            font-size: 12px;
            font-weight: 700;
            color: ${p => p.theme.textGrey};
        }
    }
    .commentWrapper {
        background: #f1f5fc;
        border-radius: 7px;
    }
    .actions {
        display: flex;
        place-content: flex-end;
    }
    .comment {
        padding: 12px 24px;
        width: 100%;
    }
}

.MuiPaginationItem-outlined {
    border: none !important;
    color: #595C5E;
}

.MuiPaginationItem-page.Mui-selected {
    background-color: #C0B8E4 !important;
    border-radius: 8px;
    color: white;
}

.MuiPaginationItem-page {
    min-width: 22px;
    max-height: 22px;
}

.MuiOutlinedInput-notchedOutline {
    border-color: #E6E3E3;
}

.MuiTabs-indicator {
    background-color: transparent !important;
}

.MuiTab-wrapper {
    font-size: 11px !important;
    font-weight: 600;
}

.NWindow {
    .simpleGrid {
        box-shadow: none !important;
    }
    .MuiInputLabel-outlined {
        font-size: 11px !important;
    }
}

.udInput {
    border-radius: 3px;
    border: 1px solid #868686;
    max-width: 150px;
    position: absolute;
}

.imageWrapperPlaceHolder {
    width: 100%;
    height: 96%;
    background: #B8B8B8;
    margin-top: 1%;
    border-radius: 7px;
}

.wideDrop {
    fieldset {
        min-width: 115px;
    }
}

.documentIdInput {
    max-width: 125px;
    border-color: #E6E3E3;
    border: 1px solid #E6E3E3;
    border-radius: 3px;
    min-height: 20px;
}

.MuiListItem-root.MuiMenuItem-root.Verified {
    color: ${p => p.theme.green}!important;
}

.MuiListItem-root.MuiMenuItem-root.MuiMenuItem-root.Not {
    color: ${p => p.theme.red}!important;
}

.Verified {
    .MuiOutlinedInput-inputMarginDense {
        color: ${p => p.theme.green}!important;
    }
}

.Not {
    .MuiOutlinedInput-inputMarginDense {
        color: ${p => p.theme.red}!important;
    }
}

.MuiTypography-body1 {
    font-family: 'Open Sans' !important;
}

.MuiTab-root {
    min-height: 38px !important;
}

.adminPaymentComment {
    background: white;
    margin-top: 12px;
    height: 85px;
    padding: 6px 12px;
    font-size: 13px;
    .adminName,
    .commentDate {
        font-size: 10px;
        font-weight: 600;
        color: #5a5a5a;
    }
}

.MuiOutlinedInput-root.narrow {
    max-height: 30px;
}

.ubInputContainer {
    width: 100%;
    display: flex;
    align-items: center;
    margin-bottom: 12px;
    .MuiOutlinedInput-root {
        flex: 1;
    }
    .label {
        min-width: 120px;
        font-size: 13px;
        color: #535353;
    }
}

.checkBoxCont {
    img {
        width: 18px;
        max-width: 18px;
        position: absolute;
    }
}

.MuiMenuItem-root {
	font-size: 12px !important;
    font-weight: 600;
}

.MuiInputBase-root {
    font-size: 13px !important;
    font-weight: 600 !important;
    font-family: 'Open Sans' !important;
}



/*div.ag-cell[col-id="createdAt"],
div.ag-cell[col-id="updatedAt"],
div.ag-cell[col-id="country"],
div.ag-cell[col-id="registrationDate"],
div.ag-cell[col-id="referralId"],
div.ag-cell[col-id="email"],
div.ag-cell[col-id="registeredIP"] {
    padding-left: 5px !important;
}*/

.MuiInputBase-adornedEnd {
    background: white !important;
}

.imageGridView {
    .ag-header {
        background-color: #e4eef6 !important;
        .ag-header-cell-text {
            font-size: 12px !important;
        }
    }
		.selectedRow{
			border-left:2px solid blue !important;
		}
}

.rowDataContainer {
    .ddown .MuiSelect-outlined.MuiSelect-outlined {
        font-size: 12px !important;
    }
}

.CountryDropDownWrapper {
		min-width: 200px;
		margin-top: -11px;
    .MuiAutocomplete-inputRoot[class*="MuiOutlinedInput-root"] .MuiAutocomplete-input:first-child {
        padding-top: 0px;
        padding-bottom: 0px;
    }
}
.t6{
	.loadingCircle{
		top:6px !important;

	}
}
#UserDetailsWindowSimpleGridWrapper{
    min-width:1151px !important;
}
.withdrawWindowConfirmPopupWrapper{
    width: 235px;
    display: flex;
    flex-direction: column;
    height: 130px;
    align-items: center;
    padding: 12px;
    justify-content: space-between;
    padding-top: 0;
}
`;
