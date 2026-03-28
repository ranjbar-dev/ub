import React from 'react';
import crossIcon from 'images/crossIcon.svg';

import styled from 'styles/styled-components';
import { TransitionProps } from '@material-ui/core/transitions/transition';
import { Zoom, Dialog, IconButton } from '@material-ui/core';

const Transition = React.forwardRef<unknown, TransitionProps>(
  function Transition (props, ref) {
    // @ts-ignore
    return <Zoom ref={ref} {...props} {...{ timeout: 300 }} />;
  },
);
export default function PopupModal (props: {
  isOpen: boolean;
  children: any;
  onClose: Function;
}) {
  return (
    <Dialog
      TransitionComponent={Transition}
      fullScreen={false}
      open={props.isOpen}
      onClose={() => {
        props.onClose();
      }}
    >
      <CloseWrapper>
        <IconButton
          onClick={() => {
            props.onClose();
          }}
          size='small'
        >
          <img src={crossIcon} alt='' />
        </IconButton>
      </CloseWrapper>
      <ChildsWrapper>{props.children}</ChildsWrapper>
    </Dialog>
  );
}
const CloseWrapper = styled.div`
  position: absolute;
  right: 5px;
  top: 5px;
  img {
    width: 35px;
  }
`;
const ChildsWrapper = styled.div`
  .alertWrapper {
    padding: 72px 102px 24px 102px;
    min-height: 220px;
    display: flex;
    flex-direction: column;
    align-items: center;
    label {
      color: var(--blackText);
    }
    p {
      margin-top: 0;
      color: var(--blackText);
    }
  }
  .addAddress {
    .MuiInputLabel-outlined.MuiInputLabel-marginDense {
      margin-top: -0.5px !important;
    }
    .loadingCircle {
      top: 8px !important;
    }
    .MuiInputBase-root legend {
      zoom: 0.72;
      margin-left: 5px;
    }
  }
`;
