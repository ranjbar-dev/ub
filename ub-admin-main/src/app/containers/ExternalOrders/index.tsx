/**
 *
 * ExternalOrders
 *
 */

import GridTabs from 'app/components/GridTabs/GridTabs';
import TitledContainer from 'app/components/titledContainer/TitledContainer';
import { FullWidthWrapper } from 'app/components/wrappers/FullWidthWrapper';
import { translations } from 'locales/i18n';
import React, { memo, useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';
import styled from 'styled-components/macro';
import { useInjectReducer, useInjectSaga } from 'utils/redux-injectors';

import { externalOrdersSaga } from './saga';
import {
	ExternalOrdersReducer,
	sliceKey,
} from './slice';
import AllQueue from './tabPages/AllQueue';
import Executed from './tabPages/Executed';
import NetQueue from './tabPages/NetQueue';


interface Props { }

export const ExternalOrders = memo((props: Props) => {
	useInjectReducer({ key: sliceKey, reducer: ExternalOrdersReducer });
	useInjectSaga({ key: sliceKey, saga: externalOrdersSaga });




	const { t } = useTranslation();
	// const [ActivePage,setActivePage]=useState<JSX.Element>(<Executed />);
	const [ActivePage, setActivePage] = useState<JSX.Element>(<Executed />);

	const filterTabs = [
		{
			name: t(translations.CommonTitles.Executed()),
			callObject: <Executed />,
		},
		{
			name: t(translations.CommonTitles.NetQueue()),
			callObject: <NetQueue />,
		},
		{
			name: t(translations.CommonTitles.AllQueue()),
			callObject: <AllQueue />
		},

	];

	const handleTopTabChange = (e: Record<string, unknown>) => {
		setActivePage(e as unknown as JSX.Element)
	}
	return (
		<FullWidthWrapper>

			<TitledContainer
				id={'externalOrders'}
				title={t(translations.CommonTitles.ExternalOrders())}
			>
				<GridTabs styles={{ marginTop: '-31px' }} onChange={handleTopTabChange} tabs={filterTabs} />
				{ActivePage}
			</TitledContainer>
		</FullWidthWrapper>
	);
});

const Wrapper = styled.div``;
