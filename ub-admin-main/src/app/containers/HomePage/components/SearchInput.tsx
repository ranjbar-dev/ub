import React, { memo, useRef } from 'react';
import styled from 'styled-components/macro';

interface Props {
  title: string;
  onEnter: (value: string) => void;
  isLast?: boolean;
  disabled?: boolean;
}

function SearchInput(props: Props) {
  const { title, onEnter, isLast, disabled } = props;
  const inputValue = useRef('');
  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    e.persist();
    inputValue.current = e.target.value;
  };
  const handleKeyUp = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === 'Enter' && inputValue.current !== '') {
      onEnter(inputValue.current);
    }
  };
  return (
    <Wrapper>
      <div
        className={
          'holder' +
          (isLast === true ? ' last' : '') +
          (disabled === true ? ' disabled' : '')
        }
      >
        <div className="searchTitle">
          {title}
          <span className="colon">:</span>
        </div>
        <div className="inputWrapper">
          <input
            type="text"
            onKeyUp={handleKeyUp}
            onChange={handleInputChange}
          />
        </div>
      </div>
    </Wrapper>
  );
}

export default memo(SearchInput);
const Wrapper = styled.div`
  padding: 10px 34px 0 34px;
  .holder {
    display: flex;
    align-items: center;
    width: 100%;
    padding-bottom: 9px;
    border-bottom: 1px solid #cfcccc;
    &.last {
      border-bottom: none;
    }
    &.disabled {
      pointer-events: none;
      filter: grayscale(1);
    }
  }
  .searchTitle {
    min-width: 50%;
    font-size: 12px;
    align-items: center;
  }
  .inputWrapper {
    width: 50%;
    display: flex;
    align-items: center;
    input {
      width: 100%;
    }
  }
`;
