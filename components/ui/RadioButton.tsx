import React from "react";

interface RadioButtonProps {
  name?: string;
  value?: string;
  checked?: boolean;
  onChange?: (event: React.ChangeEvent<HTMLInputElement>) => void;
  disabled?: boolean;
  id?: string;
}

export default function RadioButton({
  name,
  value,
  checked,
  onChange,
  disabled = false,
  id,
}: RadioButtonProps) {
  const inputId = id || `${name}-${value}`;

  return (
    <input
      id={inputId}
      type="radio"
      name={name}
      value={value}
      checked={checked}
      onChange={onChange}
      disabled={disabled}
      className="h-4 w-4 min-w-4 text-primary rounded focus:ring-primary border-gray-300 disabled:cursor-not-allowed"
    />
  );
}
