import React, { useState, memo } from 'react';

import { Tabs, Tab } from '@material-ui/core';

import { FundsPages } from 'containers/FundsPage/constants';

import { SegmentSelectionWrapper } from 'components/wrappers/segmentSelectionWrapper';
interface SegmentOption {
  title: any;
  page: FundsPages;
}
const SegmentSelector = (props: {
  onChange: Function;
  options: SegmentOption[];
  defaultIndex?: number;
}) => {
  const [SelectedIndex, setSelectedIndex] = useState(
    props.defaultIndex ? props.defaultIndex : 0,
  );

  const handleSegmentClick = (e, index) => {
    setSelectedIndex(index);
    setTimeout(() => {
      props.onChange(props.options[index].page);
    }, 130);
  };
  return (
    //<SegmentSeletionButtonsWrapper>
    //  {props.options.map((item: SegmentOption, index: number) => {
    //    return (
    //      <Button
    //        onClick={(e) => {
    //          handleSegmentClick(e, index);
    //        }}
    //        disableRipple={true}
    //        className={`segmentButton ${
    //          index == SelectedIndex ? 'selected' : ''
    //        } ${index == 0 ? 'first' : ''} ${
    //          index == props.options.length - 1 ? 'last' : ''
    //        } `}
    //        id={index === 1 ? 'depositeTab' : index === 2 ? 'withdrawTab' : ''}
    //        key={'segment' + index}
    //      >
    //        <span>{item.title}</span>
    //      </Button>
    //    );
    //  })}
    //</SegmentSeletionButtonsWrapper>
    <SegmentSelectionWrapper className='funds'>
      <Tabs
        value={SelectedIndex}
        onChange={handleSegmentClick}
        indicatorColor='primary'
        textColor='primary'
      >
        {props.options.map((item: SegmentOption, index: number) => {
          return (
            <Tab
              disableRipple={true}
              className={`${index == 0 ? 'first' : ''} ${
                index == props.options.length - 1 ? 'last' : ''
              } `}
              id={
                index === 1 ? 'depositeTab' : index === 2 ? 'withdrawTab' : ''
              }
              key={'segment' + index}
              label={
                <>
                  <span>{item.title}</span>
                </>
              }
            />
          );
        })}
      </Tabs>
      <div className='mask'></div>
    </SegmentSelectionWrapper>
  );
};
export default memo(SegmentSelector);
