import React, { useMemo, useState } from 'react';
import {
  Select,
  ListItemAvatar,
  ListItemText,
  ListItem,
  TextField,
} from '@material-ui/core';
import ExpandMore from 'images/themedIcons/expandMore';
import { FixedSizeList, ListChildComponentProps } from 'react-window';
import styled from 'styles/styled-components';
let list;
export default function FilterableSelect (props: {
  list: any[];
  fieldName: string;
  onSelect: Function;
  hasImage?: boolean;
}) {
  const initialList = props.list;
  list = props.list;
  const fieldName = props.fieldName;
  const [IsSelectCountriesOpen, setIsSelectCountriesOpen] = React.useState(
    false,
  );
  const [ItemsList, setItemsList] = useState(initialList);
  const handleChange = e => {
    props.onSelect(e);
    setIsSelectCountriesOpen(false);
  };
  const filterList = e => {
    console.log(e);
  };
  const renderItems = (propertiess: ListChildComponentProps) => {
    const { index, style } = propertiess;
    return (
      <ListItem
        button
        onClick={() => {
          handleChange(index);
        }}
        style={style}
        key={index + 'countryList'}
      >
        {props.hasImage === true && (
          <ListItemAvatar>
            <div className='countryContainer'>
              <span>
                <img src={ItemsList[index].image} alt='' />{' '}
              </span>
            </div>
          </ListItemAvatar>
        )}
        <ListItemText primary={ItemsList[index][fieldName]} />
      </ListItem>
    );
  };
  return useMemo(
    () => (
      <Wrapper>
        <TextField
          className='inputField'
          variant='outlined'
          onChange={({ target }) => {
            filterList(target.value);
          }}
          //   onFocus={() => setIsSelectCountriesOpen(true)}
          //   onBlur={() => setIsSelectCountriesOpen(false)}
          fullWidth
          placeholder='select country'
        ></TextField>
        <Select
          fullWidth
          variant='outlined'
          className='maxSelect'
          open={IsSelectCountriesOpen}
          onOpen={() => setIsSelectCountriesOpen(true)}
          onClose={() => setIsSelectCountriesOpen(false)}
          IconComponent={ExpandMore}
          MenuProps={{
            getContentAnchorEl: null,
            anchorOrigin: {
              vertical: 'bottom',
              horizontal: 'left',
            },
          }}
          value={''}
          onChange={handleChange}
        >
          <FixedSizeList
            height={270}
            width={400}
            itemSize={35}
            itemCount={ItemsList.length}
          >
            {renderItems}
          </FixedSizeList>
        </Select>
      </Wrapper>
    ),
    [IsSelectCountriesOpen, ItemsList],
  );
}
const Wrapper = styled.div`
  width: 100%;
  position: relative;
  .inputField {
    position: absolute;
    z-index: 1;
    width: 90%;

    fieldset {
      border: none;
      background: transparent;
    }
    .MuiOutlinedInput-root {
      background: transparent !important;
    }
  }
`;
