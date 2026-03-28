/**
 *
 * BreadCrumb
 *
 */
import React from 'react';

// import styled from 'styles/styled-components';

import { FormattedMessage } from 'react-intl';
import messages from './messages';
import { Breadcrumbs, Typography, Link } from '@material-ui/core';

import { Wrapper } from './wrapper';
import { AppPages } from 'containers/App/constants';
import { useDispatch } from 'react-redux';
import { push } from 'redux-first-history';
import BreadIcon from 'images/themedIcons/breadIcon';

function BreadCrumb(props: {
  links: { pageName: string; pageLink: AppPages; last?: boolean }[];
}) {
  const dispatch = useDispatch();
  function handleClick(event: any, link: AppPages) {
    event.preventDefault();
    dispatch(push(link));
  }
  return (
    <>
      <Wrapper>
        <Breadcrumbs separator={<BreadIcon />} aria-label="breadcrumb">
          {props.links.map((item, index: number) => {
            if (item.last) {
              return (
                <Typography
                  className="lastBread"
                  key={`breadCrumb${index}`}
                  color="textPrimary"
                >
                  <FormattedMessage {...messages[item.pageName]} />
                </Typography>
              );
            }
            return (
              <Link
                key={`breadCrumb${index}`}
                className="costumeLink"
                onClick={e => handleClick(e, item.pageLink)}
              >
                <FormattedMessage {...messages[item.pageName]} />
              </Link>
            );
          })}
        </Breadcrumbs>
      </Wrapper>
    </>
  );
}

export default BreadCrumb;
