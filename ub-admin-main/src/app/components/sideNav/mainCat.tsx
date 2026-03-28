import {
ListItem,
ListItemText,
Collapse,
List,
ListItemIcon,
} from '@material-ui/core';
import {push} from 'connected-react-router';
import React,{memo} from 'react';
import {useDispatch} from 'react-redux';
import styled from 'styled-components/macro';

/** A single child page link inside a navigation category. */
interface NavChild {
  name: string;
  page?: string;
}

/**
 * Collapsible navigation category item for the side nav.
 * Renders a top-level item that expands to show child page links on click.
 *
 * @example
 * ```tsx
 * <MainCat
 *   title="Users"
 *   icon={<PeopleIcon />}
 *   childs={[{ name: 'Accounts', page: '/user-accounts' }]}
 *   index={0}
 *   isOpen={false}
 *   activeItemId=""
 *   onSelectChild={(id) => setActiveItemId(id)}
 *   onClick={(i) => setOpen(i)}
 * />
 * ```
 */
const MainCat=(props: {
title: string;
childs: NavChild[];
icon: React.ReactNode;
isOpen: boolean;
index: number;
onClick: (index: number) => void;
activeItemId: string;
onSelectChild: (id: string) => void;
}) => {
const {isOpen}=props;
const dispatch=useDispatch();
const handleClick=() => {
props.onClick(props.index);
};
return (
<Wrapper>
<ListItem button onClick={handleClick} aria-expanded={isOpen}>
<ListItemIcon>{props.icon}</ListItemIcon>
<ListItemText className="mainItem" primary={props.title} />
</ListItem>
<Collapse in={isOpen} timeout="auto">
<List component="div" disablePadding>
{props.childs.map((item,index) => {
const itemId='subCategory'+item.name.split(' ')[0];
const isActive=props.activeItemId===itemId;
return (
<ListItem
disableRipple
className={`subCategory${isActive ? ' active' : ''}`}
key={item.name}
disabled={!item.page}
id={itemId}
button
aria-current={isActive ? 'page' : undefined}
onClick={() => {
if(item.page) {
props.onSelectChild(itemId);
dispatch(push(item.page));
}
}}
>
<ListItemText className="menusubItem" primary={item.name} />
</ListItem>
);
})}
</List>
</Collapse>
</Wrapper>
);
};
export default memo(MainCat);
const Wrapper=styled.div`
  .subCategory {
    padding-left: 50px !important;
  }
  .MuiListItemIcon-root {
    min-width: unset !important;
  }
  .menuSubItem {
    span {
      font-size: 11px !important;
      font-weight: 600 !important;
    }
  }
  .subCategory.active span {
    color: ${p => p.theme.primary} !important;
    font-weight: 600 !important;
  }
`;