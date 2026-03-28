import { Dispatch } from "@reduxjs/toolkit";
import { WindowTypes } from "app/constants";
import { UserAccountsActions } from "app/containers/UserAccounts/slice";
import React from "react";

/** Row data shape for the email button cell renderer. */
interface EmailCellData {
  userId: number;
  [key: string]: unknown;
}

/**
 * Clickable email cell renderer for AG Grid.
 * Opens the UserDetails window for the associated user on click.
 *
 * @example
 * ```tsx
 * <EmailButton data={rowData} dispatch={dispatch} fieldName="userEmail" />
 * ```
 */
export const EmailButton = ({
  data,
  dispatch,
  fieldName,
}: {
  data: EmailCellData;
  dispatch: Dispatch<{ type: string; payload?: unknown }>;
  fieldName?: string;
}) => {
  const handleEmailClick = () => {
    dispatch(
      UserAccountsActions.getInitialSingleUserDataAndOpenWindowAction({
        id: data.userId,
        windowType: WindowTypes.User,
      })
    );
  };

  return (
    <div style={{ cursor: "pointer" }} onClick={handleEmailClick}>
      <span>{data[fieldName ?? "userEmail"]}</span>
    </div>
  );
};
