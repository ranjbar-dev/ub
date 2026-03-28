/*
 *
 * ContactUsPage
 *
 */
import React, { ChangeEventHandler, useEffect, useState } from 'react';
import { Helmet } from 'react-helmet';
import { useSelector, useDispatch } from 'react-redux';
import { createStructuredSelector } from 'reselect';

import { useInjectSaga } from 'utils/injectSaga';
import { useInjectReducer } from 'utils/injectReducer';
import makeSelectContactUsPage, {
  makeSelectCounter,
  makeSelectInputValue,
} from './selectors';
import reducer from './reducer';
import saga from './saga';
import styled from 'styles/styled-components';
import {
  addValueByInputAction,
  changeInputValueAction,
  decrementAction,
  getNumberFromApi,
  incrementAction,
  subtractValueByInputAction,
} from './actions';
import { css } from 'styled-components';
import { toast } from 'components/Customized/react-toastify';
import { MessageNames, Subscriber } from 'services/message_service';
import { CircularProgress } from '@material-ui/core';
import OrderList from './components';

const stateSelector = createStructuredSelector({
  contactUsPage: makeSelectContactUsPage(),
  counterValue: makeSelectCounter(),
  inputValue: makeSelectInputValue(),
});

interface Props {}

function ContactUsPage(props: Props) {
  // Warning: Add your key to RootState in types/index.d.ts file
  useInjectReducer({ key: 'contactUsPage', reducer: reducer });
  useInjectSaga({ key: 'contactUsPage', saga: saga });
  const [inputSubOrAddSelected, setInputSubOrAddSelected] = useState({
    sub: false,
    add: true,
  });
  const { contactUsPage, counterValue, inputValue } = useSelector(
    stateSelector,
  );
  const dispatch = useDispatch();

  const handleIncrementClick = () => {
    dispatch(incrementAction());
  };

  const handleDecrementClick = () => {
    dispatch(decrementAction());
  };

  const handleApiButtonClick = () => {
    dispatch(getNumberFromApi());
  };

  const [buttonLoading, setbuttonLoading] = useState(false);

  useEffect(() => {
    const Subscription = Subscriber.subscribe((message: any) => {
      if (message.name === MessageNames.SET_LOADING_TEST) {
        setbuttonLoading(true);
      }

      if (message.name === MessageNames.SET_LOADING_END) {
        setbuttonLoading(false);
      }
    });

    return () => {
      Subscription.unsubscribe();
    };
  }, [buttonLoading]);

  const handleSubmitInput = () => {
    if (inputSubOrAddSelected.add && !isNaN(parseFloat(inputValue))) {
      return dispatch(addValueByInputAction(parseFloat(inputValue)));
    }
    if (inputSubOrAddSelected.sub && !isNaN(parseFloat(inputValue))) {
      return dispatch(subtractValueByInputAction(parseFloat(inputValue)));
    }
    return toast.error('Please put a number');
  };

  const handleInputChange: ChangeEventHandler<HTMLInputElement> = (e) => {
    dispatch(changeInputValueAction(e.target.value));
  };

  const handleInputAddClick = () => {
    setInputSubOrAddSelected({ add: true, sub: false });
  };

  const handleInputSubtractClick = () => {
    setInputSubOrAddSelected({ add: false, sub: true });
  };

  return (
    <div>
      <Helmet>
        <title>ContactUsPage</title>
        <meta name="description" content="Description of ContactUsPage" />
      </Helmet>
      <Wrapper>
        {/* <div>Click to increment | decrement!</div>
        <NumberWrapper>
          <Button id="Increment" onClick={handleIncrementClick}>
            +
          </Button>
          <div>{counterValue}</div>
          <Button id="Decrement" onClick={handleDecrementClick}>
            -
          </Button>
        </NumberWrapper>
        <Input value={inputValue} type="number" onChange={handleInputChange} />
        <NumberWrapper>
          <InputButton
            onClick={handleInputAddClick}
            isSelected={inputSubOrAddSelected.add}
          >
            Addition
          </InputButton>
          <InputButton
            onClick={handleInputSubtractClick}
            isSelected={inputSubOrAddSelected.sub}
          >
            Subtract
          </InputButton>
        </NumberWrapper>
        <Button
          style={{ margin: '10px 0', fontSize: '16px' }}
          id="Submit"
          onClick={handleSubmitInput}
        >
          Submit
        </Button>
        <Button
          style={{ margin: '10px 0', fontSize: '16px' }}
          id="API"
          onClick={buttonLoading ? () => {} : handleApiButtonClick}
        >
          {buttonLoading ? <CircularProgress size={22} /> : 'API'}
        </Button> */}
        <OrderList />
      </Wrapper>
    </div>
  );
}

const Wrapper = styled.section`
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  min-width: 100vw;
  height: 90vh;
  background-image: linear-gradient(#614385, #516395);
  div {
    color: white;
  }
`;
const buttonCss = css`
  display: flex;
  cursor: pointer;
  flex-direction: row;
  justify-content: center;
  align-items: center;
  user-select: none;
  border-radius: 5px;
`;
const NumberWrapper = styled.section`
  display: flex;
  flex-direction: row;
  justify-content: space-evenly;
  align-items: center;
  width: 300px;
  margin-top: 20px;
`;
const Button = styled.div`
  width: 80px;
  height: 30px;
  font-weight: 800;
  font-size: 25px;
  background-image: linear-gradient(#2193b0, #6dd5ed);
  ${buttonCss}
`;

const Input = styled.input`
  width: 235px;
  margin-top: 20px;
`;
const InputButton = styled.div<{ isSelected: boolean }>`
  transition: all 0.2s;
  ${({ isSelected }) =>
    isSelected
      ? css`
          background-color: #f45c43;
        `
      : css`
          background-color: #2193b0;
        `}
  ${buttonCss}
  width: 80px;
  height: 30px;
`;
export default ContactUsPage;
{
  /* <FormattedMessage {...messages.header} />
        <FormattedMessage {...messages.king} /> */
}
