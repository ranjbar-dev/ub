import { currencyMap } from './sharedData';

export function queryStringer (params: any): string {
  let qs = '?';
  let counter = 0;
  for (const key in params) {
    const prefix = counter === 0 ? '' : '&';
    qs += prefix + key + '=' + params[key];
    counter++;
  }
  if (qs == '?') {
    return '';
  }
  return qs;
}

export const CurrencyFormater = (val: string) => {
  let value = val;
  if (!value) {
    return '';
  }
  if (value.includes('%')) {
    return value;
  }
  let trail = '';
  if (value.includes(' ')) {
    if (+value.split[0] < 1) {
      return value;
    }
    value = +value.split(' ')[0] + ' ' + value.split(' ')[1];
    trail = ' ' + value.split(' ')[1];
  } else {
    if (+value < 1) {
      return value + '';
    }
    value = +value + '';
  }
  if (!value.split('.')[1]) {
    value = Number(value.split(' ')[0]).toFixed(2) + trail;
  } else if (value.split(' ')[0].split('.')[1].length < 3) {
    value = Number(value.split(' ')[0]).toFixed(2) + trail;
  }

  const separated = value.split('.');
  const comaSeparated = separated[0].replace(/\B(?=(\d{3})+(?!\d))/g, ',');
  if (separated[1]) {
    return comaSeparated + '.' + separated[1];
  }
  return comaSeparated;
};
export const commaSeparatedInput = (value: any) => {
  let haveDot = false;
  let completeValue = value.internal ?? value.replace(/,/g, '');
  if (value.includes('.')) {
    haveDot = true;
    completeValue = value
      .split('.')[0]
      .replace(/,/g, '')
      .replace(/\B(?=(\d{3})+(?!\d))/g, ',');
    return completeValue + '.' + (value.split('.')[1] ?? '');
  }
  return completeValue.replace(/,/g, '').replace(/\B(?=(\d{3})+(?!\d))/g, ',');
};
export function Format (value: number | string) {
  const v = value ? Number(value) : 0;
  if (v > 0) {
    return Intl.NumberFormat('de-DE')
      .format(v)
      .split('.')
      .join(',');
  } else if (v === 0) {
    return '';
  } else {
    return (
      '(' +
      Intl.NumberFormat('de-DE')
        .format(-1 * v)
        .split('.')
        .join(',') +
      ')'
    );
  }
}

export function FormatDate (value: number | string) {
  const v = value ? value + '' : '';
  if (v.length > 0) {
    return (
      v.substring(0, 4) + '-' + v.substring(4, 6) + '-' + v.substring(6, 8)
    );
  } else {
    return '';
  }
}
export const DatePrefixer = (number: number) => {
  return number < 10 ? '0' + number : '' + number;
};

export const vw = (percent: number, viewportWidth: number) => {
  const per = (percent / 100) * viewportWidth + '';
  return Number(per.split('.')[0]);
};
// export const columnResize = (data: {
//   gridColumnApi: any;
//   resizeLimit: number;
// }) => {
//   const width = window.innerWidth;
//   if (data.resizeLimit > width) {
//     let allColumnIds: any[] = [];
//     data.gridColumnApi.getAllColumns().forEach(function(column: any) {
//       allColumnIds.push(column.colId);
//     });
//     data.gridColumnApi.autoSizeColumns(allColumnIds, true);
//   }
// };
export const censor = (item: { value: string; from: number; to: number }) => {
  const transmitted = item.value.split('');
  for (let i = item.from; i <= item.to; i++) {
    transmitted[i] = '*';
  }
  return transmitted.join('');
};
export const CopyToClipboard = (text: string) => {
  const dummy = document.createElement('textarea');
  document.body.appendChild(dummy);
  dummy.value = text;
  dummy.select();
  document.execCommand('copy');
  document.body.removeChild(dummy);
};
export const Translator = (translateData: {
  intl: any;
  containerPrefix: string;
  message: string;
}) => {
  return translateData.intl.formatMessage({
    id: translateData.containerPrefix + '.' + translateData.message,
    defaultMessage: 'ET.' + translateData.message,
  });
};

export const PairFormat = (pair: string) => {
  return pair.replace('-', '');
};

export const UpperFirstRegx = /(\b[a-z](?!\s))/g;

export const UpperFirstLetters = (s: string) => {
  return s.replace(UpperFirstRegx, function (x) {
    return x.toUpperCase();
  });
};

export const removeTrailingZeroes = (value: string) => {
  value = value ? value.toString() : '';
  if (value !== '') {
    if (value.indexOf('.') === -1) {
      return value;
    }

    let cutFrom = value.length - 1;
    do {
      if (value[cutFrom] === '0') {
        cutFrom--;
      }
    } while (value[cutFrom] === '0');
    if (value[cutFrom] === '.') {
      cutFrom--;
    }
    return value.substr(0, cutFrom + 1);
  }
  return '';
};

export const zeroFixer = (value: string) => {
  let tmp;
  if (value) {
    if (value.split('.')[0] !== '0') {
      tmp = CurrencyFormater(value);
    } else {
      tmp = removeTrailingZeroes(value);
    }
  }
  return tmp;
};
export const toFraction = (v1: string | number, to: number) => {
  let tmp = v1 + '';
  if (tmp.includes('.')) {
    const sep = tmp.split('.');
    if (sep[1].length > to) {
      sep[1] = sep[1].substring(0, to);
    }
    tmp = sep.join('.');
  }
  return tmp;
};

export function toFixedWithoutScientificNotation (x: any) {
  if (Math.abs(x) < 1.0) {
    var e = parseInt(x.toString().split('e-')[1]);
    if (e) {
      x *= Math.pow(10, e - 1);
      x = '0.' + new Array(e).join('0') + x.toString().substring(2);
    }
  } else {
    var e = parseInt(x.toString().split('+')[1]);
    if (e > 20) {
      e -= 20;
      x /= Math.pow(10, e);
      x += new Array(e + 1).join('0');
    }
  }
  return x;
}
export function numberToString (num) {
  let numStr = String(num);

  if (Math.abs(num) < 1.0) {
    const e = parseInt(num.toString().split('e-')[1]);
    if (e) {
      const negative = num < 0;
      if (negative) num *= -1;
      num *= Math.pow(10, e - 1);
      numStr = '0.' + new Array(e).join('0') + num.toString().substring(2);
      if (negative) numStr = '-' + numStr;
    }
  } else {
    let e = parseInt(num.toString().split('+')[1]);
    if (e > 20) {
      e -= 20;
      num /= Math.pow(10, e);
      numStr = num.toString() + new Array(e + 1).join('0');
    }
  }

  return numStr;
}
export const formatSmall = (val: string, toFix = 8) => {
  let value = val;
  if (!value) {
    return '';
  }
  if (value.includes('%')) {
    return value;
  }
  if (isNaN(Number(value))) {
    return '';
  }
  if (Number(value) < 0) {
    return '';
  }
  if (Number(value) < 1) {
    value = numberToString(Number(value)).toString();
    const splitted = value.split('.');
    if (splitted[1]) {
      value = splitted[0] + '.' + splitted[1].substring(0, toFix);
    }
    return value;
  } else {
    const splitted = value.split('.');
    if (splitted[1]) {
      value = splitted[0] + '.' + splitted[1].substring(0, toFix);
    }
    return CurrencyFormater(value);
  }
};
const addComma = (s: string) => {
  return s.replace(/\B(?=(\d{3})+(?!\d))/g, ',');
};
export const formatCurrencyWithMaxFraction = (
  s: string | undefined,
  maxFraction = 8,
) => {
  if (s == undefined) {
    return '';
  }
  let formatted = s;
  if (s.includes('e')) {
    return formatSmall(s, maxFraction);
  }
  if (s.includes('.')) {
    if (s.startsWith('.')) {
      return '0' + s;
    }
    const splitted = s.split('.');
    splitted[0] = addComma(splitted[0]);
    if (splitted[1]) {
      formatted = splitted[0] + '.' + splitted[1].substring(0, maxFraction);
    } else {
      formatted = splitted[0] + '.';
    }
    return formatted;
  }
  return formatted.replace(/\B(?=(\d{3})+(?!\d))/g, ',');
};
export const formatTableCell = ({ value }) => {
  const splitted = value.split(' ');
  if (splitted[0] === '0') {
    return `0.00 ${splitted[1] ?? ''}`;
  }
  return `${formatCurrencyWithMaxFraction(
    splitted[0],
    currencyMap({ code: splitted[1] })?.showDigits ?? 8,
  )} ${splitted[1] ?? ''}`;
};
