import { cn } from "@/lib/cn";
import React, { useState } from "react";
import { Eye, EyeOff } from "lucide-react";

interface InputFieldProps extends React.InputHTMLAttributes<HTMLInputElement> {
  className?: string;
  error?: string;
  wrapperClass?: string;
}

const InputField = React.forwardRef<HTMLInputElement, InputFieldProps>(
  ({ type = "text", className = "", wrapperClass, error, ...rest }, ref) => {
    const [showPassword, setShowPassword] = useState(false);
    const isPassword = type === "password";

    return (
      <div className={cn("w-full relative", wrapperClass)}>
        <input
          type={isPassword && showPassword ? "text" : type}
          ref={ref}
          aria-invalid={!!error}
          className={cn(
            "bg-white border rounded-md text-sm px-3 py-1.5 focus:outline-none focus:ring-1",
            error
              ? "border-red-500 focus:ring-red-500"
              : "border-gray-300 focus:ring-primary",
            className
          )}
          {...rest}
        />
        {isPassword && (
          <button
            type="button"
            className="absolute inset-y-0 right-2 flex items-center text-gray-500"
            onClick={() => setShowPassword((prev) => !prev)}
            tabIndex={-1}
          >
            {showPassword ? (
              <Eye className="w-4 h-4" />
            ) : (
              <EyeOff className="w-4 h-4" />
            )}
          </button>
        )}
        {error && (
          <p className="mt-1 text-sm text-red-500" role="alert">
            {error}
          </p>
        )}
      </div>
    );
  }
);

InputField.displayName = "InputField";

export default InputField;
