import { TextField } from '@material-ui/core';
import { Autocomplete } from '@material-ui/lab';
import { Country } from 'locales/types';
import React, { memo, useMemo, useState } from 'react';
import { LocalStorageKeys } from 'services/constants';
import styled from 'styled-components/macro';

interface Props {
  onChange: (country: string) => void;
  onClear: () => void;
}

/**
 * Country autocomplete filter input for AG Grid column headers.
 * Reads the countries list from localStorage and emits country id on selection.
 *
 * @example
 * ```tsx
 * <CountryFilter onChange={(id) => applyFilter(id)} onClear={() => clearFilter()} />
 * ```
 */
function CountryFilter(props: Props) {
  const { onChange } = props;
  const countries = useMemo(
    () => JSON.parse(localStorage[LocalStorageKeys.COUNTRIES]),
    [],
  );
  const [SelectedCountry, setSelectedCountry] = useState<Country | undefined>();
  const handleChange = (e: Country) => {
    onChange(e.id + '');
    return;
  };
  return (
    <Wrapper>
      <Autocomplete
        className="autoComplete"
        options={countries as Country[]}
        //popupIcon={<ExpandMore />}
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
    background: white !important;
    .MuiButtonBase-root {
      background: white !important;
    }
  }
`;

export default memo(CountryFilter);
