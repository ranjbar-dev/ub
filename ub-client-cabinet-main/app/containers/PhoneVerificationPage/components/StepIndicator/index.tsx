/**
 *
 * StepIndicator
 *
 */
import React, { memo } from 'react';

// import styled from 'styles/styled-components';

import { Step } from 'containers/PhoneVerificationPage/types';
import styled from 'styles/styled-components';

interface Props {
  steps: Step[];
  selectesStep: number;
}

function StepIndicator (props: Props) {
  const steps = props.steps;
  const selectedIndex = props.selectesStep;
  let selectedStep: Step = {
    title: steps[0].title,
    description: steps[0].description,
  };
  for (let i = 0; i < steps.length; i++) {
    if (selectedIndex === i) {
      steps[i].isSelected = true;
      selectedStep = steps[i];
    }
  }

  return (
    <Wrapper>
      <div className='mb2vh'>
        <span className='title'>{selectedStep.title}</span>
        <span className='description'>{selectedStep.description}</span>
      </div>
      <div className='indicatorsWrapper'>
        {steps.map((item, index) => {
          return (
            <div className='lineAndCircleWrapper' key={'indic' + index}>
              <div
                className={`line ${index === 0 ? 'hidden' : ''} ${
                  index <= selectedIndex ? 'completed' : ''
                } `}
              ></div>
              <div
                className={`circle ${
                  index < selectedIndex ? 'completed' : ''
                } ${index === selectedIndex ? 'selected' : ''}`}
              ></div>
            </div>
          );
        })}
      </div>
    </Wrapper>
  );
}

export default memo(StepIndicator);
const Wrapper = styled.div`
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: space-evenly;
  flex: 3.5;
  place-content: flex-end;
  .title {
    color: var(--textBlue);
  }
  .description {
    color: var(--textGrey);
  }
  .indicatorsWrapper {
    display: flex;
  }
  .hidden {
    display: none;
  }
  .lineAndCircleWrapper {
    display: flex;
    align-items: center;
  }
  .circle {
    border: 1px solid var(--textGrey);
    border-radius: 10px;
    width: 12px;
    height: 12px;
    margin: 0 3px;
    &.completed {
      border: 1px solid var(--textBlue);
      background: var(--textBlue);
    }
    &.selected {
      border: 1px solid var(--textBlue);
    }
  }
  .line {
    width: 100px;
    height: 1px;
    background: var(--textGrey);
    &.completed {
      background: var(--textBlue);
    }
  }
`;
