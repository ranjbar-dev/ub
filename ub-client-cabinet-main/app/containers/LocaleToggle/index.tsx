/*
 *
 * LanguageToggle
 *
 */

import React from 'react';
import { createSelector } from 'reselect';
import { useSelector, useDispatch } from 'react-redux';
import Wrapper from './Wrapper';
import { changeLocale } from '../LanguageProvider/actions';
import { makeSelectLocale } from '../LanguageProvider/selectors';
import { makeStyles, FormControl, Select, MenuItem } from '@material-ui/core';
import ExpandMore from 'images/themedIcons/expandMore';

const stateSelector = createSelector(makeSelectLocale(), (locale) => ({
  locale,
}));
const useStyles = makeStyles((theme) => ({
  button: {
    display: 'block',
    marginTop: theme.spacing(2),
  },
  formControl: {
    margin: theme.spacing(1),
    minWidth: 120,
  },
}));

export default function LocaleToggle() {
  const { locale } = useSelector(stateSelector);
  const dispatch = useDispatch();
  const classes = useStyles();
  const [language, setLanguage] = React.useState(locale);
  const [open, setOpen] = React.useState(false);

  const handleChange = (event) => {
    setLanguage(event.target.value);
    onLocaleToggle(event);
  };

  const handleClose = () => {
    setOpen(false);
  };

  const handleOpen = () => {
    setOpen(true);
  };
  const onLocaleToggle = (evt) => dispatch(changeLocale(evt.target.value));

  return (
    <Wrapper>
      <FormControl className={classes.formControl}>
        <Select
          IconComponent={ExpandMore}
          MenuProps={{
            getContentAnchorEl: null,
            anchorOrigin: {
              vertical: 'bottom',
              horizontal: 'left',
            },
          }}
          open={open}
          onClose={handleClose}
          onOpen={handleOpen}
          value={language}
          onChange={handleChange}
        >
          <MenuItem value={'en'}>English</MenuItem>
          {/*<MenuItem value={'ar'}>Arabic</MenuItem>
          <MenuItem value={'de'}>Dutch</MenuItem>*/}
        </Select>
      </FormControl>
    </Wrapper>
  );
}
