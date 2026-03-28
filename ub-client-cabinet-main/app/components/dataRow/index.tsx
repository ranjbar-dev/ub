import * as React from 'react';
import styled, { css } from 'styles/styled-components';
import external from 'images/themedIcons/externalLink.svg';
import TextLoader from 'components/textLoader';

const DataRow = (props: {
  title: any;
  value: any;
  boldValue?: boolean;
  small?: boolean;
  isLoading?: boolean;
  hideWhenSmall?: boolean;
  titleWidth?: string;
  style?: any;
  dense?: boolean;
  clickAddress?: string;
  className?: string;
}) => {
  const style = props.style ? props.style : {};
  return (
    <Wrapper
      hideWhenSmall={props.hideWhenSmall}
      style={style}
      className={'RowData'}
    >
      <div
        className={`${props.className ? props.className : ''} title ${
          props.small ? 'small' : ''
        }  ${props.dense ? 'dense' : ''} `}
        style={{
          width: props.titleWidth,
          minWidth: props.titleWidth,
          maxWidth: props.titleWidth,
        }}
      >
        {props.title}
      </div>
      <div
        className={`${props.className ? props.className : ''} value ${
          props.boldValue ? 'bolder' : ''
        } ${props.small ? 'small' : ''}  ${
          props.dense ? 'dense' : ''
        } ${!props.isLoading && props.clickAddress && 'clickable'}`}
        onClick={() => {
          if (!props.clickAddress) {
            return;
          } else {
            window.open(props.clickAddress);
          }
        }}
      >
        {props.isLoading === true ? (
          <TextLoader width={200} height={20} />
        ) : props.value ? (
          props.value
        ) : (
          '-'
        )}
        {props.clickAddress && <img src={external} />}
      </div>
    </Wrapper>
  );
};
export default DataRow;
const Wrapper = styled.div<{ hideWhenSmall?: boolean }>`
  display: flex;
  ${({ hideWhenSmall }) =>
    hideWhenSmall &&
    css`
      @media screen and (max-height: 725px) {
        display: none;
      }
    `}
  flex-basis: 40px;
  max-height: 27px;
  .title {
    flex: 1;
    min-width: 150px;
    max-width: 150px;
    margin-top: 1px;
    span {
      font-size: 13px !important;
    }
    &.small {
      font-size: 12px !important;
      span {
        font-size: 12px !important;
      }
    }
    &.dense {
      flex: unset;
      width: unset !important;
      min-width: unset !important;
      max-width: unset !important;
      padding-right: 10px;
      line-height: 34px;
    }
    &.noValue {
      text-decoration: line-through;
    }
  }
  .value {
    flex: 2;
    font-weight: 600;
    font-size: 13px !important;
    margin-top: 2px;
    color: var(--blackText);
    img {
      margin-top: -3px;
      margin-left: 2px;
    }
    &.clickable {
      cursor: pointer;
    }
    &.bolder {
      font-weight: 700;
    }
    &.small {
      font-size: 12px !important;
      span {
        font-size: 12px !important;
      }
    }
  }
`;
