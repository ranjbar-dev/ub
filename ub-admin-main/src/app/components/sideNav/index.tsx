import { ListItem } from '@material-ui/core';
import List from '@material-ui/core/List';
import EqualizerIcon from '@material-ui/icons/Equalizer';
import { AppPages } from 'app/constants';
import { push } from 'connected-react-router';
import React, { memo, useState } from 'react';
import { useDispatch } from 'react-redux';
import styled from 'styled-components/macro';
import { StyleConstants } from 'styles/StyleConstants';


import MainCat from './mainCat';

/** A single navigation category with its child page links. */
interface NavCategory {
  name: string;
  icon: React.ReactNode;
  childs: { name: string; page?: string }[];
}

const DASHBOARD_ID = 'subCategoryDashboard';

/**
 * Left-side navigation panel with collapsible category items.
 * Dispatches push() actions to navigate between admin pages.
 *
 * @example
 * ```tsx
 * <SideNav mainCategories={[{ name: 'Users', icon: <PeopleIcon />, childs: [...] }]} />
 * ```
 */
const SideNav = memo(
  (props: {
    mainCategories: NavCategory[];
  }) => {
    const dispatch = useDispatch();
    const [OpenMenuIndex, setOpenMenuIndex] = useState(-1);
    const [activeItemId, setActiveItemId] = useState<string>(DASHBOARD_ID);

    const handleOpenClose = (index: number) => {
      if (index === OpenMenuIndex) {
        setOpenMenuIndex(-1);
        return;
      }
      setOpenMenuIndex(index);
    };

    const handleDashboardClick = () => {
      setActiveItemId(DASHBOARD_ID);
      dispatch(push(AppPages.HomePage));
    };

    const handleSelectChild = (id: string) => {
      setActiveItemId(id);
    };

    const isDashboardActive = activeItemId === DASHBOARD_ID;

    return (
      <SideNavWrapper>
        <List
          component="nav"
          aria-label="Main navigation"
          subheader={
            <ListItem
              component="div"
              role="button"
              tabIndex={0}
              className={`HomeButton subCategory${isDashboardActive ? ' active' : ''}`}
              id={DASHBOARD_ID}
              aria-current={isDashboardActive ? 'page' : undefined}
              onClick={handleDashboardClick}
              onKeyDown={(e: React.KeyboardEvent) => {
                if (e.key === 'Enter' || e.key === ' ') {
                  e.preventDefault();
                  handleDashboardClick();
                }
              }}
            >
              <span>
                <EqualizerIcon fontSize="small" />
              </span>
              Dashboard
            </ListItem>
          }
        >
          {props.mainCategories.map((item, index) => {
            return (
              <MainCat
                index={index}
                isOpen={index === OpenMenuIndex}
                onClick={handleOpenClose}
                key={item.name}
                icon={item.icon}
                title={item.name}
                childs={item.childs}
                activeItemId={activeItemId}
                onSelectChild={handleSelectChild}
              />
            );
          })}
        </List>
      </SideNavWrapper>
    );
  },
);
export { SideNav };
const SideNavWrapper = styled.div`
  width: ${StyleConstants.SIDE_NAV_WIDTH};
  height: 100%;
  margin-top: ${StyleConstants.NAV_BAR_HEIGHT};
  overflow: auto;
  min-width: ${StyleConstants.SIDE_NAV_WIDTH};
  .HomeButton {
    cursor: pointer;
    background: ${p => p.theme.white};
    color: ${p => p.theme.blackText};
    font-weight: 600;
    font-size: 12px;
    svg {
      margin-top: -3px;
      margin-left: 1px;
    }
    span {
      margin-right: 9px !important;
    }
    &.active span {
      color: ${p => p.theme.primary};
    }
  }
  .MuiListItemText-root {
    span {
      font-weight: 500 !important;
      font-size: 12px !important;
    }
  }
  .mainItem {
    span {
      font-weight: 600 !important;
    }
  }
  .MuiListItemText-root {
    padding-left: 8px;
  }
  .MuiList-subheader {
    padding-top: 22px !important;
  }
`;