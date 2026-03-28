/**
 *
 * UserAccounts
 *
 */

import React, { memo } from "react";
import { useSelector } from "react-redux";

import UserAccountsPage from "./components/UserAccountsPage";
import VerificationPage from "./components/VerificationPage";
import { selectRouter } from "./selectors";

interface Props {}

export const UserAccounts = memo((props: Props) => {
  const router = useSelector(selectRouter);
  /////////////////////////

  const isVerificationPage =
    router && router.location.pathname.includes("verification");

  return isVerificationPage === true ? (
    <VerificationPage />
  ) : (
    <UserAccountsPage />
  );
});
