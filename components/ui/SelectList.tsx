import React from "react";
import {
  Listbox,
  ListboxButton,
  ListboxOption,
  ListboxOptions,
} from "@headlessui/react";
import { cn } from "@/lib/cn";
import { MdKeyboardArrowDown } from "react-icons/md";

interface SelectListProps {
  options?: {
    id: number;
    name: string;
    value?: string;
  }[];
  value?: { id: number; name: string; value: string } | null;
  onChange: (val: { id: number; name: string; value: string }) => void;
  className?: string;
  placeholder?: string;
}

export default function SelectList({
  options = [],
  value,
  onChange,
  className,
  placeholder = "Select an option",
}: SelectListProps) {
  return (
    <Listbox value={value} onChange={onChange}>
      <ListboxButton
        className={cn(
          "relative rounded-md border bg-white py-1.5 px-3 text-left text-sm flex items-center justify-between",
          "focus:outline-none data-[focus]:outline-2 data-[focus]:-outline-offset-2 data-[focus]:outline-white/25",
          className
        )}
      >
        {value?.name || <span className="text-gray-400">{placeholder}</span>}
        <MdKeyboardArrowDown size={20} className="text-gray-500 ml-2" />
      </ListboxButton>
      <ListboxOptions
        anchor="bottom"
        transition
        className={cn(
          "w-[var(--button-width)] z-10 rounded-xl border bg-white p-1 [--anchor-gap:var(--spacing-1)] focus:outline-none",
          "transition duration-100 ease-in data-[leave]:data-[closed]:opacity-0"
        )}
      >
        {options.length > 0 ? (
          options.map((item) => (
            <ListboxOption
              key={item.id}
              value={item}
              className="group flex cursor-default items-center gap-2 rounded-lg py-1.5 px-3 select-none data-[focus]:bg-gray-100 data-[selected]:bg-gray-100"
            >
              <div className="text-sm font-medium">{item.name}</div>
            </ListboxOption>
          ))
        ) : (
          <div className="py-2 px-3 text-sm text-gray-400">
            No options available
          </div>
        )}
      </ListboxOptions>
    </Listbox>
  );
}
