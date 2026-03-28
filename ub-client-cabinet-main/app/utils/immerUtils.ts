import produce from 'immer';

// remove first item
export const removeFirst = ({ array }) =>
  produce(array, draft => {
    draft.shift();
  });

// add item at the beginning of the array
export const addAtFirst = ({ array, item }) =>
  produce(array, draft => {
    draft.unshift(item);
  });

// remove last item
export const removeLast = ({ array }) =>
  produce(array, draft => {
    draft.pop();
  });

export const addAtFirstAndRemoveLast = ({ array, item }) =>
  produce(array, draft => {
    draft.unshift(item);
    draft.pop();
  });
