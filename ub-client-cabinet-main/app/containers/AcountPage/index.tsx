/*
 *
 * AcountPage
 *
 */

import React, { FC, memo, useEffect, useRef } from 'react';
import { Helmet } from 'react-helmet';
import { useSelector, useDispatch } from 'react-redux';
import { createStructuredSelector } from 'reselect';

import { useInjectSaga } from 'utils/injectSaga';
import { useInjectReducer } from 'utils/injectReducer';
import makeSelectAcountPage, { makeSelectLoggedIn } from './selectors';
import reducer from './reducer';
import saga from './saga';
import { MaxWidthWrapper } from 'components/wrappers/maxWidthWrapper';
import AcountPageUserInfoIcon from 'images/themedIcons/acountPageUserInfoIcon';
import { Themes, AppPages } from 'containers/App/constants';
import styled from 'styles/styled-components';
import { Divider, Button, Card } from '@material-ui/core';

import translate from './messages';
import { FormattedMessage } from 'react-intl';
// import { loginAction } from 'containers/LoginPage/actions';
import { push } from 'redux-first-history';
import { getUserDataAction, getNewVerificationEmailAction } from './actions';
import AnimateChildren from 'components/AnimateChildren';
import { GridLoading } from 'components/grid_loading/gridLoading';
import BreadCrumb from 'components/BreadCrumb';
import { censor } from 'utils/formatters';

import AcountPageUserSecurityIcon from 'images/themedIcons/acountPageUserSecurity';
import SecurityIcon from 'images/securityIcon/securityIcon';
import { SecurityLevel, KycStatus } from './constants';
import SecurityLevelBar from 'components/securityLevelBar/securityLevelBar';
import AcountPageDepositIcon from 'images/themedIcons/acountPageDepositIcon';
import bagIcon from 'images/bag_icon.svg';
import listIcon from 'images/list_icon.svg';
import historyIcon from 'images/history_icon.svg';
import checkIcon from 'images/check_icon.svg';
import redShield from 'images/redShieldIcon.svg';
import { toast } from 'components/Customized/react-toastify';
import { storage } from 'utils/storage';
import { LocalStorageKeys } from 'services/constants';
import throttle from 'utils/throttle';

const stateSelector = createStructuredSelector({
  acountPage: makeSelectAcountPage(),
  loggedIn: makeSelectLoggedIn(),
});

interface Props { }
interface infoRowProps {
  title: any;
  bigTitle?: boolean;
  value?: any;
  titleIcon?: any;
  description?: string;
  actionButton?: { buttonTitle: any; onClick: any; disabled?: boolean };
}

const InfoRow: FC<infoRowProps> = props => {
  const {
    title,
    actionButton,
    bigTitle,
    description,
    titleIcon,
    value,
  } = props;

  return (
    <div className='infoRow'>
      <div className={`rowItemTitle ${bigTitle ? 'bigTitle' : ''}`}>
        <span>{titleIcon}</span>
        <span>{title}</span>
      </div>
      <div className='rowItemValue'>{value}</div>
      {description && <span className='itemDescription'>{description}</span>}
      {actionButton && (
        <div className='rowActionButton'>
          <Button
            className={`${actionButton.disabled ? 'disabledOutline' : ''}`}
            disabled={actionButton.disabled}
            onClick={actionButton ? actionButton.onClick : () => { }}
            variant={'outlined'}
            size='small'
            color='primary'
          >
            {actionButton.buttonTitle}
          </Button>
        </div>
      )}
    </div>
  );
};

const throttleFunc = throttle(5000, false, functionToRun => {
  functionToRun();
});

function AcountPage(props: Props) {
  useInjectReducer({ key: 'acountPage', reducer: reducer });
  useInjectSaga({ key: 'acountPage', saga: saga });
  const { acountPage } = useSelector(stateSelector);
  const dispatch = useDispatch();
  //effect
  useEffect(() => {
    // if (!acountPage || !acountPage.userData) {
    setTimeout(() => {
      dispatch(getUserDataAction());
    }, 0);
    // }
  }, []);
  //page Data
  const isLoading = acountPage ? acountPage.isLoading : true;
  const passwordToShow = '**********';

  let emailToShow: string = '';
  let userId: string = '';
  let has2fa: boolean = false;
  let phoneToShow: string = '';
  let securityLevel: SecurityLevel = SecurityLevel.LOW;
  let securityLevelMessage: string = '';
  let emailVerified: boolean = false;
  let kycStatus: KycStatus = KycStatus.INCOMPLETE;
  let kycDescription: string = '';
  let kycLevel: string = 'none';
  let profileStatus: string = '';
  let googleAuth: boolean = false;
  let isVerified: boolean = false;
  if (acountPage && acountPage.userData) {
    emailToShow = censor({
      value: acountPage.userData.email,
      from: 3,
      to: acountPage.userData.email.length - 11,
    });
    userId = acountPage.userData.ubId;
    phoneToShow = censor({
      value: acountPage.userData.phone,
      from: 5,
      to: acountPage.userData.phone.length - 4,
    });
    securityLevel = acountPage.userData.securityLevel;
    securityLevelMessage = acountPage.userData.securityLevelMessage;
    emailVerified = acountPage.userData.isAccountVerified;
    kycStatus = acountPage.userData.kycStatus;
    has2fa = acountPage.userData.has2fa;
    kycDescription = acountPage.userData.kycLevelMessage;
    googleAuth = acountPage.userData.google2faEnabled;
    isVerified = acountPage.userData.isAccountVerified;
    kycLevel = acountPage.userData.kycLevel;
    profileStatus = acountPage.userData.profileStatus;
  }
  // useRenderCount('account page')
  //navigation Methods
  const gotoChangePassword = () => {
    dispatch(push(AppPages.ChangePassword));
  };
  const gotoPhoneVerification = () => {
    dispatch(push(AppPages.PhoneVerification));
  };

  const howManySent = useRef(0);
  const getNewVerificationEmail = () => {
    if (
      howManySent.current < 2 ||
      storage.read(LocalStorageKeys.CanSendNewEmail) !== false
    ) {
      toast.warn('Please Check Your Email For Verification Link');
      dispatch(getNewVerificationEmailAction());
      howManySent.current++;
    } else {
      storage.write(LocalStorageKeys.CanSendNewEmail, false);
    }
  };

  return (
    <div>
      <Helmet>
        <title>Account</title>
        <meta name='description' content='Description of AccountPage' />
      </Helmet>
      <MaxWidthWrapper>
        {isLoading === true ? <GridLoading /> : <div />}
        <>
          <BreadCrumb
            links={[
              { pageName: 'home', pageLink: AppPages.HomePage },
              {
                pageName: 'acountAndSecurity',
                pageLink: AppPages.AcountPage,
                last: true,
              },
            ]}
          />
          <CardWrapper>
            <AnimateChildren isLoading={isLoading}>
              {/* first card */}
              <Card className={`card`} elevation={0}>
                <div className='cardInfoWrapper info1'>
                  <div className='iconContainer'>
                    <AcountPageUserInfoIcon theme={Themes.LIGHT} />
                  </div>
                  <div className='infoContainer'>
                    <InfoRow
                      value={userId}
                      title={<FormattedMessage {...translate.userId} />}
                    ></InfoRow>
                    <Divider className='divider' />
                    <InfoRow
                      value={emailToShow}
                      title={<FormattedMessage {...translate.emailAddress} />}
                    ></InfoRow>

                    <Divider className='divider' />
                    <InfoRow
                      value={passwordToShow}
                      title={<FormattedMessage {...translate.password} />}
                      actionButton={{
                        buttonTitle: <FormattedMessage {...translate.change} />,
                        onClick: () => {
                          throttleFunc(() => {
                            if (emailVerified) {
                              gotoChangePassword();
                            } else {
                              toast.warn('Your email is not verified');
                            }
                          });
                        },
                      }}
                    ></InfoRow>
                  </div>
                </div>
              </Card>
              {/* second card */}
              <Card className={`card `} elevation={0}>
                <div className='cardInfoWrapper info2'>
                  <div className='iconContainer'>
                    <AcountPageUserSecurityIcon theme={Themes.LIGHT} />
                  </div>
                  <div className='infoContainer'>
                    <InfoRow
                      value={<SecurityLevelBar level={securityLevel} />}
                      title={<FormattedMessage {...translate.securityLevel} />}
                      description={securityLevelMessage}
                      titleIcon={
                        securityLevel != SecurityLevel.LOW ? (
                          <SecurityIcon level={securityLevel} />
                        ) : (
                          <img src={redShield} style={{ margin: 0 }} alt='' />
                        )
                      }
                    ></InfoRow>
                    <Divider className='divider' />
                    <InfoRow
                      value={emailToShow}
                      title={<FormattedMessage {...translate.email} />}
                      titleIcon={
                        emailVerified ? (
                          <img className='noMargin' src={checkIcon} />
                        ) : (
                          <div className='colorDot orange'></div>
                        )
                      }
                      actionButton={{
                        buttonTitle: emailVerified ? (
                          <FormattedMessage {...translate.verified} />
                        ) : (
                          <FormattedMessage {...translate.verify} />
                        ),
                        disabled: emailVerified,
                        onClick: () => {
                          throttleFunc(() => {
                            getNewVerificationEmail();
                          });
                        },
                      }}
                    ></InfoRow>

                    <Divider className='divider' />
                    <InfoRow
                      value={phoneToShow}
                      title={<FormattedMessage {...translate.phone} />}
                      titleIcon={
                        phoneToShow !== '' ? (
                          <img className='noMargin' src={checkIcon} />
                        ) : (
                          <div className='colorDot orange'></div>
                        )
                      }
                      actionButton={{
                        buttonTitle:
                          phoneToShow !== '' ? (
                            <FormattedMessage {...translate.change} />
                          ) : (
                            <FormattedMessage {...translate.verify} />
                          ),
                        onClick: () => {
                          throttleFunc(() => {
                            if (emailVerified) {
                              gotoPhoneVerification();
                            } else {
                              toast.warn('Your email is not verified');
                            }
                          });
                        },
                      }}
                    ></InfoRow>
                    <Divider className='divider' />
                    <InfoRow
                      titleIcon={
                        googleAuth ? (
                          <img className='noMargin' src={checkIcon} />
                        ) : (
                          <div className='colorDot orange'></div>
                        )
                      }
                      title={
                        <FormattedMessage {...translate.googleAuthentication} />
                      }
                      actionButton={{
                        buttonTitle: googleAuth ? (
                          <FormattedMessage {...translate.disable} />
                        ) : (
                          <FormattedMessage {...translate.enable} />
                        ),
                        // disabled: googleAuth,
                        onClick: () => {
                          throttleFunc(() => {
                            if (emailVerified) {
                              dispatch(push(AppPages.GoogleAuthentication));
                            } else {
                              toast.warn('Your email is not verified');
                            }
                          });
                        },
                      }}
                    ></InfoRow>
                    <Divider className='divider' />

                    <InfoRow
                      titleIcon={
                        profileStatus == 'confirmed' ? (
                          <img className='noMargin' src={checkIcon} />
                        ) : (
                          <div className='colorDot orange'></div>
                        )
                      }
                      title={
                        <FormattedMessage {...translate.identityVerification} />
                      }
                      description={kycDescription}
                      actionButton={{
                        buttonTitle:
                          profileStatus == 'confirmed' ? (
                            <FormattedMessage {...translate.verified} />
                          ) : (
                            <FormattedMessage {...translate.verify} />
                          ),
                        disabled: profileStatus == 'confirmed',
                        onClick: () => {
                          throttleFunc(() => {
                            if (emailVerified) {
                              dispatch(push(AppPages.UserInfo));
                            } else {
                              toast.warn('Your email is not verified');
                            }
                          });
                        },
                      }}
                    ></InfoRow>
                  </div>
                </div>
              </Card>

              {/* third card */}
              <Card className={`card`} elevation={0}>
                <div className='cardInfoWrapper info3'>
                  <div className='iconContainer'>
                    <AcountPageDepositIcon theme={Themes.LIGHT} />
                  </div>
                  <div className='infoContainer'>
                    <InfoRow
                      title={
                        <FormattedMessage
                          {...translate.withdrawAddressManagement}
                        />
                      }
                      titleIcon={<img src={bagIcon} />}
                      bigTitle={true}
                      actionButton={{
                        buttonTitle: <FormattedMessage {...translate.manage} />,
                        onClick: () => {
                          throttleFunc(() => {
                            if (emailVerified) {
                              dispatch(push(AppPages.AddressManagement));
                            } else {
                              toast.warn('Your email is not verified');
                            }
                          });
                        },
                      }}
                    ></InfoRow>
                    <Divider className='divider' />
                    <InfoRow
                      title={<FormattedMessage {...translate.loginWhiteList} />}
                      titleIcon={<img src={listIcon} />}
                      actionButton={{
                        disabled: true,
                        buttonTitle: (
                          <FormattedMessage {...translate.disabled} />
                        ),
                        onClick: () => { },
                      }}
                    ></InfoRow>

                    <Divider className='divider' />
                    <InfoRow
                      title={
                        <FormattedMessage {...translate.recentLoginHistory} />
                      }
                      titleIcon={<img src={historyIcon} />}
                      actionButton={{
                        disabled: true,
                        buttonTitle: (
                          <FormattedMessage {...translate.disabled} />
                        ),
                        onClick: () => { },
                      }}
                    ></InfoRow>
                  </div>
                </div>
              </Card>
            </AnimateChildren>
          </CardWrapper>
        </>
      </MaxWidthWrapper>
    </div>
  );
}

export default memo(AcountPage);
const CardWrapper = styled.div`
  .animm0 {
    min-width: 920px;
  }
  .cardInfoWrapper {
    display: flex;
    align-items: center;
    .iconContainer {
      width: 15vw;
      min-width: 200px;
      display: flex;
      justify-content: center;
      svg {
        min-width: 140px;
        max-width: 140px;
      }
    }
    &.info1 {
      min-height: 216px;
    }
    &.info2 {
      min-height: 354px;
    }
    &.info3 {
      min-height: 216px;
    }
  }

  @media screen and (max-height: 900px) {
    --sub: 10px;
    .cardInfoWrapper {
      &.info1 {
        min-height: calc(23vh - var(--sub));
      }
      &.info2 {
        min-height: calc(38vh - var(--sub));
      }
      &.info3 {
        min-height: calc(23vh - var(--sub));
      }
    }
  }
  @media screen and (max-height: 690px) {
    --sub: 6px;
    .cardInfoWrapper {
      &.info1 {
        min-height: 160px;
      }
      &.info2 {
        min-height: 270px;
      }
      &.info3 {
        min-height: 160px;
      }
    }
  }
  .infoContainer {
    width: 100%;
    padding-right: 2vw;
    display: flex;
    flex-direction: column;
    align-self: stretch;
    justify-content: space-evenly;
  }
  .divider {
    background: #c1c1c1;
  }
  .infoRow {
    display: flex;
    flex: 1;
    align-items: center;
    .rowItemTitle {
      flex: 2;
      align-self: center;
      min-width: 180px;
      display: flex;
      align-items: center;
      &.bigTitle {
        min-width: 270px;
      }
      img {
        margin: -1px 4px 0 4px;
      }
      span {
        color: var(--blackText);
        font-weight: 600;
        font-size: 13px;
      }
    }
    .rowItemValue {
      flex: 10;
      align-self: center;
      color: var(--blackText);
      font-weight: 600;
      font-size: 13px;
      span {
        font-size: 13px;
        color: var(--blackText);
      }
    }
    .itemDescription {
      color: var(--textGrey);
      flex: 20;
      align-self: center;
      font-size: 12px !important;
    }
  }
  .card {
    margin-bottom: 1vh;
    border-radius: 10px;
  }
  .MuiButton-outlinedPrimary {
    color: var(--textBlue);
    border: 1px solid var(--textBlue);
    &:hover {
      border: 1px solid var(--textBlue);
    }
  }
  .MuiButton-outlinedSizeSmall {
    /* padding: 4px 17px 1px; */
    /* font-size: 12px; */
    /* min-width: 85px; */
    padding: 0;
    border-radius: 7px !important;
  }
  .noMargin {
    margin: 0 !important;
  }
  .colorDot {
    width: 7px;
    height: 7px;

    margin: 8px;
    border-radius: 14px;
    &.orange {
      background: var(--orange);
    }
  }
  .rowActionButton {
    .MuiButtonBase-root {
      max-width: 64px !important;
      min-height: 24px !important;
      max-height: 24px !important;
    }
    span {
      font-weight: 600;
      font-size: 12px !important;
    }
  }
`;
