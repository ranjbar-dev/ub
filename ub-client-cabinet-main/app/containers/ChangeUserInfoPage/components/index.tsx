import React, { useState, useEffect } from 'react';
import translate from '../messages';
import { FormattedMessage } from 'react-intl';
import Title from './title';
import WithTitle from './withTitle';
import InputWithValidator from '../../../components/inputWithValidator';
import { UserProfileData } from 'containers/DocumentVerificationPage/types';
import RadioButtons from './radioButtons';
import BirthDaySelector from './birthDaySelector';
import SelectRegion from './selectRegion';
import { Button, Card } from '@material-ui/core';
import { Buttons, AppPages } from 'containers/App/constants';
import styled from 'styles/styled-components';
import {
  MessageService,
  MessageNames,
  Subscriber,
} from 'services/message_service';
import { Country } from 'containers/PhoneVerificationPage/constants';
import IsLoadingWithText from 'components/isLoadingWithText/isLoadingWithText';
import { useDispatch } from 'react-redux';
import { updateUserProfileDataAction } from '../actions';
import { replace, push } from 'redux-first-history';
import { MaxContainer } from 'components/wrappers/maxContainer';

let dataToSend: UserProfileData = {
  country_id: 799,
};
let filledBefore = false;
export default function ComponentsWrapper (props: { data: UserProfileData }) {
  const [CanSubmit, setCanSubmit] = useState(false);
  const [PageData, setPageData] = useState(props.data);
  const [IsLoading, setIsLoading] = useState(false);
  const dispatch = useDispatch();
  if (PageData.firstName && PageData.firstName !== '') {
    filledBefore = true;
  }
  if (PageData.changed === true) {
    filledBefore = false;
  }
  const handleChange = (data: { field: string; value: string }) => {
    dataToSend[data.field] = data.value;
    if (!dataToSend.country_id) {
      dataToSend.country_id = 799;
    }
    const oldState = { ...PageData };
    oldState[data.field] = data.value;
    oldState.changed = true;
    if (oldState[data.field] == null || oldState[data.field].length === 0) {
      MessageService.send({
        name: MessageNames.SET_INPUT_ERROR,
        value: data.field,
        payload: <FormattedMessage {...translate.Required} />,
      });
      setCanSubmit(false);
    } else {
      MessageService.send({
        name: MessageNames.SET_INPUT_ERROR,
        value: data.field,
        payload: null,
      });
    }

    setPageData(oldState);

    validateArray([
      dataToSend.firstName,
      dataToSend.lastName,
      dataToSend.postalCode,
      dataToSend.address,
    ]);
  };
  const validateArray = (data: any[]) => {
    const expected = data.length;
    let counter = 0;
    for (let i = 0; i < data.length; i++) {
      if (data[i] && data[i].length > 0) {
        counter++;
      }
    }
    if (counter === expected) {
      setCanSubmit(true);
      return;
    }
    setCanSubmit(false);
  };
  const onCountrySelect = (e: Country) => {
    dataToSend.country = e.id;
    dataToSend.country_id = e.id;
    setPageData({ ...PageData, changed: true });
  };
  const handleSubmit = () => {
    const data = {
      first_name: dataToSend.firstName,
      last_name: dataToSend.lastName,
      gender: dataToSend.gender,
      date_of_birth: dataToSend.dateOfBirth,
      address: dataToSend.address,
      region_and_city: dataToSend.city,
      postal_code: dataToSend.postalCode,
      country_id: dataToSend.country_id,
    };
    dispatch(updateUserProfileDataAction(data));
  };
  useEffect(() => {
    dataToSend = props.data;
    dataToSend.country_id = PageData.country;
    dataToSend.city = PageData.regionAndCity;
    validateArray([
      dataToSend.firstName,
      dataToSend.lastName,
      dataToSend.postalCode,
      dataToSend.address,
    ]);
    const Subscription = Subscriber.subscribe((message: any) => {
      if (message.name === MessageNames.SETLOADING) {
        setIsLoading(message.payload);
      }
    });
    return () => {
      dataToSend = {};
      Subscription.unsubscribe();
    };
  }, []);
  return (
    <MainWrapper>
      <MaxContainer>
        <Title flex={3} />
        <WithTitle
          flex={40}
          title={<FormattedMessage {...translate.Basicinfo} />}
        >
          <InputWithValidator
            uniqueName='firstName'
            initialValue={PageData.firstName}
            label={<FormattedMessage {...translate.FirstName} />}
            onChange={value => {
              handleChange({ field: 'firstName', value });
            }}
          />
          <InputWithValidator
            uniqueName='lastName'
            initialValue={PageData.lastName}
            label={<FormattedMessage {...translate.LastName} />}
            onChange={value => {
              handleChange({ field: 'lastName', value });
            }}
          />
          <RadioButtons
            initialValue={PageData.gender}
            onChange={value => {
              handleChange({ field: 'gender', value });
            }}
          />
          <BirthDaySelector
            initialValue={
              PageData.dateOfBirth ? PageData.dateOfBirth : 'Birthday'
            }
            onDateSelect={value => {
              handleChange({ field: 'dateOfBirth', value });
            }}
          />
        </WithTitle>
        <WithTitle
          flex={60}
          title={<FormattedMessage {...translate.Residentialaddress} />}
        >
          <SelectRegion
            initialCountryId={PageData.country ? PageData.country : -1}
            onCountrySelect={onCountrySelect}
          />
          <InputWithValidator
            uniqueName='city'
            initialValue={PageData.regionAndCity}
            label={<FormattedMessage {...translate.city} />}
            onChange={value => {
              handleChange({ field: 'city', value });
            }}
          />
          <InputWithValidator
            uniqueName='postalCode'
            initialValue={PageData.postalCode}
            label={<FormattedMessage {...translate.postalCode} />}
            onChange={value => {
              handleChange({ field: 'postalCode', value });
            }}
          />
          <InputWithValidator
            uniqueName='address'
            initialValue={PageData.address}
            rows={2}
            label={<FormattedMessage {...translate.address} />}
            onChange={value => {
              handleChange({ field: 'address', value });
            }}
          />
        </WithTitle>
        <div className='buttonsWrapper'>
          <Button
            variant='contained'
            onClick={
              filledBefore == false
                ? handleSubmit
                : () => {
                    dispatch(push(AppPages.DocumentVerification));
                  }
            }
            disabled={!CanSubmit}
            color='primary'
          >
            <IsLoadingWithText
              isLoading={IsLoading}
              text={
                filledBefore === false ? (
                  <FormattedMessage {...translate.BeginVerification} />
                ) : (
                  <FormattedMessage {...translate.next} />
                )
              }
            />
          </Button>
          <Button
            onClick={() => {
              dispatch(replace(AppPages.AcountPage));
            }}
            className={`fitted mt12 ${Buttons.CancelButton}`}
          >
            <FormattedMessage {...translate.cancel} />
          </Button>
        </div>
      </MaxContainer>
    </MainWrapper>
  );
}
const MainWrapper = styled(Card)`
  display: flex;
  justify-content: center;
  align-items: center;
  flex-direction: column;
  border-radius: 10px !important;
  box-shadow: none !important;
  height: calc(98vh - 115px);
  padding: 48px;
  overflow: auto !important;

  .buttonsWrapper {
    display: flex;
    flex-direction: column;
    flex: 55;
    justify-content: flex-start;
    align-items: center;
    margin-top: 16px;
    .fitted {
      max-width: fit-content;
      min-width: max-content;
      span {
        color: var(--textGrey);
      }
    }
  }
  .loadingCircle {
    top: 8px !important;
  }

  .MuiInputBase-root {
    .MuiSvgIcon-root path {
      fill: var(--textGrey) !important;
    }
  }
  .Mui-checked path {
    fill: var(--textBlue) !important;
  }
  .inputWithValidator {
    margin: 1vh 0 0 0;
  }
`;
