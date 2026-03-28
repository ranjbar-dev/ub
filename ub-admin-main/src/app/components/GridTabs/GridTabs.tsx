import {GridTopTab} from 'locales/types';
import React,{memo,useState,useCallback,CSSProperties} from 'react';
import styled from 'styled-components/macro';

interface Props {
	tabs: GridTopTab[];
	/** Called with the tab's callObject when a tab is selected. */
	onChange: (callObject: Record<string, unknown>) => void;
	styles?: CSSProperties
}

/**
 * Tab bar rendered above a SimpleGrid, allowing the user to switch between
 * pre-defined server-side filter presets (e.g. "All", "Pending", "Active").
 *
 * @example
 * ```tsx
 * <GridTabs
 *   tabs={[{ name: 'All', callObject: {} }, { name: 'Pending', callObject: { status: 'pending' } }]}
 *   onChange={(params) => applyParams(params)}
 * />
 * ```
 */
function GridTabs(props: Props) {
	const {tabs,onChange,styles}=props;
	const [SelectedIndex,setSelectedIndex]=useState(0);
	const handleTabClick=useCallback(
		(index: number) => {
			setSelectedIndex(index);
			onChange(tabs[index].callObject as Record<string, unknown>);
		},
		[onChange,tabs],
	);
	return (
		<Wrapper style={styles} >
			{tabs.map((item,index) => {
				return (
					<div
						key={item.name}
						onClick={() => handleTabClick(index)}
						className={`button ${SelectedIndex===index? 'active':''}`}
					>
						<span style={{marginTop: '-9px'}}> {item.name}</span>
					</div>
				);
			})}
		</Wrapper>
	);
}

export default memo(GridTabs);
const Wrapper=styled.div`
  display: flex;
  min-height: 40px;
  box-shadow: 16px 0px 20px 9px ${p => p.theme.greyBackground};
  margin-top: -20px;
  margin-bottom: 10px;
  background: ${p => p.theme.greyBackground};
  /*border-bottom: 1px solid #eee;*/
  .button {
    padding: 5px 15px;
    display: flex;

    transition-property: color, border;
    transition-duration: 0.5s, 0.2s;
    color: ${p => p.theme.textGrey};
    font-size: 12px;
    font-weight: 600;
    background: ${p => p.theme.tabBackground};
    cursor: pointer;
    &.active {
      box-shadow: none !important;
      color: ${p => p.theme.blackText};
      background: ${p => p.theme.white};
    }
    border-bottom: none !important;
    border-top-right-radius: 5px;
    border-top-left-radius: 5px;
    align-items: center;
    box-shadow: -1px -3px 5px 0px rgba(0, 0, 0, 0.04);
  }
`;
