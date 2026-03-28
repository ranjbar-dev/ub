//import * as Decimal from 'utils/decimal.js';
import { add, divide, multiply, subtract } from 'precise-math';

export const removeComma = (val: any) => {
  if (val === '' || val === '0' || val === '0.00' || val === 0) {
    return '0';
  }
  const value = val.internal ?? val;
  return Number(value.replace(/,/g, ''));
};
export const IsNumber = val => {
  if (
    !Number(val) &&
    !removeComma(val) &&
    val !== '' &&
    val !== '0' &&
    !val.includes('.')
  ) {
    return false;
  }
  return true;
};

const validateNumber = (number: string | number) => {
  let tmp = number;
  if (typeof tmp === 'string') {
    tmp = tmp.replace(/,/g, '');
    const splited = tmp.split('.');
    if (splited.length > 2) {
      tmp = splited[0] + '.' + splited[1];
    }
  }
  if (isNaN(Number(tmp)) || !isFinite(Number(tmp))) {
    tmp = 0;
  }
  return Number(tmp);
};

export const numberCheck = (v1, v2) => {
  const nv1 = validateNumber(v1);
  const nv2 = validateNumber(v2);

  return [nv1, nv2];
};

export const Multiply = (v1, v2) => {
  const [nv1, nv2] = numberCheck(v1, v2);
  return multiply(nv1, nv2) + '';
  //Decimal(nv1).mul(nv2).internal
};
export const Divide = (v1, v2) => {
  const [nv1, nv2] = numberCheck(v1, v2);
  if (nv2 === 0) {
    return '0';
  }
  return divide(nv1, nv2) + '';
  // Decimal(nv1).div(nv2).internal
};

export const Subtract = (v1, v2) => {
  const [nv1, nv2] = numberCheck(v1, v2);
  return subtract(nv1, nv2) + '';
  // Decimal(nv1).sub(nv2).internal
};

export const Add = (v1, v2) => {
  const [nv1, nv2] = numberCheck(v1, v2);
  return add(nv1, nv2) + '';

  //Decimal(nv1).plus(nv2).internal
};

export const formatNumber = function (n: any) {
  const clean = (n.internal ?? n).replace(/,/g, '');
  const separated = clean.split('.');
  const comaSeparated = separated[0].replace(/\B(?=(\d{3})+(?!\d))/g, ',');
  if (clean.includes('.')) {
    return comaSeparated + '.' + separated[1];
  }
  return comaSeparated;
};

export const toFix = (value: string | number, to: number) => {
  return Number(value).toFixed(to); //(Math.floor(Number(Number(value)*(10*to)))/(10*to))+''
};
export const toFixFloor = (value: string | number, to: number) => {
  return Math.floor(Number(Number(value) * (10 * to))) / (10 * to) + '';
};
