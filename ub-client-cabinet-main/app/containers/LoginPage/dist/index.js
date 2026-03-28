"use strict";
/*
 *
 * LoginPage
 *
 */
var __assign = (this && this.__assign) || function () {
    __assign = Object.assign || function(t) {
        for (var s, i = 1, n = arguments.length; i < n; i++) {
            s = arguments[i];
            for (var p in s) if (Object.prototype.hasOwnProperty.call(s, p))
                t[p] = s[p];
        }
        return t;
    };
    return __assign.apply(this, arguments);
};
exports.__esModule = true;
var react_1 = require("react");
var react_helmet_1 = require("react-helmet");
var fullPageWrapper_1 = require("components/wrappers/fullPageWrapper");
var LocaleToggle_1 = require("containers/LocaleToggle");
var loginBody_1 = require("./loginBody");
function LoginPage(props) {
    var isPopup = props.isPopup;
    return (react_1["default"].createElement(react_1["default"].Fragment, null, !isPopup ? (react_1["default"].createElement(fullPageWrapper_1.FullPageWrapper, null,
        react_1["default"].createElement(react_helmet_1.Helmet, null,
            react_1["default"].createElement("title", null, "Login Page"),
            react_1["default"].createElement("meta", { name: "description", content: "Description of LoginPage" })),
        react_1["default"].createElement("div", { className: "head darkTheme WhiteHeader" },
            react_1["default"].createElement(LocaleToggle_1["default"], null)),
        react_1["default"].createElement("div", { className: "body" }, react_1["default"].createElement(loginBody_1["default"], null)))) : (react_1["default"].createElement(loginBody_1["default"], __assign({}, props)))));
}
exports["default"] = react_1.memo(LoginPage);
