import {
  Checkbox,
  FormControl,
  FormLabel,
  FormGroup,
  FormControlLabel,
  OutlinedInput,
  InputAdornment,
  DialogTitle,
  DialogContent,
  DialogContentText,
  DialogActions,
  Button,
  makeStyles,
  CircularProgress,
} from '@material-ui/core';
import React, { memo } from 'react';

import DateFilter from './components/GridFilter/DateFilter';
import PopupModal from './components/materialModal/modal';
import PaginationComponent from './components/PaginationComponent/PaginationComponent';
import VerificationPage from './containers/UserAccounts/components/VerificationPage';
import EditDropDown from './containers/UserDetails/components/EditDropDown';

interface Props {}
const materialClasses = makeStyles({
  loadingIndicator: {
    color: 'white',
  },
});
function ForceStyles(props: Props) {
  const {} = props;
  const classes = materialClasses();
  return (
    <div className="stylePs" style={{ display: 'none' }}>
      <PopupModal isOpen={false} onClose={() => {}}>
        a
      </PopupModal>
      <FormControl component="fieldset">
        <FormLabel component="legend">Assign responsibility</FormLabel>
        <FormGroup>
          <FormControlLabel
            control={<Checkbox checked={true} />}
            label={'item.name'}
          />
        </FormGroup>
        <OutlinedInput
          endAdornment={<InputAdornment position="end">{''}</InputAdornment>}
        />
      </FormControl>
      <DialogTitle id="alert-dialog-title"></DialogTitle>
      <DialogContent>
        <DialogContentText id="alert-dialog-description"></DialogContentText>
      </DialogContent>
      <DialogActions>
        <Button onClick={() => {}} color="primary">
          Disagree
        </Button>
        <Button onClick={() => {}} color="primary">
          Agree
        </Button>
      </DialogActions>
      <CircularProgress size={14} className={classes.loadingIndicator} />;
      <DateFilter onDateSelect={e => {}} title="" />
      <PaginationComponent size={2} onPageChange={() => {}} />
      <EditDropDown
        initialValue={{ name: '', id: '' }}
        onSelect={(e: string) => {}}
        options={[{ name: '', id: '' }]}
      />
    </div>
  );
}

export default memo(ForceStyles);
