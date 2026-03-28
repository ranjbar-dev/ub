import { makeStyles, CircularProgress } from '@material-ui/core';
import React from 'react';

/**
 * Small circular loading spinner for use inside button components.
 * Colour adapts based on the `lightBackground` flag.
 *
 * @example
 * ```tsx
 * <LoadingInButton lightBackground={false} />
 * ```
 */
export default function LoadingInButton(props: { lightBackground?: boolean }) {
  const materialClasses = makeStyles({
    loadingIndicator: {
      color: props.lightBackground === true ? 'blue' : 'white',
    },
  });
  const classes = materialClasses();

  return <CircularProgress size={14} className={classes.loadingIndicator} />;
}
