import React, { memo } from 'react';
import styled from 'styles/styled-components';
import { IUserProfileImage } from 'containers/DocumentVerificationPage/types';
import VerifiedResidence from 'images/themedIcons/verifiedResidence';
import VerifiedIdentity from 'images/themedIcons/verifeidIdentity';
import { FormattedMessage } from 'react-intl';
import translate from '../../messages';
interface Props {
  uploadedImage: IUserProfileImage;
}

function Verified (props: Props) {
  const { uploadedImage } = props;

  return (
    <Wrapper>
      <div className='withIconWrapper confirmedWrapper'>
        {uploadedImage.type === 'address' ? (
          <VerifiedResidence />
        ) : (
          <VerifiedIdentity />
        )}

        <FormattedMessage {...translate.YourDocumentHasBeenVerified} />
      </div>
    </Wrapper>
  );
}

export default memo(Verified);
const Wrapper = styled.div`
  height: 100%;
  .confirmedWrapper {
    height: 86%;
    justify-content: center;
    align-items: center;
    text-align: center;
  }
  span {
    font-size: 14px !important;
    padding: 0px 10px;
  }
`;
