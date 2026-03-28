import React, { useState } from 'react';
import { Tooltip, Button } from '@material-ui/core';
import { FormattedMessage } from 'react-intl';
import translate from '../messages';
import { QrCode } from '../types';
import { CopyToClipboard } from 'utils/formatters';
import styled from 'styles/styled-components';

export default function CodeButton (props: { qrCode: QrCode }) {
  const [IsTooltipOpen, setIsTooltipOpen] = useState(false);

  const handleCopyClick = () => {
    setIsTooltipOpen(true);
    setTimeout(() => {
      setIsTooltipOpen(false);
    }, 1500);
    CopyToClipboard(props.qrCode.code || '');
  };
  return (
    <Wrapper>
      <Tooltip
        PopperProps={{
          disablePortal: true,
        }}
        placement='top'
        open={IsTooltipOpen}
        disableFocusListener
        disableHoverListener
        disableTouchListener
        title={<FormattedMessage {...translate.codeCopiedToClipboard} />}
      >
        <Button onClick={handleCopyClick} className='code'>
          {props.qrCode.code}
        </Button>
      </Tooltip>
    </Wrapper>
  );
}
const Wrapper = styled.div`
  span {
    color: White !important;
  }
  .code {
    background: var(--greyBackground);
    span {
      font-size: 15px;
      font-weight: 600;
      color: var(--textBlue) !important;
    }
  }
`;
