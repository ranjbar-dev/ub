/**
 *
 * FilledOrders
 *
 */

import { ICellRendererParams, ValueFormatterParams, CellClassParams } from "ag-grid-community";
import { EmailButton } from "app/components/clickableEmail";
import { CellRenderer } from "app/components/renderer";
import { SimpleGrid } from "app/components/SimpleGrid/SimpleGrid";
import TitledContainer from "app/components/titledContainer/TitledContainer";
import { FullWidthWrapper } from "app/components/wrappers/FullWidthWrapper";
import { translations } from "locales/i18n";
import React, { memo, useMemo } from "react";
import { useTranslation } from "react-i18next";
import { useDispatch, useSelector } from "react-redux";
import { MessageNames } from "services/messageService";
import { CurrencyFormater } from "utils/formatters";
import { useInjectReducer, useInjectSaga } from "utils/redux-injectors";
import { stateStyler, cellColorAndNameFormatter } from "utils/stylers";

import { filledOrdersSaga } from "./saga";
import { FilledOrdersReducer, sliceKey, FilledOrdersActions } from "./slice";
import { selectFilledOrdersData } from "./selectors";

interface Props {}

export const FilledOrders = memo((props: Props) => {
  useInjectReducer({ key: sliceKey, reducer: FilledOrdersReducer });
  useInjectSaga({ key: sliceKey, saga: filledOrdersSaga });
  const dispatch = useDispatch();
  const filledOrdersData = useSelector(selectFilledOrdersData);

  const { t } = useTranslation();

  const staticRows = useMemo(
    () => [
      {
        headerName: t(translations.CommonTitles.ID()),
        field: "id",
        maxWidth: 100,
      },
      {
        headerName: t(translations.CommonTitles.UserId()),
        field: "userId",
        maxWidth: 130,
      },
      {
        headerName: t(translations.CommonTitles.UserEmail()),
        field: "userEmail",
        cellRenderer: ({ data, rowIndex }: ICellRendererParams) =>
          CellRenderer(<EmailButton dispatch={dispatch} data={data} />),
      },
      {
        headerName: t(translations.CommonTitles.RequestDate()),
        field: "createdAt",
      },
      {
        headerName: t(translations.CommonTitles.Pair()),
        field: "pair",
      },
      {
        headerName: t(translations.CommonTitles.Type()),
        field: "type",
        ...cellColorAndNameFormatter("type"),
      },
      {
        headerName: t(translations.CommonTitles.Side()),
        field: "side",
        cellStyle: (params: CellClassParams) => {
          return {
            color: stateStyler(params.data.side),
            textTransform: "capitalize",
          };
        },
        maxWidth: 110,
      },
      {
        headerName: t(translations.CommonTitles.Price()),
        field: "price",
        valueFormatter: (params: ValueFormatterParams) => {
          return CurrencyFormater(params.data.price);
        },
      },
      {
        headerName: t(translations.CommonTitles.Amount()),
        field: "amount",
        valueFormatter: (params: ValueFormatterParams) => {
          return CurrencyFormater(params.data.amount);
        },
      },
      {
        headerName: t(translations.CommonTitles.Total()),
        field: "total",
        valueFormatter: (params: ValueFormatterParams) => {
          return CurrencyFormater(params.data.total);
        },
      },
      {
        headerName: t(translations.CommonTitles.Filled()),
        field: "executed",
      },
      {
        headerName: t(translations.CommonTitles.TriggerCondition()),
        field: "triggerCondition",
      },
    ],
    []
  );
  return (
    <FullWidthWrapper>
      <TitledContainer
        id={"filledOrders"}
        title={t(translations.CommonTitles.FilledOrders())}
      >
        <SimpleGrid
          containerId="filledOrders"
          additionalInitialParams={{ status: "filled" }}
          arrayFieldName="orders"
          immutableId="id"
          filters={{
            dropDownCols: [
              {
                id: "side",
                substituteId: "type",
                options: [
                  {
                    name: "Buy",
                    value: "buy",
                  },
                  {
                    name: "Sell",
                    value: "sell",
                  },
                ],
              },
            ],
          }}
          //  onRowClick={handleRowClick}
          initialAction={FilledOrdersActions.GetFilledOrdersAction}
          messageName={MessageNames.SET_FILLED_ORDERS_DATA}
          externalData={filledOrdersData}
          staticRows={staticRows}
        />
      </TitledContainer>
    </FullWidthWrapper>
  );
});
