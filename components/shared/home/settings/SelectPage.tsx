import React, { useState } from "react";
import { Menu, MenuButton, MenuItem, MenuItems } from "@headlessui/react";
import { MdKeyboardArrowDown } from "react-icons/md";

type Option = {
  id: number;
  name: string;
};

const options: Option[] = [
  { id: 1, name: "Checkout" },
  { id: 2, name: "Cart" },
];

export default function SelectPage() {
  const [selectedItem, setSelectedItem] = useState<Option | null>(null);

  return (
    <Menu>
      <MenuButton className="inline-flex items-center justify-between text-gray-600 border gap-2 rounded-md py-1.5 px-3 text-sm min-w-40 font-medium focus:outline-none data-[focus]:outline-1 data-[focus]:outline-white">
        {selectedItem?.name || "Select Option"}
        <MdKeyboardArrowDown className="size-5" />
      </MenuButton>
      <MenuItems
        transition
        anchor="bottom end"
        className="w-52 origin-top-right rounded-xl border bg-white p-1 text-sm transition duration-100 ease-out [--anchor-gap:var(--spacing-1)] focus:outline-none data-[closed]:scale-95 data-[closed]:opacity-0"
      >
        {options.map((item) => (
          <MenuItem key={item.id}>
            <button
              onClick={() => setSelectedItem(item)}
              className="group flex w-full items-center gap-2 rounded-lg py-1.5 px-3 data-[focus]:bg-gray-100"
            >
              {item.name}
            </button>
          </MenuItem>
        ))}
      </MenuItems>
    </Menu>
  );
}
