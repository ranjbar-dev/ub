import React, { useRef, useState } from 'react';

import { FormattedMessage } from 'react-intl';
import { Currency } from '../types';
import translate from '../messages';

import { Select, TextField, ListItemText, MenuItem } from '@material-ui/core';

import PulsingButton from 'components/pulsingButton';
import styled from 'styles/styled-components';
import { toast } from 'components/Customized/react-toastify/cjs/react-toastify';
import ExpandMore from 'images/themedIcons/expandMore';

export default function CreateWrapper(props: {
  currencyList: Currency[];
  onCreateClick: Function;
}) {
  const currencyList = props.currencyList;
  const [SelectedCurrencyIndex, setSelectedCurrencyIndex] = useState(-1);
  const [SelectedNetworkIndex, setSelectedNetworkIndex] = useState(-1);
  const network = useRef<any>(null);
  const [Address, setAddress] = useState('');
  const [Label, setLabel] = useState('');
  const onCreateClick = () => {
    if (SelectedCurrencyIndex === -1) {
      toast.warn('Please select currency');
      return;
    }
    if (otherNetworks() && otherNetworks().length > 0) {
      if (SelectedNetworkIndex === -1 || network.current === null) {
        toast.warn('Please select a network');
        return;
      }
    }
    if (Address === '') {
      toast.warn('Please enter the address');
      return;
    }
    if (Label === '') {
      toast.warn('Please enter a label');
      return;
    }

    const dataToSend = {
      address: Address,
      label: Label,
      ...(network.current && { network: network.current }),
      code: currencyList[SelectedCurrencyIndex].code,
    };

    props.onCreateClick(dataToSend);
  };
  const handleCurrencyChange = (e: any) => {
    setSelectedCurrencyIndex(e.target.value);

    network.current = null;
    setSelectedNetworkIndex(-1);
  };

  let allNetworks: any[] = [];
  const networks =
    SelectedCurrencyIndex !== -1
      ? currencyList[SelectedCurrencyIndex].otherBlockChainNetworks
      : [];
  if (networks && networks.length > 0) {
    //@ts-ignore
    allNetworks = [...networks];
    allNetworks.unshift({
      //@ts-ignore
      code: currencyList[SelectedCurrencyIndex].mainNetwork,
    });
  }
  const handleNetworkChange = (e: any) => {
    setSelectedNetworkIndex(e.target.value);
    network.current =
      //@ts-ignore
      allNetworks[
        e.target.value
        //@ts-ignore
      ].code ?? null;
  };
  const otherNetworks = () => {
    return allNetworks;
  };

  return (
    <Wrapper>
      <div className="coin">
        <Select
          className="select"
          margin="dense"
          fullWidth
          IconComponent={ExpandMore}
          MenuProps={{
            getContentAnchorEl: null,
            anchorOrigin: {
              vertical: 'bottom',
              horizontal: 'left',
            },
          }}
          variant="outlined"
          value={SelectedCurrencyIndex}
          onChange={handleCurrencyChange}
        >
          <MenuItem value={-1}>
            <FormattedMessage {...translate.coin} />
          </MenuItem>
          {currencyList.map((item: Currency, index: number) => {
            return (
              <MenuItem key={'currencyList' + index} value={index}>
                <ListItemText className="addressCoin" primary={item.name} />
              </MenuItem>
            );
          })}
        </Select>
      </div>
      {otherNetworks() && otherNetworks().length > 0 && (
        <div className="coin" style={{ marginLeft: '10px' }}>
          <Select
            className="select"
            margin="dense"
            fullWidth
            IconComponent={ExpandMore}
            MenuProps={{
              getContentAnchorEl: null,
              anchorOrigin: {
                vertical: 'bottom',
                horizontal: 'left',
              },
            }}
            variant="outlined"
            value={SelectedNetworkIndex}
            onChange={handleNetworkChange}
          >
            <MenuItem value={-1}>
              <FormattedMessage {...translate.network} />
            </MenuItem>
            {otherNetworks().map((item: any, index: number) => {
              return (
                <MenuItem key={'otherNetworks' + index} value={index}>
                  <ListItemText
                    className="addressCoin"
                    primary={item.name ?? item.code}
                  />
                </MenuItem>
              );
            })}
          </Select>
        </div>
      )}
      <div className="label1">
        <TextField
          margin="dense"
          onChange={(e) => setAddress(e.target.value)}
          fullWidth
          variant="outlined"
          label={<FormattedMessage {...translate.address} />}
        />
      </div>
      <div className="label2">
        <TextField
          margin="dense"
          onChange={(e) => setLabel(e.target.value)}
          fullWidth
          variant="outlined"
          label={<FormattedMessage {...translate.label} />}
        />
      </div>
      <div className="buttonWrapper">
        <PulsingButton
          title={<FormattedMessage {...translate.create} />}
          onClick={onCreateClick}
        />
      </div>
    </Wrapper>
  );
}
const Wrapper = styled.div`
  display: flex;
  align-items: center;
  height: 75%;
  .coin {
    flex: 4;
  }
  .label1 {
    flex: 6;
    padding: 0 10px;
  }
  .label2 {
    flex: 4;
    padding: 0 10px 0 0;
  }
  .buttonWrapper {
    flex: 5;
    margin-top: -5px;
  }
  .MuiInputBase-root {
    max-height: 32px;
    &.select {
      margin-top: -2px;
      margin-bottom: -5px;
    }
  }

  .MuiInputLabel-outlined.MuiInputLabel-marginDense {
    transform: translate(14px, 7px) scale(1);
  }
  .MuiInputLabel-outlined.MuiInputLabel-shrink {
    transform: translate(14px, -6px) scale(0.75);
  }
  input {
    font-size: 12px;
  }
`;
