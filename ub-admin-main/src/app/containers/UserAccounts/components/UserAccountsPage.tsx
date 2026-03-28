import React, { memo, useCallback, useMemo } from 'react';
import { FullWidthWrapper } from 'app/components/wrappers/FullWidthWrapper';
import TitledContainer from 'app/components/titledContainer/TitledContainer';
import { useTranslation } from 'react-i18next';
import { translations } from 'locales/i18n';
import { SimpleGrid } from 'app/components/SimpleGrid/SimpleGrid';
import { UserAccountsActions } from '../slice';
import { MessageNames, GridNames } from 'services/messageService';
import { CellRenderer } from 'app/components/renderer';
import { RowClickedEvent, ValueFormatterParams, ICellRendererParams } from 'ag-grid-community';
import { useDispatch, useSelector } from 'react-redux';
import { WindowTypes } from 'app/constants';
import { selectUserAccountsData } from '../selectors';

interface Props {}
function UserAccountsPage(props: Props) {
  const { t } = useTranslation();
  const dispatch = useDispatch();
  const userAccountsData = useSelector(selectUserAccountsData);
  const handleRowClick = useCallback((e: RowClickedEvent) => {
    dispatch(
      UserAccountsActions.getInitialSingleUserDataAndOpenWindowAction({
        id: e.data.id,
        windowType: WindowTypes.User,
      }),
    );
  }, []);
  const staticRows = useMemo(
    () => [
      {
        headerName: t(translations.Grid.IDNO()),
        field: 'id',
        maxWidth: 150,
      },
      {
        headerName: t(translations.Grid.EmailAddress()),
        field: 'email',
        minWidth: 250,
      },
      {
        headerName: t(translations.Grid.FirstName()),
        field: 'firstName',
      },
      {
        headerName: t(translations.Grid.LastName()),
        field: 'lastName',
      },
      {
        headerName: t(translations.Grid.Country()),
        field: 'country',
        valueFormatter: (params: ValueFormatterParams) => {
          return params.data.countryFullName;
        },
        minWidth: 110,
      },
      {
        headerName: t(translations.Grid.ReferKey()),
        field: 'referKey',
        maxWidth: 130,
        minWidth: 100,
      },
      {
        headerName: t(translations.Grid.ReferalID()),
        field: 'referralId',
      },
      {
        headerName: t(translations.Grid.Manager()),
        field: 'manager',
        maxWidth: 130,
        minWidth: 100,
      },
      {
        headerName: t(translations.Grid.RegisterDate()),
        field: 'registrationDate',
        minWidth: 120,
      },
      {
        headerName: t(translations.Grid.RegisterIP()),
        field: 'registeredIP',
      },
      {
        width: 0,
        field: 'loading',
        cellRenderer: ({ data, rowIndex }: ICellRendererParams) =>
          CellRenderer(
            <>
              <div
                className="loadingGridRow"
                id={'loading' + data.id}
                style={{ display: 'none' }}
              ></div>
            </>,
          ),
      },
    ],

    [],
  );
  return (
    <FullWidthWrapper>
      <TitledContainer
        id="userAccounts"
        title={t(translations.CommonTitles.UserAccounts())}
      >
        <SimpleGrid
          containerId="userAccounts"
          additionalInitialParams={{}}
          arrayFieldName="users"
          immutableId="id"
          filters={{
            countryCols: ['country'],
            dateCols: ['registrationDate'],
          }}
          onRowClick={handleRowClick}
          initialAction={UserAccountsActions.GetInitialUserAccountsAction}
          messageName={MessageNames.SET_USER_ACCOUNTS}
          externalData={userAccountsData}
          staticRows={staticRows}
        />
      </TitledContainer>
    </FullWidthWrapper>
  );
}

export default memo(UserAccountsPage);
