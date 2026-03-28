import React from 'react';
import { FormattedMessage } from 'react-intl';
import translate from '../../../messages';

import styled from 'styles/styled-components';

const Description = (props: { value }) => (
  <div className='desc'>
    <div className='blueDot'></div>
    {props.value}
  </div>
);
export default function Descriptions () {
  return (
    <Wrapper>
      <Description
        value={
          <FormattedMessage
            {...translate.AcceptabledocumentsNationalIDDriverLicensePassport}
          />
        }
      />
      <Description
        value={
          <FormattedMessage
            {...translate.Thedocumentmustbeissuedwithinthelastsixmonths}
          />
        }
      />
      <Description
        value={
          <FormattedMessage
            {...translate.Theresolutionoftheuploadeddocumentsmustbeatleast300dpi}
          />
        }
      />
      <Description
        value={
          <FormattedMessage
            {...translate.AcceptablefileformatsGIFJPEGJPGPNGBMPorPDF}
          />
        }
      />
      <Description
        value={
          <FormattedMessage
            {...translate.Electronicdocumentsarenotacceptable}
          />
        }
      />
    </Wrapper>
  );
}
const Wrapper = styled.div`
  width: 100%;
  padding: 0 51px 0 41px;
  .desc {
    display: flex;
    align-items: center;
    margin: 10px 0;
    .blueDot {
      width: 6px;
      height: 6px;
      background: var(--textBlue);
      border-radius: 100%;
      margin: 0 7px;
      margin-top: 1px;
    }
    span {
      color: var(--textGrey);
      font-weight: 600;
      font-size: 13px;
    }
  }
`;
