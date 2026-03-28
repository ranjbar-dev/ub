/**
 *
 * OpenOrders
 *
 */

import { Dispatch } from "@reduxjs/toolkit";
import { ColDef, ICellRendererParams, ValueFormatterParams, CellClassParams } from "ag-grid-community";
import { EmailButton } from "app/components/clickableEmail";
import IsLoadingWithTextAuto from "app/components/isLoadingWithText/isLoadingWithTextAuto";
import { CellRenderer } from "app/components/renderer";
import { SimpleGrid } from "app/components/SimpleGrid/SimpleGrid";
import TitledContainer from "app/components/titledContainer/TitledContainer";
import { FullWidthWrapper } from "app/components/wrappers/FullWidthWrapper";
import { Buttons, WindowTypes } from "app/constants";
import { translations } from "locales/i18n";
import React, { memo, useMemo, useCallback } from "react";
import { useTranslation } from "react-i18next";
import { useDispatch, useSelector } from "react-redux";
import { MessageNames } from "services/messageService";
import styled from "styled-components/macro";
import { CurrencyFormater } from "utils/formatters";
import { useInjectReducer, useInjectSaga } from "utils/redux-injectors";
import { stateStyler, cellColorAndNameFormatter } from "utils/stylers";

import { openOrdersSaga } from "./saga";
import { OpenOrdersReducer, sliceKey, OpenOrdersActions } from "./slice";
import { selectOpenOrdersData } from "./selectors";
import { UserAccountsActions } from "../UserAccounts/slice";

interface Props {}

export const OpenOrders = memo((props: Props) => {
  useInjectReducer({ key: sliceKey, reducer: OpenOrdersReducer });
  useInjectSaga({ key: sliceKey, saga: openOrdersSaga });
  const dispatch = useDispatch();
  const openOrdersData = useSelector(selectOpenOrdersData);
  const { t } = useTranslation();
  const handleCancelClick = useCallback(
    (data: Record<string, unknown>) => {
      dispatch(
        OpenOrdersActions.CancelOpenOrderAction({ id: data.id, type: "order" })
      );
    },
    [dispatch]
  );
  const handleFullFillClick = useCallback(
    (data: Record<string, unknown>) => {
      dispatch(
        OpenOrdersActions.FullFillOpenOrderAction({
          id: data.id,
          type: "order",
        })
      );
    },
    [dispatch]
  );
  const staticRows: ColDef[] = useMemo(
    () => [
      {
        headerName: t(translations.CommonTitles.ID()),
        field: "id",
        maxWidth: 90,
        minWidth: 70,
      },
      {
        headerName: t(translations.CommonTitles.UserId()),
        field: "userId",
        maxWidth: 140,
        minWidth: 90,
      },
      {
        headerName: t(translations.CommonTitles.UserEmail()),
        field: "userEmail",
        minWidth: 90,
        cellRenderer: ({ data, rowIndex }: ICellRendererParams) =>
          CellRenderer(<EmailButton dispatch={dispatch} data={data} />),
      },
      {
        headerName: t(translations.CommonTitles.RequestDate()),
        field: "createdAt",
        minWidth: 160,
      },
      {
        headerName: t(translations.CommonTitles.Pair()),
        field: "pair",
        maxWidth: 150,
        minWidth: 90,
      },
      {
        headerName: t(translations.CommonTitles.Type()),
        field: "type",
        maxWidth: 120,
        minWidth: 90,
        ...cellColorAndNameFormatter("type"),
      },
      {
        headerName: t(translations.CommonTitles.Side()),
        field: "side",
        maxWidth: 120,
        minWidth: 90,
        cellStyle: (params: CellClassParams) => {
          return {
            color: stateStyler(params.data.side),
            textTransform: "capitalize",
          };
        },
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
      {
        headerName: "",
        field: "actions",
        minWidth: 170,
        cellRenderer: ({ data, rowIndex }: ICellRendererParams) =>
          CellRenderer(
            <div className="gridActionButton">
              {/*<Button className={Buttons.SkyBlueButton}>
                {t(translations.CommonTitles.FullFilled())}
              </Button>*/}
              <IsLoadingWithTextAuto
                text={t(translations.CommonTitles.FullFilled())}
                className={Buttons.SkyBlueButton}
                loadingId={"fullFillButton" + data.id}
                onClick={() => {
                  handleFullFillClick(data);
                }}
              />
              <IsLoadingWithTextAuto
                text={t(translations.CommonTitles.Cancel())}
                className={Buttons.RedButton}
                loadingId={"cancelButton" + data.id}
                onClick={() => {
                  handleCancelClick(data);
                }}
              />
            </div>
          ),
      },
    ],

    []
  );

  return (
    <FullWidthWrapper>
      <TitledContainer
        id={"openOrders"}
        title={t(translations.CommonTitles.OpenOrders())}
      >
        <SimpleGrid
          containerId="openOrders"
          additionalInitialParams={{}}
          arrayFieldName="orders"
          immutableId="id"
          filters={{
            hiddenCols: ["actions"],
            disabledCols: ["price"],
            dateCols: ["createdAt"],
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
          initialAction={OpenOrdersActions.GetOpenOrdersAction}
          messageName={MessageNames.SET_OPEN_ORDERS_PAGE_DATA}
          externalData={openOrdersData}
          staticRows={staticRows}
        />
      </TitledContainer>
    </FullWidthWrapper>
  );
});

const Wrapper = styled.div``;
