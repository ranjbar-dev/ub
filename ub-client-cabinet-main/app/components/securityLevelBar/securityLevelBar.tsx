import React from 'react';
import styled from 'styles/styled-components';
import { SecurityLevel } from 'containers/AcountPage/constants';

export default function SecurityLevelBar (props: { level: SecurityLevel }) {
  let backgroundColor = '';
  let width = '';
  switch (props.level) {
    case SecurityLevel.LOW:
      backgroundColor = 'red';
      width = '30%';
      break;
    case SecurityLevel.MEDIUM:
      backgroundColor = '#f69906';
      width = '60%';
      break;
    case SecurityLevel.High:
      backgroundColor = 'green';
      width = '100%';
      break;

    default:
      break;
  }
  return (
    <>
      <SecurityLevelBarWrapper>
        <div className='secBar'>
          <div className='wrapper'>
            <div
              className='bar'
              style={{
                width: `${width}`,
                background: `${backgroundColor}`,
              }}
            ></div>
          </div>
          <span style={{ color: `${backgroundColor}` }}>{props.level}</span>
        </div>
      </SecurityLevelBarWrapper>
    </>
  );
}
const SecurityLevelBarWrapper = styled.div`
  .secBar {
    width: 190px;
    display: flex;
    justify-content: space-between;
    align-content: center;
    span {
      min-width: 65px;
      margin-top: 1px;
      font-weight: 400 !important;
      margin-left: 5px;
      text-align: center;
    }
  }
  .wrapper {
    width: 8vw;
    height: 10px;
    background: var(--oddRows);
    border: 1px solid var(--lightGrey);
    border-radius: 10px;
    margin-top: 6px;
    .bar {
      border-radius: 10px;
      margin-top: -1px;
      height: 10px;
    }
  }
`;
