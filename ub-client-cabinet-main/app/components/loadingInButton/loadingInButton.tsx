import React from 'react';
import { makeStyles, CircularProgress } from '@material-ui/core';

const materialClasses = makeStyles({
  loadingIndicator: {
    color: 'white',
  },
});
export default function LoadingInButton() {
  const classes = materialClasses();

  return <CircularProgress size={14} className={classes.loadingIndicator} />;
}
