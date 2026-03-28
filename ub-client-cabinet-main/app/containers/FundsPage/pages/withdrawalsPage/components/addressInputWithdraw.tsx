import React, { useState, useEffect } from 'react';
import { TextField } from '@material-ui/core';
import { injectIntl } from 'react-intl';
import PopupModal from 'components/materialModal/modal';

import { IconButton } from '@material-ui/core';
import styled from 'styles/styled-components';
import { Button } from '@material-ui/core';
import { useDispatch } from 'react-redux';
import { addNewAddressAction } from 'containers/FundsPage/actions';
import IsLoadingWithText from 'components/isLoadingWithText/isLoadingWithText';
import { Subscriber, MessageNames } from 'services/message_service';
import { AddIcon } from 'images/themedIcons/expandMore';
import { toast } from 'components/Customized/react-toastify';

const AddressInputWithdraw = (props: {
  onChange: Function;
  formerWithdrawAddresses: any[];
  intl: any;
  selectedAddress: string;
  code: string;
  network: string;
  onAddressExists: Function;
}) => {
  const dispatch = useDispatch();
  const [IsAddressInputOpen, setIsAddressInputOpen] = useState(false);
  const [IsAddButtonVisible, setIsAddButtonVisible] = useState(false);
  const [AddressLabel, setAddressLabel] = useState('');
  const [IsLoading, setIsLoading] = useState(false);
  const intl = props.intl;
  const onCreateClick = () => {
    if (!AddressLabel) {
      toast.warn('please enter a label for address');
      return;
    }
    const dataToSend = {
      address: props.selectedAddress,
      label: AddressLabel,
      code: props.code,
      ...(props.network && { network: props.network }),
    };
    dispatch(addNewAddressAction(dataToSend));
  };
  useEffect(() => {
    const Subscription = Subscriber.subscribe((message: any) => {
      if (message.name === MessageNames.SET_POPUP_LOADING) {
        setIsLoading(message.payload);
      } else if (message.name === MessageNames.CLOSE_MODAL) {
        setIsAddressInputOpen(false);
      }
      if (message.name === MessageNames.ADDITIONAL_ACTION) {
        if (message.payload.title === 'removeAddButton') {
          setIsAddButtonVisible(false);
        }
      }
    });
    return () => {
      Subscription.unsubscribe();
    };
  }, []);
  const handleAddressChange = e => {
    const addressArray: any = [];
    props.formerWithdrawAddresses.forEach(item => {
      if (item.address === e.target.value) {
        addressArray.push(item);
      }
    });
    if (addressArray.length > 0) {
      if (IsAddButtonVisible === true) {
        setIsAddButtonVisible(false);
      }
      props.onAddressExists(addressArray);
      return;
    }
    if (e.target.value.length > 15) {
      if (IsAddButtonVisible === false) {
        setIsAddButtonVisible(true);
      }
    } else {
      if (IsAddButtonVisible === true) {
        setIsAddButtonVisible(false);
      }
    }
    props.onChange(e.target.value);
  };

  return (
    <Wrapper>
      <PopupModal
        isOpen={IsAddressInputOpen}
        onClose={() => {
          setIsAddressInputOpen(false);
        }}
      >
        <div className='alertWrapper addAddress'>
          <TextField
            variant='outlined'
            margin='dense'
            fullWidth
            onChange={({ target }) => {
              setAddressLabel(target.value);
            }}
            label={intl.formatMessage({
              id: 'containers.FundsPage.label',
              defaultMessage: 'ET.Label',
            })}
          />
          <p>
            {intl.formatMessage({
              id: 'containers.FundsPage.Thelabelmustnotexceed20characters',
              defaultMessage: 'ET.Label',
            })}
          </p>
          <Button onClick={onCreateClick} color='primary' variant='contained'>
            <IsLoadingWithText
              text={intl.formatMessage({
                id: 'containers.FundsPage.create',
                defaultMessage: 'ET.Label',
              })}
              isLoading={IsLoading}
            />
          </Button>
        </div>
      </PopupModal>
      <TextField
        placeholder={
          props.formerWithdrawAddresses.length > 0
            ? intl.formatMessage({
                id: 'containers.FundsPage.TypeorSelectwithdrawaladdress',
                defaultMessage: 'Type or Select withdrawal address',
              })
            : intl.formatMessage({
                id: 'containers.FundsPage.Typewithdrawaladdress',
                defaultMessage: 'Type withdrawal address',
              })
        }
        fullWidth
        value={props.selectedAddress}
        onChange={handleAddressChange}
        className='textField'
      />
      <div className='addButton'>
        {IsAddButtonVisible === true && (
          <IconButton
            onClick={() => {
              setIsAddressInputOpen(true);
            }}
            color='primary'
            aria-label='upload picture'
            component='span'
            size='small'
          >
            <AddIcon />
          </IconButton>
        )}
      </div>
    </Wrapper>
  );
};
export default injectIntl(AddressInputWithdraw);
const Wrapper = styled.div`
  .addButton {
    position: absolute;
    right: 0;
    top: 7px;
  }
  input {
    font-weight: 600;
  }
`;
