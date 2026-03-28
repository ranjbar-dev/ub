import React from 'react';
import translate from '../messages';
import { FormattedMessage } from 'react-intl';
import styled from 'styles/styled-components';

export default function Title (props) {
  return (
    <Wrapper style={{ flex: props.flex }}>
      <FormattedMessage
        {...translate.TheinformationyoufillinmustbeconsistentwiththeinformationinyourIDdocuments}
      />
    </Wrapper>
  );
}
const Wrapper = styled.div`
  span {
    font-size: 13px;
    color: var(--blackText);
  }
`;
