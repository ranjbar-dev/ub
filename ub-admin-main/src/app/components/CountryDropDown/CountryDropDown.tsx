import { TextField } from '@material-ui/core';
import { Autocomplete } from '@material-ui/lab';
import { Country } from 'locales/types';
import React, { memo, useMemo, useState } from 'react';
import { LocalStorageKeys } from 'services/constants';
import styled from 'styled-components/macro';

interface Props {
  onChange: (country: number) => void;
  onClear: () => void;
  initialCountryId?: number;
}

/**
 * Country autocomplete dropdown for form fields.
 * Reads the countries list from localStorage and lets the user search by name.
 *
 * @example
 * ```tsx
 * <CountryDropDown
 *   initialCountryId={user.countryId}
 *   onChange={(id) => setCountry(id)}
 *   onClear={() => setCountry(undefined)}
 * />
 * ```
 */
function CountryDropDown(props: Props) {
  const { onChange, initialCountryId } = props;
  const countries: Country[] = useMemo(
    () => JSON.parse(localStorage[LocalStorageKeys.COUNTRIES]),
    [],
  );
  const [SelectedCountry, setSelectedCountry] = useState<Country | undefined>(
    initialCountryId ? countries.find(item => item.id === initialCountryId) : undefined,
  );
  const handleChange = (e: Country) => {
    onChange(e.id);
    return;
  };
  return (
    <Wrapper className="CountryDropDownWrapper">
      <Autocomplete
        className="autoComplete"
        options={countries as Country[]}
        disablePortal
        //popupIcon={<ExpandMore />}
        defaultValue={SelectedCountry ?? null}
        autoHighlight
        getOptionLabel={option => {
          return option.fullName;
        }}
        onChange={(e, value) => {
          if (value && value.code) {
            setSelectedCountry(value);
            handleChange(value);
          } else {
            props.onClear();
            setSelectedCountry(undefined);
          }
        }}
        renderOption={option => (
          <div className="countryContainer">
            <span className="flag">
              <img src={option.image} alt="" />
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
            placeholder="Country"
            variant="outlined"
            inputProps={{
              ...params.inputProps,
              autoComplete: 'new-password',
            }}
          />
        )}
      />
    </Wrapper>
  );
}
const Wrapper = styled.div`
  .MuiInputBase-adornedEnd {
    background: #f9f9f9 !important;
    .MuiButtonBase-root {
      background: #f9f9f9 !important;
    }
  }
`;

export default memo(CountryDropDown);
