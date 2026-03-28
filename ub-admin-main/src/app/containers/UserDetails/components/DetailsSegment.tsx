import BorderColorOutlinedIcon from "@material-ui/icons/BorderColorOutlined";
import SaveOutlinedIcon from "@material-ui/icons/SaveOutlined";
import CountryDropDown from "app/components/CountryDropDown/CountryDropDown";
import IsLoadingWithTextAuto from "app/components/isLoadingWithText/isLoadingWithTextAuto";
import { Buttons, WindowTypes } from "app/constants";
import {
  UserAccountsActions,
  useUserAccountsSlice,
} from "app/containers/UserAccounts/slice";
import { InitialUserDetails } from "app/containers/UserAccounts/types";
import { translations } from "locales/i18n";
import React, { memo, useRef, useCallback, useState } from "react";
import { useTranslation } from "react-i18next";
import { useDispatch } from "react-redux";
import styled from "styled-components/macro";
import { CurrencyFormater } from "utils/formatters";

import { UserDetailsActions } from "../slice";
import EditDropDown from "./EditDropDown";
interface Props {
  data: InitialUserDetails;
  inValueClick?: () => void;
}

function DetailsSegment(props: Props) {
  const { data } = props;
  const [IsEdditing, setIsEdditing] = useState(false);
  const { t } = useTranslation();
  const dispatch = useDispatch();
  const dataTosend = useRef<Record<string, unknown>>({ id: data.id });

  const DataRow = useCallback(
    (propes: {
      title: string;
      value: React.ReactNode;
      editable?: boolean;
      kkey?: string;
      className?: string;
      onClick?: () => void;
    }) => {
      return (
        <div className="detailRow">
          <div className="title">
            <span className="title">{propes.title} : </span>
          </div>
          {propes.editable && IsEdditing === false && (
            <span className="value">{propes.value}</span>
          )}
          {propes.editable && IsEdditing === true && (
            <span className="value">
              <input
                type="text"
                className="udInput"
                placeholder={typeof propes.value === 'string' ? propes.value : undefined}
                onChange={(e) => {
                  if (propes.kkey) {
                    dataTosend.current[propes.kkey] = e.target.value;
                  }
                }}
              />
            </span>
          )}

          {!propes.editable && (
            <span
              onClick={propes.onClick}
              className={"value" + " " + propes.className ?? ""}
            >
              {typeof propes.value === "string"
                ? propes.value.toLowerCase()
                : propes.value}
            </span>
          )}
        </div>
      );
    },
    [IsEdditing]
  );
  const handleSubmitClick = async () => {
    //console.log(dataTosend.current);
    await dispatch(UserDetailsActions.UpdateUserDataAction(dataTosend.current));
    setIsEdditing((IsEdditing) => false);
  };
  const handleEditClick = () => {
    setIsEdditing((IsEdditing) => !IsEdditing);
  };

  const handleIdentityClick = async () => {
    dispatch(
      UserAccountsActions.getInitialSingleUserDataAndOpenWindowAction({
        id: props.data.id,
        windowType: WindowTypes.Verification,
      })
    );
  };
  const addressClickable =
    props.data.addressConfirmationStatus === "INCOMPLETE";
  const identityClickable =
    props.data.identityConfirmationStatus === "INCOMPLETE";
  return (
    <Wrapper>
      <div className="column">
        <DataRow
          title={t(translations.CommonTitles.Name())}
          //  value={<EditableValue initialValue={props.data.fullName} />}
          value={props.data.fullName}
          kkey={"full_name"}
          editable={true}
        />
        <DataRow
          title={t(translations.CommonTitles.Email())}
          value={props.data.email}
          kkey="email"
          //  editable={true}
        />
        <DataRow
          title={t(translations.CommonTitles.Gender())}
          value={props.data.gender}
        />
        <DataRow
          title={t(translations.CommonTitles.Birthdate())}
          value={props.data.birthDate}
          kkey="birthDate"
          editable={true}
        />
        <DataRow
          title={t(translations.CommonTitles.Mobile())}
          value={props.data.mobile}
          editable={true}
          kkey="phone"
        />
        <DataRow
          title={t(translations.CommonTitles.Registration())}
          value={props.data.registrationDate}
        />
        <DataRow
          title={t(translations.CommonTitles.RegisteredIP())}
          value={props.data.registeredIp}
        />
        <DataRow
          title={t(translations.CommonTitles.Country())}
          value={
            IsEdditing === true ? (
              <CountryDropDown
                initialCountryId={props.data.countryId}
                onChange={(countryId) => {
                  dataTosend.current["country_id"] = countryId;
                }}
                onClear={() => {}}
              />
            ) : (
              props.data.country
            )
          }
        />
        <DataRow
          title={t(translations.CommonTitles.City())}
          value={props.data.city}
          kkey="city"
          editable={true}
        />
        <DataRow
          title={t(translations.CommonTitles.PostalCode())}
          value={props.data.postalCode}
          kkey="postalCode"
          editable={true}
        />
        <DataRow
          title={t(translations.CommonTitles.Address())}
          value={props.data.address}
          kkey="address"
          editable={true}
        />
      </div>
      <div className="column">
        <DataRow
          onClick={identityClickable ? handleIdentityClick : () => {}}
          title={t(translations.CommonTitles.Identity())}
          value={props.data.identityConfirmationStatus}
          className={`${props.data.identityConfirmationStatus}`}
        />
        <DataRow
          onClick={addressClickable ? handleIdentityClick : () => {}}
          title={t(translations.CommonTitles.Address())}
          value={props.data.addressConfirmationStatus}
          className={`${props.data.addressConfirmationStatus}`}
        />
        {/*<DataRow
          title={t(translations.CommonTitles.Mobile())}
          value={
            <EditDropDown
              onSelect={(e: string) => {
                dataTosend.current['phone'] = e;
              }}
              initialValue={{
                name: data.mobile ?? 'Not Verified',
                id: data.mobile,
              }}
              options={[
                {
                  name: 'Verified',
                  id: 'Verified',
                },
                {
                  name: 'Not Verified',
                  id: 'null',
                },
              ]}
            />
          }
        />*/}
        <DataRow
          title={t(translations.CommonTitles.UserLevel())}
          value={
            <EditDropDown
              onSelect={(e: string) => {
                dataTosend.current["user_level_id"] = Number(e);
              }}
              initialValue={{
                ...(data.metaData.userLevels.find(
                  (item) => Number(item.id) === Number(data.userLevelId)
                ) || data.metaData.userLevels[0]),
                id: String((data.metaData.userLevels.find(
                  (item) => Number(item.id) === Number(data.userLevelId)
                ) || data.metaData.userLevels[0])?.id)
              }}
              options={data.metaData.userLevels}
            />
          }
        />
        <DataRow
          title={t(translations.CommonTitles.Group())}
          value={
            <EditDropDown
              onSelect={(e: string) => {
                dataTosend.current["group_id"] = Number(e);
              }}
              initialValue={data.metaData.userGroups.find(
                (item) => item.id === data.groupId
              )}
              options={data.metaData.userGroups}
            />
          }
        />
        <DataRow title={t(translations.CommonTitles.Manager())} value={""} />
        <DataRow
          title={t(translations.CommonTitles.Status())}
          value={
            <EditDropDown
              onSelect={(e: string) => {
                dataTosend.current["status"] = e;
              }}
              initialValue={{ name: data.status, id: data.status }}
              options={data.metaData.userStatuses.map((item, index) => {
                return {
                  name: item.name.toLowerCase(),
                  id: item.name,
                };
              })}
            />
          }
        />
        <DataRow
          title={t(translations.CommonTitles.TrustLevel())}
          value={
            <EditDropDown
              onSelect={(e: string) => {
                dataTosend.current["trust_level"] = Number(e);
              }}
              initialValue={{
                name: data.trustLevel + "",
                id: data.trustLevel + "",
              }}
              options={Array.from(Array(10).keys()).map((item, index) => {
                return {
                  name: index + "",
                  id: index + "",
                };
              })}
            />
          }
        />
        <DataRow
          title={t(translations.CommonTitles.SystemID())}
          value={props.data.systemId}
        />
        <DataRow
          title={t(translations.CommonTitles.ReferralID())}
          value={props.data.referralId}
        />
        <DataRow
          title={t(translations.CommonTitles.ReferKey())}
          value={props.data.referKey}
        />
      </div>
      <div className="column">
        <DataRow
          title={t(translations.CommonTitles.TotalDeposit())}
          value={props.data.totalDeposit}
        />
        <DataRow
          title={t(translations.CommonTitles.TotalWithdraw())}
          value={props.data.totalWithdraw}
        />
        <DataRow
          title={t(translations.CommonTitles.TotalCommissions())}
          value={props.data.totalCommissions}
        />
        <DataRow
          title={t(translations.CommonTitles.TotalBalances())}
          value={CurrencyFormater(props.data.totalBalance)}
        />
        <DataRow
          title={t(translations.CommonTitles.TotalOntrade())}
          value={props.data.totalOnTrade}
        />
      </div>
      <div className="submitContainer">
        <IsLoadingWithTextAuto
          onClick={handleEditClick}
          text={t(translations.CommonTitles.Edit())}
          loadingId={"edittButton" + data.id}
          className={Buttons.BlackButton}
          icon={<BorderColorOutlinedIcon />}
        />
        <IsLoadingWithTextAuto
          onClick={handleSubmitClick}
          text={t(translations.CommonTitles.SaveChanges())}
          loadingId={"userEdit"}
          className={Buttons.SkyBlueButton}
          icon={<SaveOutlinedIcon />}
        />
      </div>
    </Wrapper>
  );
}

export default memo(DetailsSegment);
const Wrapper = styled.div`
  display: flex;
  padding: 10px 40px;
  background: rgb(255, 255, 255);
  width: calc(100% - 24px);
  border-radius: 7px;
  min-height: calc(100% + 53px);
  margin-top: -20px;
  border: 1px solid rgb(222, 222, 222);
  margin-left: 12px;

  .column {
    display: flex;
    flex-direction: column;
    min-width: 310px;
  }
  .detailRow {
    margin: 10px 0px;
    display: flex;
    .title,
    .value {
      color: ${(p) => p.theme.blackText};
      min-width: 90px;
      max-width: 160px;
      font-size: 12px;
      text-transform: capitalize;
      &.REJECTED,
      &.INCOMPLETE {
        color: ${(p) => p.theme.red};
      }
      &.VERIFIED,
      &.CONFIRMED {
        color: ${(p) => p.theme.green};
      }
      &.INCOMPLETE {
        cursor: pointer;
      }
    }
    .value {
      font-weight: 600;
    }
  }
  .submitContainer {
    position: absolute;
    bottom: 60px;
    right: 65px;
    button {
      margin-right: 12px;
    }
  }
`;
