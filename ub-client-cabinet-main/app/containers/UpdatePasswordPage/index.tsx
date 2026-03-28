/*
 *
 * UpdatePasswordPage
 *
 */

import React, { memo } from 'react';
import { useSelector } from 'react-redux';
import { createStructuredSelector } from 'reselect';

import { useInjectSaga } from 'utils/injectSaga';
import { useInjectReducer } from 'utils/injectReducer';
import { makeSelectLocation } from './selectors';
import reducer from './reducer';
import saga from './saga';
import { Card } from '@material-ui/core';
import StepSelector from './pages/stepSelector';

const stateSelector = createStructuredSelector({
  location: makeSelectLocation(),
});

interface Props {}

function UpdatePasswordPage(props: Props) {
  useInjectReducer({ key: 'updatePasswordPage', reducer: reducer });
  useInjectSaga({ key: 'updatePasswordPage', saga: saga });
  const { location } = useSelector(stateSelector);
  const loc: any = location;
  const query = loc.query;
  return (
    <Card style={{ height: '100vh' }}>
      <StepSelector {...query} />
    </Card>
  );
}

export default memo(UpdatePasswordPage);
