import React, { memo, useState, useRef } from 'react';
import '@vaadin/vaadin-date-picker';
import styled from 'styled-components/macro';
import { filterHeight } from 'utils/gridUtilities';

interface Props {
  title: string;
  /** Called with the selected date string (ISO format), or empty string on clear. */
  onDateSelect: (value: string) => void;
}

/**
 * Date picker filter input for AG Grid column headers.
 * Uses a Vaadin date-picker web component under the hood with a
 * styled clickable overlay for consistent appearance.
 *
 * @example
 * ```tsx
 * <DateFilter title="Created At" onDateSelect={(date) => applyFilter(date)} />
 * ```
 */
function DateFilter(props: Props) {
  const { onDateSelect, title } = props;
  const [Date, setDate] = useState(title);
  const DatePicker = useRef<HTMLElement & { open(): void }>(null);

  customElements.whenDefined('vaadin-date-picker').then(function () {
    if (DatePicker.current) {
      DatePicker.current.addEventListener('change', function (event: Event) {
        const target = event.target as HTMLInputElement;
        if (target && 'value' in target) {
          setDate(target.value);
          onDateSelect(target.value);
        }
      });
    }
  });

  return (
    <Wrapper>
      <div
        className="fake"
        onClick={() => {
          try {
            if (DatePicker.current) {
              DatePicker.current.open();
            }
          } catch {
            //console.log(BcpDatePicker.current);

            //BcpDatePicker.current.click();
            //BcpDatePicker.current.focus();
          }
        }}
      >
        {Date}
        {Date !== title && (
          <div
            onClick={e => {
              e.stopPropagation();
              onDateSelect('');
              setDate(title);
            }}
            className="x"
          >
            X
          </div>
        )}
      </div>
      {/*<TextField
        onChange={e => {
          console.log(e.target.value);
        }}
        ref={BcpDatePicker}
        id="datetime-local"
        label={title}
        type="datetime-local"
        variant="outlined"
        margin="dense"
        //defaultValue="2017-05-24T10:30"
        InputLabelProps={{
          shrink: true,
        }}
      />*/}
      <div
        style={{
          position: 'absolute',
          opacity: '0',
          pointerEvents: 'none',
          top: '0',
          height: '40px',
        }}
      >
        <vaadin-date-picker
          ref={DatePicker}
          placeholder="Start Time"
        ></vaadin-date-picker>
      </div>
    </Wrapper>
  );
}
const Wrapper = styled.div`
  .fake {
    cursor: pointer;
    display: flex;
    font-size: 13px;
    line-height: 28px;
    width: 100%;
    padding-left: 10px;
    height: ${filterHeight}px;
    border: 1px solid #e6e3e3;
    border-radius: 5px;
    color: ${p => p.theme.textGrey} !important;
    position: relative;
  }
  .x {
    position: absolute;
    right: 10px;
  }
`;
export default memo(DateFilter);
