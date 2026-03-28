import React, { memo } from 'react';
import styled from 'styled-components/macro';

interface Props {
  title: string;
  children: React.ReactNode;
}

function DashboardCard(props: Props) {
  const { title, children } = props;

  return (
    <Wrapper>
      <div className="cardTitle">{title}</div>
      <div className="content">{children}</div>
    </Wrapper>
  );
}

export default memo(DashboardCard);
const Wrapper = styled.div`
  width: 310px;
  height: 230px;
  background: white;
  border-radius: 7px;
  margin: 24px;
  box-shadow: 0 0 11px 0px rgb(0 0 0 / 6%);
  padding-top: 12px;
  .cardTitle {
    width: 100%;
    background: #e4eef6;
    /* margin-top: 12px; */
    padding: 12px 24px;
    font-size: 15px;
    font-weight: 600;
  }
`;
