import { Switch } from "@headlessui/react";
import { useState } from "react";

interface ToggleSwitchProps {
  checked?: boolean;
  onChange?: (value: boolean) => void;
  className?: string;
}

export default function ToggleSwitch({
  checked,
  onChange,
  className = "",
}: ToggleSwitchProps) {
  const [internalChecked, setInternalChecked] = useState(false);
  const isControlled = typeof checked === "boolean";
  const enabled = isControlled ? checked : internalChecked;

  const handleChange = (val: boolean) => {
    if (!isControlled) setInternalChecked(val);
    onChange?.(val);
  };

  return (
    <Switch
      checked={enabled}
      onChange={handleChange}
      className={`group relative flex h-6 w-10 min-w-10 rounded-full bg-gray-200 p-1 transition-colors duration-200 ease-in-out focus:outline-none data-[focus]:outline-1 data-[focus]:outline-white data-[checked]:bg-primary ${className}`}
    >
      <span
        aria-hidden="true"
        className="pointer-events-none inline-block size-4 translate-x-0 rounded-full bg-white ring-0 shadow-lg transition duration-200 ease-in-out group-data-[checked]:translate-x-4"
      />
    </Switch>
  );
}
