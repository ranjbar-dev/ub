/*
 *
 * EmailAuthentication
 *
 */

import React, { useEffect } from 'react';
import { useSelector, useDispatch } from 'react-redux';
import { createStructuredSelector } from 'reselect';

import { useInjectSaga } from 'utils/injectSaga';
import { useInjectReducer } from 'utils/injectReducer';
import { makeSelectLocation } from './selectors';
import reducer from './reducer';
import saga from './saga';
import { acountActivationAction } from './actions';
import StepSelector from './pages/stepSelector';
import { Card } from '@material-ui/core';

const stateSelector = createStructuredSelector({
  location: makeSelectLocation(),
});

interface Props {}

function EmailVerificationPage(props: Props) {
  useInjectReducer({ key: 'emailAuthentication', reducer: reducer });
  useInjectSaga({ key: 'emailAuthentication', saga: saga });
  const dispatch = useDispatch();
  const { location } = useSelector(stateSelector);

  useEffect(() => {
    const loc: any = location;
    if (loc && loc.query && loc.query.code) {
      dispatch(acountActivationAction({ code: loc.query.code }));
    }
    return () => {};
  }, []);

  return (
    <Card>
      <StepSelector />
    </Card>
  );
}

export default EmailVerificationPage;
