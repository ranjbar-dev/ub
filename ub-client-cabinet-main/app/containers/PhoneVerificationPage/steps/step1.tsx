import React, { useState, useEffect } from 'react';
import { MainIconWrapper } from 'components/wrappers/mainIconWrapper';
import { CenterInputsWrapper } from 'components/wrappers/centerInputsWrapper';
import { CenterButtonsWrapper } from 'components/wrappers/centerButtonsWrapper';

import { Country } from '../constants';

import translate from './messages';
import { FormattedMessage } from 'react-intl';
import Autocomplete from '@material-ui/lab/Autocomplete';
import PhoneVStep1Icon from 'images/themedIcons/phoneVStep1Icon';
import { Buttons } from 'containers/App/constants';
import InputWithValidator from 'components/inputWithValidator';
import { PhoneValidator } from '../validators/phoneValidator';
import StreamLoadingButton from 'components/streamLoadingButton';
import { TextField } from '@material-ui/core';
import { Button } from '@material-ui/core';
import ExpandMore from 'components/icons/expandMore';
let fields = {
  pre: { isValid: false, value: '' },
  phone: { isValid: false, value: '' },
};

export default function Step1 (props: {
  onCancel: Function;
  onSubmit: Function;
  countries: Country[];
  submitIsLoading: boolean;
}) {
  const countries = props.countries;
  const [CanSubmit, setCanSubmit] = useState(false);
  const [SelectedCountry, setSelectedCountry]: [any, any] = useState({});
  const isFieldValid = (properties: {
    fieldName: string;
    isValid: boolean;
    value: string;
  }) => {
    fields[properties.fieldName].isValid = properties.isValid;
    fields[properties.fieldName].value = properties.value;
    if (fields.phone.isValid === true) {
      setCanSubmit(true);
    } else {
      setCanSubmit(false);
    }
  };

  useEffect(() => {
    return () => {
      fields = {
        pre: { isValid: false, value: '' },
        phone: { isValid: false, value: '' },
      };
    };
  }, []);

  const handleSubmit = () => {
    props.onSubmit({
      country: SelectedCountry,
      phoneNumber: fields.phone.value,
    });
  };
  const handleCancelButton = () => {
    props.onCancel();
  };

  return (
    <>
      <MainIconWrapper className='fl fl9 minimized'>
        <PhoneVStep1Icon />
      </MainIconWrapper>
      <CenterInputsWrapper
        className='noPadd pt3'
        style={{ flex: 6, position: 'relative' }}
      >

        <Autocomplete
          className={`autoComplete`}
          options={countries as Country[]}
          disablePortal
          popupIcon={<ExpandMore />}
          noOptionsText={'no results'}
          autoHighlight
          getOptionLabel={option => {
            return option.fullName;
          }}
          onChange={(e, value) => {
            if (value && value.code) {
              setSelectedCountry(value);
            } else {
              setSelectedCountry({});
            }
          }}
          renderOption={option => (
            <div className='countryContainer'>
              <span className='flag'>
                <img src={option.image} alt='' />
              </span>
              <span>
                {option.fullName.length < 30
                  ? option.fullName
                  : option.fullName.slice(0, 28) + '...'}
              </span>
            </div>
          )}
          renderInput={params => (
            <TextField
              {...params}
              className={`countryInput ${SelectedCountry.id && 'filled'}`}
              placeholder='Select Country'
              variant='outlined'
              inputProps={{
                ...params.inputProps,
                autoComplete: 'new-password',
              }}
            />
          )}
        />
        <div className='selectedFlag'>
          {SelectedCountry && (
            <img className='flag' src={SelectedCountry.image} />
          )}
        </div>
        <InputWithValidator
          inputType='number'
          label={<FormattedMessage {...translate.enterPhoneNumber} />}
          startComponent={
            SelectedCountry.code ? (
              <div className='preNumber'>+{SelectedCountry.code}</div>
            ) : null
          }
          onChange={(phone: string) => {
            isFieldValid({
              fieldName: 'phone',
              isValid: PhoneValidator({
                uniqueInputId: 'phone',
                value: phone,
              }),
              value: phone,
            });
          }}
          uniqueName='phone'
          throttleTime={500}
        />
      </CenterInputsWrapper>
      <CenterButtonsWrapper>
        <StreamLoadingButton
          disabled={!CanSubmit || !SelectedCountry.id}
          onClick={handleSubmit}
          className='ubButton nmt'
          variant='contained'
          color='primary'
          text={<FormattedMessage {...translate.getSMS} />}
        />
        <Button className={Buttons.CancelButton} onClick={handleCancelButton}>
          <FormattedMessage {...translate.cancel} />
        </Button>
      </CenterButtonsWrapper>
    </>
  );
}
