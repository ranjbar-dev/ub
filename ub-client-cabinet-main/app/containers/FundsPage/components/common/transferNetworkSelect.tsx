import { Button, Tooltip, Zoom } from '@material-ui/core';
import { IOtherNetworks } from 'containers/FundsPage/types';
import InfoIcon from 'images/info';
import React, { useState } from 'react';
import { FormattedMessage } from 'react-intl';
import styled from 'styles/styled-components';
import translate from '../../messages';

interface Props {
  onNetworkSelect: (network: IOtherNetworks) => void;
  otherNetworks: IOtherNetworks[];
  type: 'withdraw' | 'deposit';
  completedNetworkName?: string;
  mainNetwork?: string;
}

const TransferNetworkSelect = (props: Props) => {
  const [selectedIndex, setselectedIndex] = useState<number>(0);
  const [IsTooltipOpen, setIsTooltipOpen] = useState<boolean>(false);
  const {
    onNetworkSelect,
    otherNetworks,
    type,
    completedNetworkName,
    mainNetwork,
  } = props;

  /*
  
 address: "TQJe5ski1DGGTA4KA5bWkTDLjACVquaTac"
code: "TRX"
completedNetworkName: "Tron(TRX) trc20"
supportsDeposit: true
supportsWithdraw: true
 
  */
  //@ts-ignore
  const buttons = otherNetworks.map(item => {
    if (
      (type === 'deposit' && item.supportsDeposit === true) ||
      (type === 'withdraw' && item.supportsWithdraw === true)
    ) {
      const small = item.completedNetworkName.split(' ')[
        item.completedNetworkName.split(' ').length - 1
      ];
      return {
        //@ts-ignore
        name: item.completedNetworkName.replace(small, ''),
        small,
        value: item.code,
        address: item.address,
      };
    }
  });
  if (buttons.length > 0) {
    //  @ts-ignore
    const small = completedNetworkName.split(' ')[
      //  @ts-ignore
      completedNetworkName.split(' ').length - 1
    ];
    buttons.unshift({
      //  @ts-ignore
      name: completedNetworkName?.replace(small, ''),
      //  @ts-ignore
      small: completedNetworkName?.split(' ')[
        completedNetworkName.split(' ').length - 1
      ],
      //  @ts-ignore
      value: mainNetwork,
    });
  }
  const handleClick = (index: number) => {
    //@ts-ignore
    onNetworkSelect(buttons[index]);
    setselectedIndex(index);
  };
  const handleOpenTooltip = () => {
    setIsTooltipOpen(true);
  };

  const handleCloseTooltip = () => {
    setIsTooltipOpen(false);
  };

  return (
    <Wrapper>
      <Title>
        <FormattedMessage {...translate.transferNetwork} />

        <Tooltip
          PopperProps={{
            disablePortal: true,
          }}
          arrow
          TransitionComponent={Zoom}
          placement='top'
          //  open
          open={IsTooltipOpen}
          disableFocusListener
          disableHoverListener
          disableTouchListener
          className='infoTooltip'
          title={'select the transfer network'}
        >
          <IconWrapper
            onMouseEnter={handleOpenTooltip}
            onMouseLeave={handleCloseTooltip}
          >
            <InfoIcon />
          </IconWrapper>
        </Tooltip>
      </Title>
      <ButtonsWrapper>
        {buttons.map((item, i) => {
          return (
            <Button
              key={i}
              onClick={() => {
                handleClick(i);
              }}
              variant={'outlined'}
              size='small'
              color='primary'
              className={`${i === buttons.length - 1 ? 'last' : ''} ${
                i === selectedIndex ? 'active' : ''
              }`}
            >
              {i === selectedIndex && <Indicator />}
              {/*@ts-ignore*/}
              {item.name}
              {/*@ts-ignore*/}
              <Small>{item.small.toUpperCase()}</Small>
            </Button>
          );
        })}
      </ButtonsWrapper>
    </Wrapper>
  );
};
const Indicator = styled.div`
  width: 7px;
  height: 7px;
  background: var(--textBlue) !important;
  margin-left: -14px;
  margin-right: 7px;
  border-radius: 10px;
`;
const IconWrapper = styled.div`
  cursor: pointer;
  margin-top: -4px;
`;
const Small = styled.span`
  font-size: 9px !important;
  padding-left: 5px;
  margin-top: 1px;
`;
const ButtonsWrapper = styled.div`
  display: flex;
  .MuiButton-label {
    font-size: 11px !important;
  }

  button {
    flex: 1;
    margin-right: 3%;
    border: 1px solid var(--textGrey) !important;

    &.last {
      margin-right: 0%;
    }
    span {
      color: var(--textGrey) !important;
    }
    &.active {
      border: 1px solid var(--textBlue) !important;
      span {
        color: var(--textBlue) !important;
      }
    }
  }
`;
const Title = styled.div`
  height: 24px;
  display: flex;
  align-items: center;
  span {
    font-size: 9px;
    color: var(--textGrey);
    margin-right: 8px;
  }
  .MuiTooltip-popper {
    top: 15px !important;
    .MuiTooltip-arrow {
      color: rgba(97, 97, 97, 0.9) !important;
    }
  }
`;

const Wrapper = styled.div``;
export default TransferNetworkSelect;
