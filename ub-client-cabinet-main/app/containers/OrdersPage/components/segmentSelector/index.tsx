import React, { useState, memo, useCallback } from 'react';
import { Tabs, Tab } from '@material-ui/core';

import { OrderPages } from 'containers/OrdersPage/constants';
import { SegmentSelectionWrapper } from 'components/wrappers/segmentSelectionWrapper';
interface SegmentOption {
  title: any;
  page: OrderPages;
}
const SegmentSelector = (props: {
  onChange: Function;
  options: SegmentOption[];
}) => {
  const [SelectedIndex, setSelectedIndex] = useState(0);

  const handleSegmentClick = useCallback((e, index) => {
    setSelectedIndex(index);
    setTimeout(() => {
      props.onChange(props.options[index].page);
    }, 130);
  }, []);
  return (
    //<SegmentSeletionButtonsWrapper className="orders">
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
    //        key={'segment' + index}
    //      >
    //        <span>{item.title}</span>
    //      </Button>
    //    );
    //  })}
    //</SegmentSeletionButtonsWrapper>

    <SegmentSelectionWrapper>
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
              className={`${index === 0 ? 'first' : ''} ${
                index === props.options.length - 1 ? 'last' : ''
              }  `}
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
