import { cn } from "@/lib/cn";
import Link from "next/link";
import React, { ReactNode } from "react";

type ButtonProps = {
  children: ReactNode;
  onClick?: () => void;
  type?: "button" | "submit" | "reset";
  variant?: "primary" | "secondary" | "outline";
  disabled?: boolean;
  className?: string;
  link?: boolean;
  src?: string;
};

const Button: React.FC<ButtonProps> = ({
  children,
  onClick,
  type = "button",
  variant = "primary",
  disabled = false,
  className = "",
  link,
  src,
}) => {
  const baseStyle =
    "px-4 py-2 rounded-lg text-sm font-medium transition duration-200 ease-in-out focus:outline-none flex items-center gap-2";
  const variants: Record<string, string> = {
    primary: "bg-primary text-white hover:bg-blue-700",
    secondary: "border hover:bg-gray-100",
    outline: "border border-gray-400 text-gray-700 hover:bg-gray-100",
  };

  return !link ? (
    <button
      type={type}
      onClick={onClick}
      disabled={disabled}
      className={cn(
        "",
        baseStyle,
        variants[variant],
        disabled ? "opacity-50 cursor-not-allowed" : "",
        className
      )}
    >
      {children}
    </button>
  ) : (
    <Link
      href={src || ""}
      className={cn(
        "",
        baseStyle,
        variants[variant],
        disabled ? "opacity-50 cursor-not-allowed" : "",
        className
      )}
    >
      {children}
    </Link>
  );
};

export default Button;
