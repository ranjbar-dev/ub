import React, { useState, useEffect } from 'react';
import styled from 'styles/styled-components';
import { TextField } from '@material-ui/core';
import { LocalStorageKeys } from 'services/constants';
import { Country } from 'containers/PhoneVerificationPage/constants';
import ExpandMore from 'components/icons/expandMore';
import Autocomplete from '@material-ui/lab/Autocomplete';
let initial;
export default function SelectRegion (props: {
  initialCountryId: number;
  onCountrySelect: Function;
}) {
  const [SelectedCountry, setSelectedCountry]: [any, any] = useState({
    id: 799,
    name: 'USA',
    fullName: 'United States of America',
    code: '1',
    image: 'https://app.unitedbit.com/assets/images/country-logo/us.png',
  });
  const countries: Country[] = localStorage[LocalStorageKeys.COUNTRIES]
    ? JSON.parse(localStorage[LocalStorageKeys.COUNTRIES])
    : [];
  useEffect(() => {
    for (let i = 0; i < countries.length; i++) {
      if (countries[i].id === props.initialCountryId) {
        setSelectedCountry(countries[i]);
      }
    }
    return () => {};
  }, []);
  const handleChange = e => {
    props.onCountrySelect(e);
    return;
  };
  return (
    <Wrapper>
      <Autocomplete
        className='autoComplete'
        options={countries as Country[]}
        // defaultValue={SelectedCountry}
        value={SelectedCountry}
        popupIcon={<ExpandMore />}
        noOptionsText='no results'
        autoHighlight
        getOptionLabel={option => {
          return option.fullName;
        }}
        onChange={(e, value) => {
          if (value && value.code) {
            setSelectedCountry(value);
            handleChange(value);
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
    </Wrapper>
  );
}
const Wrapper = styled.div`
  margin-top: 24px;
  position: relative;
  .selectedFlag {
    position: absolute;
    top: 6px;
    left: 10px;
    img {
      width: 25px;
      border-radius: 3px;
    }
  }
`;
