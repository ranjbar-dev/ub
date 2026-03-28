import { useState, useEffect } from 'react';
export type loginData = {
  email?: string;
  password?: string;
};

export const UseFormValidation = (initialState: any, validate: any): any => {
  const [values, setValues] = useState(initialState);
  const [errors, setErrors]: [any, any] = useState();
  const [isSubmitting, setSubmitting] = useState(false);
  useEffect(() => {
    setSubmitting(false);
  }, [errors]);
  const handleBlur = () => {
    const validationErrors = validate(values);
    setErrors(validationErrors);
  };
  const handleChange = (event: any) => {
    event.preventDefault();
    setValues({
      ...values,
      [event.target.name]: event.target.value,
    });
  };
  const handleSubmit = () => {
    const validationErrors = validate(values);
    setErrors(validationErrors);
    setSubmitting(true);
  };
  let hasError: boolean = false;
  if (errors != null) {
    const errorsList: any = Object.values(errors);
    for (let i = 0; i < errorsList.length; i++) {
      if (errorsList[i].length > 0) {
        hasError = true;
      }
    }
  }
  if (typeof errors == 'undefined') {
    hasError = true;
  }
  return {
    values,
    handleChange,
    handleSubmit,
    handleBlur,
    errors,
    isSubmitting,
    hasError,
  };
};
