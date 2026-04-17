import React from "react";

type CheckboxProps = {
  className?: string;
  checked?: boolean;
  onChange?: (event: React.ChangeEvent<HTMLInputElement>) => void;
  disabled?: boolean;
  id?: string;
  value?: string | number;
  label?: string;
};

const Checkbox: React.FC<CheckboxProps> = ({
  className = "",
  checked = false,
  onChange,
  disabled = false,
  id,
  value,
  label,
}) => {
  return (
    <label
      htmlFor={id}
      className={`inline-flex items-center cursor-pointer gap-2 ${
        disabled ? "opacity-50 cursor-not-allowed" : ""
      } ${className}`}
    >
      <input
        id={id}
        type="checkbox"
        checked={checked}
        onChange={onChange}
        disabled={disabled}
        value={value}
        className="peer sr-only"
      />
      <div
        className={`w-[17px] h-[17px] flex items-center justify-center border rounded-sm border-gray-400 
          peer-checked:bg-primary peer-checked:border-primary 
          peer-focus:ring-1 peer-focus:ring-primary
          transition-colors duration-200`}
      >
        {checked && (
          <svg
            className="w-3 h-3 text-white"
            fill="none"
            stroke="currentColor"
            strokeWidth="3"
            viewBox="0 0 24 24"
          >
            <path d="M5 13l4 4L19 7" />
          </svg>
        )}
      </div>
      {label && <span>{label}</span>}
    </label>
  );
};

export default Checkbox;
