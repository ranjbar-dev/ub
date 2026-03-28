/*
 *
 * DocumentVerificationPage
 *
 */

import React, { memo, useEffect, useState } from 'react';
import { Helmet } from 'react-helmet';
import { useSelector, useDispatch } from 'react-redux';
import { createStructuredSelector } from 'reselect';

import { useInjectSaga } from 'utils/injectSaga';
import { useInjectReducer } from 'utils/injectReducer';
import makeSelectDocumentVerificationPage, {
  makeSelectUserProfileData,
  makeSelectIsLoadingUserProfileData,
} from './selectors';
import AnimateChildren from 'components/AnimateChildren';

import reducer from './reducer';
import saga from './saga';
import { MaxWidthWrapper } from 'components/wrappers/maxWidthWrapper';
import styled from 'styles/styled-components';
import { Card } from '@material-ui/core';
import ProofOfIdentity from './components/proofs';
import { getUserProfileAction } from './actions';
import { AppPages } from 'containers/App/constants';
import BreadCrumb from 'components/BreadCrumb';
import PopupModal from 'components/materialModal/modal';
import { Subscriber, MessageNames } from 'services/message_service';

const stateSelector = createStructuredSelector({
  documentVerificationPage: makeSelectDocumentVerificationPage(),
  userProfileData: makeSelectUserProfileData(),
  isLoading: makeSelectIsLoadingUserProfileData(),
});
interface Props {}

function DocumentVerificationPage (props: Props) {
  // Warning: Add your key to RootState in types/index.d.ts file
  useInjectReducer({ key: 'documentVerificationPage', reducer: reducer });
  useInjectSaga({ key: 'documentVerificationPage', saga: saga });

  const { userProfileData, isLoading } = useSelector(stateSelector);
  const [IsAlertOpen, setIsAlertOpen] = useState(false);
  const [AlertContent, setAlertContent] = useState(<div></div>);
  const dispatch = useDispatch();

  useEffect(() => {
    dispatch(getUserProfileAction({}));
    const Subscription = Subscriber.subscribe((message: any) => {
      if (message.name === MessageNames.OPEN_ALERT) {
        setAlertContent(message.payload);
        setIsAlertOpen(true);
      }
    });
    return () => {
      Subscription.unsubscribe();
    };
  }, []);

  return (
    <>
      <Helmet>
        <title>Document Verification</title>
        <meta
          name='description'
          content='Description of DocumentVerificationPage'
        />
      </Helmet>
      <PopupModal
        isOpen={IsAlertOpen}
        onClose={() => {
          setIsAlertOpen(false);
        }}
      >
        {AlertContent}
      </PopupModal>
      <MaxWidthWrapper>
        <BreadCrumb
          links={[
            { pageName: 'home', pageLink: AppPages.HomePage },
            {
              pageName: 'acountAndSecurity',
              pageLink: AppPages.AcountPage,
            },
            {
              pageName: 'DocumentVerification',
              pageLink: AppPages.DocumentVerification,
              last: true,
            },
          ]}
        />
        <MainWrapper>
          <AnimateChildren isLoading={isLoading}>
            <HalfCard>
              <ProofOfIdentity
                userProfileData={userProfileData}
                documentType={'identity'}
              />
            </HalfCard>
            <HalfCard>
              {/* <ProofOfResidence userProfileData={userProfileData} /> */}
              <ProofOfIdentity
                userProfileData={userProfileData}
                documentType={'address'}
              />
            </HalfCard>
          </AnimateChildren>
        </MainWrapper>
      </MaxWidthWrapper>
    </>
  );
}
export default memo(DocumentVerificationPage);
const MainWrapper = styled.div`
  width: calc(100% + 12px);
  display: flex;
  overflow: auto;
  justify-content: space-between;
  .animm0 {
    flex: 1;
    margin-right: 12px;
  }
`;
const HalfCard = styled(Card)`
  height: calc(98vh - 115px);
  min-height: 685px;
  min-width: 500px;
  max-height: var(--maxHeight);
  border-radius: 10px !important;
  box-shadow: none !important;
  .MuiSelect-outlined.MuiSelect-outlined {
    padding-top: 15px;
    font-size: 13px;
  }
  .MuiButton-outlined.Mui-disabled {
    border: 1px solid var(--darkGrey) !important;
    color: var(--textGrey) !important;
    background: var(--oddRows) !important;
    font-size: 13px !important;
    border-radius: 7px !important;
  }
`;
